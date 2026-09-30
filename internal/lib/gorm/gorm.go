package gormv2

import (
	"context"
	"errors"
	"fmt"
	"github.com/SisyphusSQ/go-starter/v2/config"
	"github.com/SisyphusSQ/go-starter/v2/internal/health"
	"github.com/SisyphusSQ/go-starter/v2/internal/lib/log"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"go.uber.org/fx"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net"
	"strconv"
	"time"
)

type Engine struct{ gorm *gorm.DB }

func New(lifecycle fx.Lifecycle, c config.Config) (*Engine, error) {
	conf, err := normalizeConfig(c.Database)
	if err != nil {
		return nil, err
	}
	location, err := time.LoadLocation(conf.TimeZone)
	if err != nil {
		return nil, fmt.Errorf("load database timezone: %w", err)
	}
	engine := &Engine{}
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			d := mysqlDriver.NewConfig()
			d.User = conf.User
			d.Passwd = conf.Password
			d.Net = "tcp"
			d.Addr = net.JoinHostPort(conf.Host, strconv.Itoa(conf.Port))
			d.DBName = conf.Database
			d.ParseTime = true
			d.ClientFoundRows = true
			d.Loc = location
			d.Params = map[string]string{"charset": conf.Charset}
			d.Timeout = c.ContextTimeout
			d.ReadTimeout = c.ContextTimeout
			d.WriteTimeout = c.ContextTimeout
			db, err := gorm.Open(gormMySQL.New(gormMySQL.Config{DSN: d.FormatDSN(), SkipInitializeWithVersion: true}), &gorm.Config{DisableAutomaticPing: true, TranslateError: true})
			if err != nil {
				return fmt.Errorf("open mysql database: %w", err)
			}
			sqlDB, err := db.DB()
			if err != nil {
				return fmt.Errorf("get sql database: %w", err)
			}
			sqlDB.SetConnMaxLifetime(conf.ConnMaxLifetime)
			sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
			sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
			if err := sqlDB.PingContext(ctx); err != nil {
				_ = sqlDB.Close()
				return fmt.Errorf("ping mysql database: %w", err)
			}
			engine.gorm = db
			engine.wrapLog()
			return nil
		},
		OnStop: func(context.Context) error {
			if engine.gorm == nil {
				return nil
			}
			db, err := engine.gorm.DB()
			if err != nil {
				return err
			}
			return db.Close()
		},
	})
	return engine, nil
}
func Readiness(engine *Engine) health.Check {
	return health.Check{Name: "mysql", Enabled: true, Ping: func(ctx context.Context) error {
		if engine.gorm == nil {
			return errors.New("mysql not started")
		}
		db, err := engine.gorm.DB()
		if err != nil {
			return err
		}
		return db.PingContext(ctx)
	}}
}
func (e *Engine) Connect() *gorm.DB { return e.gorm }
func (e *Engine) SetLogMode(enabled bool) {
	if !enabled {
		e.gorm.Logger = e.gorm.Logger.LogMode(LogLevelSilent)
	}
}
func (e *Engine) SetLogLevel(level LogLevel) { e.gorm.Logger = e.gorm.Logger.LogMode(level) }
func (e *Engine) wrapLog() {
	if log.Logger == nil {
		return
	}
	e.gorm.Logger = queryLogger{level: logger.Warn}
}

type transactionKey struct{}
type transactionContext struct {
	engine *Engine
	db     *gorm.DB
}

// DB 在事务回调中复用同一个连接；repository 必须使用此入口。
func (e *Engine) DB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(transactionKey{}).(transactionContext); ok {
		if tx.engine == e {
			return tx.db.WithContext(ctx)
		}
		db := e.gorm.WithContext(ctx)
		_ = db.AddError(errors.New("cross-database transaction context is unsupported"))
		return db
	}
	return e.gorm.WithContext(ctx)
}

// Transaction 只支持本数据库事务，不支持嵌套或跨数据库事务。
func (e *Engine) Transaction(ctx context.Context, fn func(context.Context) error) error {
	if ctx.Value(transactionKey{}) != nil {
		return errors.New("nested or cross-database transaction is unsupported")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if e.gorm == nil {
		return errors.New("mysql not started")
	}
	return e.gorm.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := fn(context.WithValue(ctx, transactionKey{}, transactionContext{engine: e, db: tx})); err != nil {
			return err
		}
		return ctx.Err()
	})
}
