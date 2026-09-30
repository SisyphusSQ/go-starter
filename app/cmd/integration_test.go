//go:build integration

package cmd

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"uuid"

	"fmt"
	"github.com/labstack/echo/v5"

	"github.com/SisyphusSQ/go-starter/v2/utils"

	"github.com/SisyphusSQ/go-starter/v2/config"
	apphttp "github.com/SisyphusSQ/go-starter/v2/internal/http"
	"github.com/SisyphusSQ/go-starter/v2/internal/lib/log"
	"go.uber.org/fx"

	gormv2 "github.com/SisyphusSQ/go-starter/v2/internal/lib/gorm"

	"github.com/SisyphusSQ/go-starter/v2/internal/lib/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/SisyphusSQ/go-starter/v2/internal/lib/redis"
)

type integrationDependencies struct {
	fx.In
	Server *apphttp.Server
	DB     *gormv2.Engine  `optional:"true"`
	Mongo  *mongodb.Client `optional:"true"`
	Redis  *redis.Client   `optional:"true"`
}

func TestIsolatedComponents(t *testing.T) {
	if os.Getenv("APP_TEST_INTEGRATION") != "1" {
		t.Skip("isolated integration environment not enabled")
	}
	file := os.Getenv("APP_TEST_CONFIG")
	if !filepath.IsAbs(file) {
		t.Fatal("APP_TEST_CONFIG must be an absolute isolated config path")
	}
	cfg, err := config.Load(file)
	if err != nil {
		t.Fatal(err)
	}
	cfg.Server.Address = "127.0.0.1:0"
	if err := log.New(cfg); err != nil {
		t.Fatal(err)
	}
	defer log.Sync()
	var deps integrationDependencies
	app := fx.New(fx.NopLogger, inject(cfg), fx.Populate(&deps))
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()
	if err := app.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := app.Stop(closeCtx); err != nil {
			t.Error(err)
		}
	})
	recorder := httptest.NewRecorder()
	deps.Server.Echo.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/ready", nil))
	if recorder.Code != 200 {
		t.Fatalf("ready=%d", recorder.Code)
	}
	name := "starter_probe_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	tested := false

	if deps.DB != nil {
		tested = true
		if err := deps.DB.DB(ctx).Exec("CREATE TABLE " + name + " (id BIGINT PRIMARY KEY)").Error; err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := deps.DB.DB(closeCtx).Exec("DROP TABLE " + name).Error; err != nil {
				t.Error(err)
			}
		})
		rollback := errors.New("intentional rollback")
		err := deps.DB.Transaction(ctx, func(txctx context.Context) error {
			if err := deps.DB.DB(txctx).Exec("INSERT INTO "+name+" (id) VALUES (?)", 1).Error; err != nil {
				return err
			}
			return rollback
		})
		if !errors.Is(err, rollback) {
			t.Fatalf("rollback=%v", err)
		}
		var count int64
		if err := deps.DB.DB(ctx).Table(name).Count(&count).Error; err != nil || count != 0 {
			t.Fatalf("rollback left %d rows: %v", count, err)
		}
		if err := deps.DB.DB(ctx).Exec("INSERT INTO "+name+" (id) VALUES (?)", 1).Error; err != nil {
			t.Fatal(err)
		}
		unchanged := deps.DB.DB(ctx).Exec("UPDATE "+name+" SET id = ? WHERE id = ?", 1, 1)
		if unchanged.Error != nil || unchanged.RowsAffected != 1 {
			t.Fatalf("unchanged existing row must match once: rows=%d err=%v", unchanged.RowsAffected, unchanged.Error)
		}

	}

	if deps.Mongo != nil {
		tested = true
		collection := deps.Mongo.Collection(name)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := collection.Drop(ctx); err != nil {
				t.Error(err)
			}
		})
		inserted, err := collection.InsertOne(ctx, bson.M{"value": "probe"})
		if err != nil {
			t.Fatal(err)
		}
		var result bson.M
		if err := collection.FindOne(ctx, bson.M{"_id": inserted.InsertedID}).Decode(&result); err != nil || result["value"] != "probe" {
			t.Fatalf("mongo read failed: %v", err)
		}

	}

	if deps.Redis != nil {
		tested = true
		if err := deps.Redis.Set(ctx, name, "probe", time.Minute).Err(); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := deps.Redis.Del(ctx, name).Err(); err != nil {
				t.Error(err)
			}
		})
		value, err := deps.Redis.Get(ctx, name).Result()
		if err != nil || value != "probe" {
			t.Fatalf("redis read failed: %v", err)
		}

		authcfg := cfg
		authcfg.Key.Type = "jwt"
		authcfg.Key.JWT = config.JWTConfig{Secret: uuid.New().String(), Expire: time.Minute, Issuer: "isolated-test", Namespace: name}
		token, err := utils.GenerateToken(1, "probe@example.com", authcfg.Key.JWT)
		if err != nil {
			t.Fatal(err)
		}
		key := fmt.Sprintf("%s:jwt:user:1", name)
		t.Cleanup(func() {
			closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := deps.Redis.Del(closeCtx, key).Err(); err != nil {
				t.Error(err)
			}
		})
		if err := deps.Redis.Set(ctx, key, token, time.Minute).Err(); err != nil {
			t.Fatal(err)
		}
		server := echo.New()
		middleware := apphttp.InitMiddleware(authcfg, deps.Redis)
		server.HTTPErrorHandler = middleware.ErrorHandler
		server.Use(middleware.Logger, middleware.Auth)
		server.GET("/private", func(c *echo.Context) error { return c.NoContent(204) })
		authenticated := func(want int) {
			t.Helper()
			req := httptest.NewRequest(http.MethodGet, "/private", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)
			if rec.Code != want {
				t.Fatalf("JWT status=%d want=%d", rec.Code, want)
			}
		}
		authenticated(204)
		if err := deps.Redis.Del(ctx, key).Err(); err != nil {
			t.Fatal(err)
		}
		authenticated(401)

	}

	_ = name
	_ = errors.New
	if !tested {
		t.Fatal("no isolated database or Redis enabled")
	}
}
