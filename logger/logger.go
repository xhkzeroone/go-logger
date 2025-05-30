package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

var registeredMessageFormater MessageFormater = &DefaultMessageFormater{}
var registeredDefaultFunctionNameFormatter FunctionNameFormatter = &DefaultFunctionNameFormatter{}

func RegisterMessageFormater(m MessageFormater) {
	registeredMessageFormater = m
}
func RegisterFunctionNameFormatter(m FunctionNameFormatter) {
	registeredDefaultFunctionNameFormatter = m
}

func GetMessageFormater() MessageFormater {
	return registeredMessageFormater
}

func GetFunctionNameFormatter() FunctionNameFormatter {
	return registeredDefaultFunctionNameFormatter
}

var Log *logrus.Logger
var userOnce sync.Once

func Init() error {
	userOnce.Do(func() {
		dir, _ := os.Getwd()
		cfg, err := LoadLogConfig(filepath.Join(dir, "/log-config.xml"))
		if err != nil {
			cfg = &LogConfig{
				TimestampFormat: "2006-01-02 15:04:05",
				Pattern:         "%timestamp% | %level% | %requestId% | %file%:%line% | %function% | %message%",
				Level:           "info",
			}
		}

		level, err := logrus.ParseLevel(cfg.Level)
		if err != nil {
			level = logrus.InfoLevel
		}
		Log = logrus.New()
		Log.SetReportCaller(true)
		Log.SetLevel(level)
		Log.SetFormatter(&DynamicFormatter{
			Pattern:               cfg.Pattern,
			TimestampFormat:       cfg.TimestampFormat,
			MsgFormatter:          GetMessageFormater(),
			FunctionNameFormatter: GetFunctionNameFormatter(),
		})
	})
	return nil
}
