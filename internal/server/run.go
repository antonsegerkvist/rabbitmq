package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rabbitmq/api"
	"github.com/rabbitmq/web"
)

//
// Run starts the public server on port 8080.
//
func Run() {

	router := mux.NewRouter()

	api.Mount(router)
	web.Mount(router)

	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Fatal(server.ListenAndServe())

}
