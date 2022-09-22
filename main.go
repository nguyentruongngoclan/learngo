package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/nguyentruongngoclan/learngo/handlers"
)

func main() {
	logger := log.New(os.Stdout, "product-api", log.LstdFlags)
	productHandler := handlers.NewProduct(logger)

	serveMux := mux.NewRouter()
	// getRouter will serve only GET methods
	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", productHandler.GetProducts)
	// putRouter will serve only PUT methods
	putRouter := serveMux.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", productHandler.UpdateProducts)
	putRouter.Use(productHandler.MiddlewareValidateProduct)
	// postRouter will serve only POST methods
	postRouter := serveMux.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", productHandler.AddProduct)
	postRouter.Use(productHandler.MiddlewareValidateProduct)

	server := &http.Server{
		Addr:         ":9090",
		Handler:      serveMux,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	sig := <-signalChannel
	logger.Println("Received terminate, graceful shutdown", sig)

	timeoutContext, _ := context.WithTimeout(context.Background(), 30*time.Second)

	server.Shutdown(timeoutContext)
}
