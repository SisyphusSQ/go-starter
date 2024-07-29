package gormv2

import (
	"errors"
	"fmt"
	"go-starter/config"
	"go-starter/internal/lib/log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	defaultDatabase     = "mysql"
	MySQLConnTmpl       = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s"
	DefaultMaxOpenConns = 200
	DefaultMaxIdleConns = 60
	DefaultMaxLeftTime  = 300 * time.Second
	Charset             = "utf8mb4"
	MPort               = 3306
	TimeZone            = "Local"
	gormEngine          *Engine
)

type Engine struct {
	gorm *gorm.DB
}

// New 实例化新的Gorm实例
func New(config config.Config) *Engine {
	var (
		err      error
		db       *gorm.DB
		conf     = config.Database
		gormConf = &gorm.Config{}
	)

	if config.Database.Driver == "" || config.Database.Driver == "mysql" {
		err = authConfig(conf)
		if err != nil {
			panic(err)
		}
		if strings.TrimSpace(conf.Charset) == "" {
			conf.Charset = Charset
		}
		if strings.TrimSpace(conf.TimeZone) == "" {
			conf.TimeZone = TimeZone
		}

		dsn := fmt.Sprintf(MySQLConnTmpl, conf.User, conf.Password, conf.Host, conf.Port, conf.Database, conf.Charset, conf.TimeZone)
		db, err = gorm.Open(mysql.Open(dsn), gormConf)
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New(fmt.Sprintf("Not support type(%s)", conf.Driver)))
	}

	gormEngine = &Engine{db}
	gormEngine.wrapLog()
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxLifetime(conf.MaxLeftTime)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	return gormEngine
}

func (db *Engine) Connect() *gorm.DB {
	return db.gorm
}

func (db *Engine) SetLogMode(mode bool) {
	if !mode {
		db.gorm.Logger.LogMode(LogLevelSilent)
	}
}

func (db *Engine) SetLogLevel(level LogLevel) {
	db.gorm.Logger.LogMode(level)
}

func (db *Engine) wrapLog() {
	if log.Logger == nil {
		return
	}

	newLogger := logger.New(
		log.Logger,
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)
	db.gorm.Logger = newLogger
}
