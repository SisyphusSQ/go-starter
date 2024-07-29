package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

var (
	configFile = "config/config.yml"
	configType = "yml"
)

type (
	Config struct {
		Debug          bool     `mapstructure:"debug"`
		ContextTimeout int      `mapstructure:"contextTimeout"`
		Server         Server   `mapstructure:"server"`
		Database       Database `mapstructure:"database"`
		Log            Log      `mapstructure:"log"`
		Key            Key      `mapstructure:"key"`
		Cron           Cron     `mapstructure:"cron"`
	}

	Server struct {
		Address string `mapstructure:"address"`
	}

	Database struct {
		Driver       string        `mapstructure:"driver"`
		Host         string        `mapstructure:"host"`
		Port         int           `mapstructure:"port"`
		User         string        `mapstructure:"username"`
		Password     string        `mapstructure:"password"`
		Database     string        `mapstructure:"database"`
		MaxIdleConns int           `mapstructure:"maxIdleConns"`
		MaxLeftTime  time.Duration `mapstructure:"maxLeftTime"`
		MaxOpenConns int           `mapstructure:"maxOpenConns"`
		Charset      string        `mapstructure:"charset"`
		TimeZone     string        `mapstructure:"timeZone"`
		Name         string        `mapstructure:"name"`
	}

	Log struct {
		FileName       string `toml:"fileName"`
		LogLevel       uint   `toml:"logLevel"`
		MaxSizeMb      int    `toml:"maxSizeMB"`
		MaxBackupCount int    `toml:"maxBackupCount"`
		MaxKeepDays    int    `toml:"maxKeepDays"`
	}

	Key struct {
		AccessKey string `mapstructure:"accessKey"`
		SecretKey string `mapstructure:"secretKey"`
	}

	Cron struct {
		On bool `mapstructure:"on"`
	}
)

func NewConfig() Config {
	conf := &Config{}
	err := viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable decode into config struct, %v", err)
	}
	return *conf
}

func InitConfig() {
	viper.SetConfigType(configType)
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()

	if err != nil {
		fmt.Println(err.Error())
	}
}

func SetConfigFile(file string) {
	configFile = file
}
