package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nguyentruongngoclan/learngo/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(responseWriter http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle GET products")
	productList := data.GetProducts()
	err := productList.ToJSON(responseWriter)
	if err != nil {
		http.Error(responseWriter, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(responseWriter http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle POST product")
	prod := request.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}

func (p *Products) UpdateProducts(responseWriter http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle PUT product")
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(responseWriter, "Unable to convert id to number", http.StatusBadRequest)
		return
	}
	prod := request.Context().Value(KeyProduct{}).(*data.Product)
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(responseWriter, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(responseWriter, "Product not found", http.StatusInternalServerError)
		return
	}
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
