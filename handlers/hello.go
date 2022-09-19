package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	h.l.Println("Hello world")
	data, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(responseWriter, "Oops, something went wrong", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(responseWriter, "Hello %s\n", data)
}
