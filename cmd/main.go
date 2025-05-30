package main

import (
	"context"
	"github.com/xhkzeroone/go-logger/logger"
	"time"
)

func main() {
	logger.RegisterSensitiveMessageFormater()
	err := logger.Init()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	log := logger.WithContext(ctx)
	log.Info("Ứng dụng đã khởi động")
	logger.Log.WithField("logger", "main").Info("Ứng dụng đã khởi động")
	logger.Log.WithField("logger", "main").Info("+84225898023")

	someFunc()
}

func someFunc() {
	logger.Log.WithField("logger", "main.someFunc").Warn("Gọi hàm someFunc")
}
