package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/alonelegion/go_microservices/product_api/handlers"
)

func main() {
	l := log.New(os.Stdout, "product_api", log.LstdFlags)

	// Create the handlers
	// Создание нового обработчика
	ph := handlers.NewProducts(l)

	// Create a new serve mux and register the handlers
	// Создание нового serve mux и регистрация обработчиков
	sm := mux.NewRouter()

	// handlers for API
	// обработчики для API
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts)
	putRouter.Use(ph.MiddlewareProductValidation)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	// create a new server
	// создание нового сервера
	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// start the server
	// запуск сервера
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	// перехват сигнала или прерывание и корректное выключение сервера
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	// Block until a signal is received
	// Блокировка до получения сигнала
	sig := <-sigChan
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds
	// for current operations to complete
	// корректно завершается работа сервера, ожидая не более 30 секунд
	// для завершения текущих операций
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	_ = s.Shutdown(tc)
}
