package handlers

import (
	"context"
	"github.com/alonelegion/go_microservices/product_api/data"
	"net/http"
)

// MiddlewareValidateProduct validates the product in the request and calls next
// if ok
// MiddlewareValidateProduct проверяет товар в запросе и вызывает следующий,
// если все в порядке
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		prod := &data.Product{}

		err := data.FromJSON(prod, req.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			w.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, w)
			return
		}

		// validate the product
		// валидация товара
		errs := p.v.Validate(prod)
		if len(errs) != 0 {
			p.l.Println("[ERROR] validating product", errs)

			// return the validation messages as an array
			// возврат сообщения проверки в виде массива
			w.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)
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
