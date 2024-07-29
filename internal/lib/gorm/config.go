package gormv2

import (
	"fmt"
	"go-starter/config"
	"time"
)

type GormConfig struct {
	Alias        string        `toml:"alias" json:"alias"`
	Type         string        `toml:"type" json:"type"`
	Server       string        `toml:"server" json:"server"`
	Port         int           `toml:"port" json:"port"`
	Database     string        `toml:"database" json:"database"`
	User         string        `toml:"user" json:"user"`
	Password     string        `toml:"password" json:"password"`
	MaxIdleConns int           `toml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int           `toml:"maxOpenConns" json:"maxOpenConns"`
	Charset      string        `toml:"charset" json:"charset"`
	TimeZone     string        `toml:"timezone" json:"timezone"`
	MaxLeftTime  time.Duration `toml:"maxLeftTime" json:"maxLeftTime"`
}

func authConfig(conf config.Database) (err error) {
	if len(conf.Name) == 0 {
		conf.Name = defaultDatabase
	}

	if conf.Port == 0 {
		conf.Port = MPort
	}

	if len(conf.User) == 0 || len(conf.Password) == 0 {
		err = fmt.Errorf("User or  Password is empty")
		return
	}

	if len(conf.Host) == 0 {
		err = fmt.Errorf("server addr is empty")
		return
	}

	if len(conf.Database) == 0 {
		err = fmt.Errorf("database is empty")
		return
	}

	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = DefaultMaxIdleConns
	}

	if conf.MaxLeftTime == 0 {
		conf.MaxLeftTime = DefaultMaxLeftTime
	}

	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = DefaultMaxOpenConns
	}

	return
}
