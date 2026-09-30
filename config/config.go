package config

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

type (
	Config struct {
		Debug          bool          `mapstructure:"debug"`
		ContextTimeout time.Duration `mapstructure:"contextTimeout"`
		Server         Server        `mapstructure:"server"`
		Database       Database      `mapstructure:"database"`
		Log            Log           `mapstructure:"log"`
		Key            Key           `mapstructure:"key"`
		Cron           Cron          `mapstructure:"cron"`
		Redis          Redis         `mapstructure:"redis"`
		MongoDB        MongoDB       `mapstructure:"mongodb"`
		Prometheus     HTTPClient    `mapstructure:"prometheus"`
		Lark           Lark          `mapstructure:"lark"`
	}

	Server struct {
		Metrics           bool          `mapstructure:"metrics"`
		ReadinessTimeout  time.Duration `mapstructure:"readinessTimeout"`
		Address           string        `mapstructure:"address"`
		ReadHeaderTimeout time.Duration `mapstructure:"readHeaderTimeout"`
		ReadTimeout       time.Duration `mapstructure:"readTimeout"`
		WriteTimeout      time.Duration `mapstructure:"writeTimeout"`
		IdleTimeout       time.Duration `mapstructure:"idleTimeout"`
		ShutdownTimeout   time.Duration `mapstructure:"shutdownTimeout"`
		MaxBodyBytes      int64         `mapstructure:"maxBodyBytes"`
		CORSOrigins       []string      `mapstructure:"corsOrigins"`
	}
	Database struct {
		Enabled         bool          `mapstructure:"enabled"`
		Driver          string        `mapstructure:"driver"`
		Host            string        `mapstructure:"host"`
		Port            int           `mapstructure:"port"`
		User            string        `mapstructure:"username"`
		Password        string        `mapstructure:"password"`
		Database        string        `mapstructure:"database"`
		MaxIdleConns    int           `mapstructure:"maxIdleConns"`
		MaxOpenConns    int           `mapstructure:"maxOpenConns"`
		ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
		Charset         string        `mapstructure:"charset"`
		TimeZone        string        `mapstructure:"timeZone"`
		Name            string        `mapstructure:"name"`
	}

	Log struct {
		Output         string        `mapstructure:"output"`
		FileName       string        `mapstructure:"fileName"`
		LogLevel       zapcore.Level `mapstructure:"logLevel"`
		MaxSizeMb      int           `mapstructure:"maxSizeMB"`
		MaxBackupCount int           `mapstructure:"maxBackupCount"`
		MaxKeepDays    int           `mapstructure:"maxKeepDays"`
	}

	Key struct {
		Type  string    `mapstructure:"type"`
		Basic BasicAuth `mapstructure:"basic"`
		AK    AKAuth    `mapstructure:"ak"`
		JWT   JWTConfig `mapstructure:"jwt"`
	}

	BasicAuth struct {
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
	}

	AKAuth struct {
		AccessKey string `mapstructure:"accessKey"`
		SecretKey string `mapstructure:"secretKey"`
	}

	JWTConfig struct {
		Namespace string        `mapstructure:"namespace"`
		Secret    string        `mapstructure:"secret"`
		Expire    time.Duration `mapstructure:"expire"`
		Issuer    string        `mapstructure:"issuer"`
	}

	Cron struct {
		On bool `mapstructure:"on"`
	}

	Redis struct {
		Enabled    bool `mapstructure:"enabled"`
		PoolConfig `mapstructure:"pool"`

		Name         string        `mapstructure:"name"`
		Proto        string        `mapstructure:"proto"`
		Addr         string        `mapstructure:"addr"`
		Auth         string        `mapstructure:"auth"`
		DialTimeout  time.Duration `mapstructure:"dialTimeout"`
		ReadTimeout  time.Duration `mapstructure:"readTimeout"`
		WriteTimeout time.Duration `mapstructure:"writeTimeout"`
		DB           int           `mapstructure:"db"`
		SlowLog      time.Duration `mapstructure:"slowLog"`
	}

	PoolConfig struct {
		Active      int           `mapstructure:"active"`
		Idle        int           `mapstructure:"idle"`
		WaitTimeout time.Duration `mapstructure:"waitTimeout"`
		Wait        bool          `mapstructure:"wait"`
	}
	MongoDB struct {
		Enabled        bool          `mapstructure:"enabled"`
		URI            string        `mapstructure:"uri"`
		AuthSource     string        `mapstructure:"authSource"`
		User           string        `mapstructure:"user"`
		Password       string        `mapstructure:"password"`
		Database       string        `mapstructure:"database"`
		MaxPoolSize    uint64        `mapstructure:"maxPoolSize"`
		MinPoolSize    uint64        `mapstructure:"minPoolSize"`
		ConnectTimeout time.Duration `mapstructure:"connectTimeout"`
	}

	HTTPClient struct {
		Enabled bool          `mapstructure:"enabled"`
		URL     string        `mapstructure:"url"`
		Token   string        `mapstructure:"token"`
		Timeout time.Duration `mapstructure:"timeout"`
	}

	Lark struct {
		Enabled   bool          `mapstructure:"enabled"`
		AppID     string        `mapstructure:"appID"`
		AppSecret string        `mapstructure:"appSecret"`
		Timeout   time.Duration `mapstructure:"timeout"`
	}
)

func Load(file string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(file)
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := bindEnvironment(v); err != nil {
		return Config{}, fmt.Errorf("bind environment: %w", err)
	}
	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", file, err)
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return Config{}, fmt.Errorf("decode config %q: %w", file, err)
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, fmt.Errorf("validate config %q: %w", file, err)
	}
	return cfg, nil
}

func bindEnvironment(v *viper.Viper) error {
	keys := []string{
		"debug",
		"contextTimeout",
		"server.address",
		"server.metrics",
		"server.readinessTimeout",
		"server.readHeaderTimeout",
		"server.readTimeout",
		"server.writeTimeout",
		"server.idleTimeout",
		"server.shutdownTimeout",
		"server.maxBodyBytes",
		"server.corsOrigins",
		"key.type",
		"key.basic.user",
		"key.basic.password",
		"key.ak.accessKey",
		"key.ak.secretKey",
		"key.jwt.namespace",
		"key.jwt.secret",
		"key.jwt.expire",
		"key.jwt.issuer",

		"redis.enabled",
		"redis.addr",
		"redis.auth",
		"redis.name",
		"redis.proto",
		"redis.db",
		"redis.dialTimeout",
		"redis.readTimeout",
		"redis.writeTimeout",
		"redis.slowLog",
		"redis.pool.active",
		"redis.pool.idle",
		"redis.pool.waitTimeout",
		"redis.pool.wait",

		"cron.on",

		"log.output",
		"log.fileName",
		"log.logLevel",
		"log.maxSizeMB",
		"log.maxBackupCount",
		"log.maxKeepDays",
		"database.enabled",
		"database.driver",
		"database.host",
		"database.port",
		"database.username",
		"database.password",
		"database.database",
		"database.maxIdleConns",
		"database.maxOpenConns",
		"database.connMaxLifetime",
		"database.charset",
		"database.timeZone",
		"database.name",
		"mongodb.enabled",
		"mongodb.uri",
		"mongodb.authSource",
		"mongodb.user",
		"mongodb.password",
		"mongodb.database",
		"mongodb.maxPoolSize",
		"mongodb.minPoolSize",
		"mongodb.connectTimeout",

		"prometheus.enabled",
		"prometheus.url",
		"prometheus.token",
		"prometheus.timeout",

		"lark.enabled",
		"lark.appID",
		"lark.appSecret",
		"lark.timeout",
	}
	for _, key := range keys {
		if err := v.BindEnv(key); err != nil {
			return err
		}
	}
	return nil
}

func (c Config) Validate() error {
	if _, _, err := net.SplitHostPort(c.Server.Address); err != nil {
		return fmt.Errorf("server.address must be host:port")
	}
	if c.ContextTimeout <= 0 || c.Server.ReadinessTimeout <= 0 {
		return fmt.Errorf("context and readiness timeouts must be positive")
	}
	if c.Server.ReadHeaderTimeout <= 0 || c.Server.ReadTimeout <= 0 || c.Server.WriteTimeout <= 0 || c.Server.IdleTimeout <= 0 || c.Server.ShutdownTimeout <= 0 {
		return fmt.Errorf("server timeouts must be positive")
	}
	if c.Server.MaxBodyBytes <= 0 {
		return fmt.Errorf("server.maxBodyBytes must be positive")
	}
	for _, origin := range c.Server.CORSOrigins {
		u, err := url.Parse(origin)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || strings.ContainsAny(u.Host, "*?") || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
			return fmt.Errorf("server.corsOrigins must contain exact HTTP(S) origins")
		}
	}
	if c.Log.Output != "stdout" && c.Log.Output != "file" && c.Log.Output != "both" {
		return fmt.Errorf("log.output must be stdout,file,both")
	}
	if c.Log.Output != "stdout" && (c.Log.FileName == "" || c.Log.MaxSizeMb <= 0 || c.Log.MaxBackupCount <= 0 || c.Log.MaxKeepDays <= 0) {
		return fmt.Errorf("file logging requires a path and positive rotation limits")
	}

	if c.Redis.Enabled {
		if c.Redis.Addr == "" || c.Redis.DialTimeout <= 0 || c.Redis.ReadTimeout <= 0 || c.Redis.WriteTimeout <= 0 || c.Redis.WaitTimeout <= 0 || c.Redis.Active <= 0 || c.Redis.Idle < 0 || c.Redis.Idle > c.Redis.Active || c.Redis.DB < 0 {
			return fmt.Errorf("invalid enabled redis configuration")
		}
	}

	if c.Database.Enabled {
		if c.Database.Host == "" || c.Database.User == "" || c.Database.Database == "" || c.Database.Port <= 0 || c.Database.Port > 65535 || c.Database.MaxOpenConns <= 0 || c.Database.MaxIdleConns < 0 || c.Database.MaxIdleConns > c.Database.MaxOpenConns || c.Database.ConnMaxLifetime <= 0 {
			return fmt.Errorf("invalid enabled database configuration")
		}
	}

	if c.MongoDB.Enabled {
		if (!strings.HasPrefix(c.MongoDB.URI, "mongodb://") && !strings.HasPrefix(c.MongoDB.URI, "mongodb+srv://")) || c.MongoDB.Database == "" || c.MongoDB.ConnectTimeout <= 0 || c.MongoDB.MaxPoolSize == 0 || c.MongoDB.MinPoolSize > c.MongoDB.MaxPoolSize {
			return fmt.Errorf("invalid enabled mongodb configuration")
		}
	}

	if c.Lark.Enabled && (c.Lark.AppID == "" || c.Lark.AppSecret == "" || c.Lark.Timeout <= 0) {
		return fmt.Errorf("enabled lark requires credentials and timeout")
	}

	if c.Prometheus.Enabled {
		u, err := url.Parse(c.Prometheus.URL)
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || c.Prometheus.Timeout <= 0 {
			return fmt.Errorf("invalid enabled prometheus configuration")
		}
	}

	switch c.Key.Type {
	case "none":
		host, _, _ := net.SplitHostPort(c.Server.Address)
		ip := net.ParseIP(host)
		if !c.Debug || ip == nil || !ip.IsLoopback() {
			return fmt.Errorf("key.type none requires debug and a loopback listener")
		}
	case "basic":
		if c.Key.Basic.User == "" || c.Key.Basic.Password == "" {
			return fmt.Errorf("basic authentication requires user and password")
		}
	case "key":
		if c.Key.AK.AccessKey == "" || c.Key.AK.SecretKey == "" {
			return fmt.Errorf("key authentication requires accessKey and secretKey")
		}

	case "jwt":
		if !c.Redis.Enabled || len(c.Key.JWT.Secret) < 32 || c.Key.JWT.Expire <= 0 || c.Key.JWT.Issuer == "" || c.Key.JWT.Namespace == "" {
			return fmt.Errorf("jwt requires enabled redis, a 32-byte secret, expiry, issuer and namespace")
		}

	default:
		return fmt.Errorf("unsupported key.type %q", c.Key.Type)
	}
	return nil
}
