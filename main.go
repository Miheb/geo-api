package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/campus-iot/geo-api/api"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/tdoa", api.GetTdoa)

	fmt.Printf("Server started")

	log.Fatal(http.ListenAndServe(":8081", r))
}
