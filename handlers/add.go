package handlers

import (
	"net/http"

	"github.com/nguyentruongngoclan/learngo/data"
)

func (p *Products) AddProduct(responseWriter http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle POST product")
	prod := request.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
}
