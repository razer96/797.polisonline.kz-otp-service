package logger

import (
	"github.com/hashicorp/go-hclog"
)

var Logger = hclog.New(&hclog.LoggerOptions{
	Name:  "insurance-otp-service",
	Level: hclog.LevelFromString("INFO"),
})

func GetLogger() hclog.Logger {
	return Logger
}
