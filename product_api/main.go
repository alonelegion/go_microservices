package main

import (
	"context"
	"github.com/alonelegion/env"
	protos "github.com/alonelegion/go_microservices/currency/protos/currency"
	"github.com/alonelegion/go_microservices/product_api/data"
	"github.com/go-openapi/runtime/middleware"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	go_handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/alonelegion/go_microservices/product_api/handlers"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":8080", "Bind address for the server")

func main() {
	env.Parse()

	l := log.New(os.Stdout, "product_api", log.LstdFlags)
	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:8082", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create client
	cc := protos.NewCurrencyClient(conn)

	// Create the handlers
	// Создание нового обработчика
	ph := handlers.NewProducts(l, v, cc)

	// Create a new serve mux and register the handlers
	// Создание нового serve mux и регистрация обработчиков
	sm := mux.NewRouter()

	// handlers for API
	// обработчики для API
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", ph.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", ph.Update)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", ph.Create)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	// Обработчик для документации
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	ch := go_handlers.CORS(go_handlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// create a new server
	// создание нового сервера
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	// запуск сервера
	go func() {
		l.Println("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	// перехват сигнала или прерывание и корректное выключение сервера
	sigChan := make(chan os.Signal, 1)
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
