package handlers

import (
	"github.com/alonelegion/go_microservices/product_api/data"
	"net/http"
)

// swagger:route POST /products products createProduct
// Create a new product
//
// response:
//  200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
// Create обрабатывает запросы POST для добавления новых товаров
func (p *Products) Create(w http.ResponseWriter, req *http.Request) {
	prod := req.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("[DEBUG] Inserting product: %#v\n", prod)
	data.AddProduct(prod)
}
