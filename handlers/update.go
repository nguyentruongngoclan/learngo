package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nguyentruongngoclan/learngo/data"
)

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
