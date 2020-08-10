package handlers

import (
	"github.com/alonelegion/go_microservices/product_api/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

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
