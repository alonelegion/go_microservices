package handlers

import (
	"log"
	"net/http"

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
