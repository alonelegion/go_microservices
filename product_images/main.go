package main

import (
	"github.com/alonelegion/env"
	"github.com/alonelegion/go_microservices/product_images/files"
	hclog "github.com/hashicorp/go-hclog"
	"os"
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
	stor, err := files.NewLocal(*basePath, 1024*1000*5)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// create the handlers
	// создание обработчиков
	fh :=
}
