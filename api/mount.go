package api

import "github.com/gorilla/mux"

//
// Mount handles mounting of api endpoints to the specified router.
//
func Mount(router *mux.Router) {

	router.HandleFunc("/v1/api/upload", handleUpload).Methods("POST")
	router.HandleFunc("/v1/api/check", handleCheck).Methods("GET")
	router.HandleFunc("/v1/api/preview/{filesize:[0-9]+}/{checksum:[0-9]+}.png", handlePreview).Methods("GET")

}
