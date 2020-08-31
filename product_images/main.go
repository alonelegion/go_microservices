package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alonelegion/env"
	"github.com/alonelegion/go_microservices/product_images/files"
	"github.com/alonelegion/go_microservices/product_images/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":8081", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {
	_ = env.Parse()
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
	fh := handlers.NewFiles(stor, l)
	mw := handlers.GzipHandler{}

	// create a new serve mux and register the handlers
	// создать новый serve mux и зарегестрировать обработчик
	sm := mux.NewRouter()

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadREST)
	ph.HandleFunc("/", fh.UploadMultipart)

	// get files
	// получить файлы
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle(
		"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
		http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))),
	)
	gh.Use(mw.GzipMiddleware)

	// create a new server
	// создать новый сервер
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	// запуск сервера
	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// graceful shutdown the server
	// корректное завершение работы сервера
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received
	// Блокировка до получения сигнала
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// graceful shutdown the server, waiting max 30 seconds
	// for current operations to complete
	// корректно выключить сервер, ожидание не более 30 секунд
	// для завершения текущих операций
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.Shutdown(ctx)
}
