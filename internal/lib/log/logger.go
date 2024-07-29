package log

import (
	"fmt"
	"os"

	"github.com/realcp1018/tinylog"

	"go-starter/config"
	"go-starter/vars"
)

var Logger *tinylog.TinyLogger

func New(config config.Config) {
	c := config.Log
	preCheck(c.LogLevel)

	Logger = tinylog.NewFileLogger(c.FileName, tinylog.LogLevel(c.LogLevel))
	if c.MaxSizeMb > 0 {
		Logger.SetFileConfig(c.FileName, c.MaxSizeMb, c.MaxBackupCount, c.MaxKeepDays)
	}
}

func preCheck(logLevel uint) {
	if logLevel < vars.DEBUG || logLevel > vars.FATAL {
		fmt.Printf("invalid log-level %d, should be [0,4]", logLevel)
		os.Exit(1)
	}
}
