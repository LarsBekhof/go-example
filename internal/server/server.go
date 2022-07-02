package server

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go-example/internal/routes"
)

func Init() {
	host := "localhost:8080"
	router := mux.NewRouter().StrictSlash(true)
	routes.RegisterRoutes(router)
	log.Fatal(http.ListenAndServe(host, router))
}
