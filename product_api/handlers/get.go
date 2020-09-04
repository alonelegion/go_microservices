package handlers

import (
	"context"
	protos "github.com/alonelegion/go_microservices/currency/protos/currency"
	"github.com/alonelegion/go_microservices/product_api/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
// ListAll обрабатывает запросы GET и возвращает все текущие товары
func (p *Products) ListAll(w http.ResponseWriter, req *http.Request) {
	p.l.Println("[DEBUG] get all records")

	w.Header().Add("Content-Type", "application/json")

	// fetch the products from the database
	// Получение товаров из базыданных
	prods := data.GetProducts()

	// serialize the list to JSON
	// сериализация списка в JSON
	err := data.ToJSON(prods, w)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
// ListSingle обрабатывает запросы GET
func (p *Products) ListSingle(w http.ResponseWriter, req *http.Request) {
	id := getProductID(req)

	p.l.Println("[DEBUG] get record id", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		w.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	// get exchange rate
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value["GBP"]),
	}
	resp, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Println("[Error] error getting new rate", err)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	prod.Price = prod.Price * resp.Rate

	err = data.ToJSON(prod, w)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}
}
