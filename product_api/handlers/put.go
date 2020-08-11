package handlers

import (
	"github.com/alonelegion/go_microservices/product_api/data"
	"net/http"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
// Update обрабатывает запросы PUT для обновления товаров
func (p *Products) Update(w http.ResponseWriter, req *http.Request) {
	// fetch the product from the context
	// получение товара из контекста
	prod := req.Context().Value(KeyProduct{}).(data.Product)
	p.l.Println("[DEBUG] updating record id", prod.ID)

	err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)

		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
