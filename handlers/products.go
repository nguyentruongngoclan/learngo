package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/nguyentruongngoclan/learngo/data"
)

// A list of products returns in the response
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:response noContent
type noContentWrapper struct {
}

// swagger:parameters deleteProduct
type productIdParameterWrapper struct {
	// The id of the product to delete from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

// Prodcuts is a http.Handler
type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

type KeyProduct struct{}

func (p Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(request.Body)
		if err != nil {
			http.Error(responseWriter, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}
		// Validate the product
		err = prod.Validate()
		if err != nil {
			http.Error(
				responseWriter,
				fmt.Sprintf("Error validate prodcut: %s", err),
				http.StatusBadRequest,
			)
			return
		}
		// Add the product to the context
		ctx := context.WithValue(request.Context(), KeyProduct{}, prod)
		request = request.WithContext(ctx)
		next.ServeHTTP(responseWriter, request)
	})
}
