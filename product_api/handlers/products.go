// Package classification Product API.
//
// Documentation for Product API
//
//     Schemes: http
//     BasePath: /
//     Version: 1.0.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	"github.com/alonelegion/go_microservices/product_api/data"
)

// Products handler for getting and updating products
// Products является обработчиком для получения и обновления товаров
type Products struct {
	l *log.Logger
}

// NewProducts returns a new products handler with the given logger
// NewProducts возвращает обработчик новых продуктов с заданным логгером
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts returns the products from data store
// GetProducts возвращает товары из хранилища данных
func (p *Products) GetProducts(w http.ResponseWriter, _ *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch the products from the database
	// Получение товаров из хранилища
	lp := data.GetProducts()

	// serialize the list to JSON
	// сериализация списка в JSON
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle POST Product")

	prod := req.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}

	p.l.Println("Handle PUT Product", id)

	prod := req.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(req.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// validate the product
		// валидация товара
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating ", err)
			http.Error(w, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		// add the product to the context
		// добавление товара в контекст
		ctx := context.WithValue(req.Context(), KeyProduct{}, prod)
		req = req.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler
		// Вызов следующего обработчика, который может быть следующий middleware в цепи, или последним
		next.ServeHTTP(w, req)
	})
}
