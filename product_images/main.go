package main

import (
	"github.com/alonelegion/env"
	hclog "github.com/hashicorp/go-hclog"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":8080", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./image_store", "Base path to save images")

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product_images",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	// Create a logger for server from the default logger
	// Создание логгера для сервера из дефолтного логгера
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// Create the storage class, use local storage
	// max filesize 5MB
	// Создание класса локального хранилища
	// максимальный размер одного файла 5MB
	stor, err :=
}
