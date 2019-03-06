package main

import (
	"fmt"
	"geo-api/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tdoa", api.GetTdoa)

	fmt.Printf("Server started")

	log.Fatal(http.ListenAndServe(":8081", r))
}
