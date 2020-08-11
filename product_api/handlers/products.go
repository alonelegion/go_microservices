package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/alonelegion/go_microservices/product_api/data"
)

// KeyProduct is a key used for the Product object in the context
// KeyProduct - это ключ, используемый для объекта Product в контексте
type KeyProduct struct{}

// Products handler for getting and updating products
// Products является обработчиком для получения и обновления товаров
type Products struct {
	l *log.Logger
	v *data.Validation
}

// NewProducts returns a new products handler with the given logger
// NewProducts возвращает обработчик новых товаров с заданным логгером
func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

// ErrInvalidProductPath is an error message when the product path is not valid
// ErrInvalidProductPath - это сообщение об ошибке, если путь к товару  недействителен.
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

// GenericError is a generic error message returned by a server
// GenericError - это общее сообщение об ошибке, возвращаемое сервером
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
// ValidationError - это набор сообщений об ошибках валидации
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getProductID returns the product ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number

// getProductID возвращает идентификатор товара из URL
// Вызывается паника, если не удается преобразовать id в целое число
// этого никогда не должно происходить, поскольку маршрутизатор гарантирует, что
// число валидное
func getProductID(r *http.Request) int {
	// parse the product id from the url
	// парсинг идентификатора товара из URL-адреса
	vars := mux.Vars(r)

	// convert the id into an integer and return
	// ковертирование идентификатора в целое число и возврат его
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}
