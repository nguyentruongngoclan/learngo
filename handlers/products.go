package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/nguyentruongngoclan/learngo/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	// Handle get
	if request.Method == http.MethodGet {
		p.getProducts(responseWriter, request)
		return
	}
	// Handle update
	if request.Method == http.MethodPost {
		p.addProduct(responseWriter, request)
		return
	}
	// Handle update
	if request.Method == http.MethodPut {
		p.l.Println("PUT")
		// Request validation
		// Expect the id in the URI
		regex := regexp.MustCompile(`/([0-9]+)`)
		group := regex.FindAllStringSubmatch(request.URL.Path, -1)
		if len(group) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(responseWriter, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(group[0]) != 2 {
			p.l.Println("Invalid URI more than one captured group")
			http.Error(responseWriter, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := group[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(responseWriter, "Unable to parse id", http.StatusInternalServerError)
		}
		p.updateProducts(id, responseWriter, request)
		p.l.Println("Got id", id)

		return
	}
	// Catch all
	responseWriter.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(responseWriter http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle GET products")
	productList := data.GetProducts()
	err := productList.ToJSON(responseWriter)
	if err != nil {
		http.Error(responseWriter, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(responseWriter http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle POST product")
	prod := &data.Product{}
	err := prod.FromJSON(request.Body)
	if err != nil {
		http.Error(responseWriter, "Unable to unmarshal json", http.StatusBadRequest)
	}
	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, responseWriter http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle PUT product")
	p.l.Println("Handle POST product")
	prod := &data.Product{}
	err := prod.FromJSON(request.Body)
	if err != nil {
		http.Error(responseWriter, "Unable to unmarshal json", http.StatusBadRequest)
	}
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
