package main

import (
	"log"
	"net/http"
)

func HandleWsConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("handle request")
}
