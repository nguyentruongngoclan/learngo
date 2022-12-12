package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/nguyentruongngoclan/learngo/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Returns a list of products
// responses:
// 	201: noContent

// DeleteProduct deletes a product from the database
func (p *Products) DeleteProduct(responseWriter http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, _ := strconv.Atoi(vars["id"])

	p.l.Println("Handle DELETE product")

	err := data.DeleteProduct(id)

	if err == data.ErrProductNotFound {
		http.Error(responseWriter, "Product not found", http.StatusNotFound)
	}
	if err != nil {
		http.Error(responseWriter, "Product not found", http.StatusInternalServerError)
		return
	}
}
