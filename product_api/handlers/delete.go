package handlers

import (
	"github.com/alonelegion/go_microservices/product_api/data"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// swagger:route DELETE /products/{id} products deleteProducts
// Return a list of products from the database
// responses:
//	201: noContent

// DeleteProduct deletes a product from database
func (p *Products) DeleteProduct(w http.ResponseWriter, req *http.Request) {
	// this will always convert because of the router
	// это всегда будет конвертироваться из-за маршрутизатора
	vars := mux.Vars(req)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle DELETE Product", id)

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}
