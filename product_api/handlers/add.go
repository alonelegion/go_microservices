package handlers

import (
	"github.com/alonelegion/go_microservices/product_api/data"
	"net/http"
)

func (p *Products) AddProduct(w http.ResponseWriter, req *http.Request) {
	p.l.Println("Handle POST Product")

	prod := req.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}
