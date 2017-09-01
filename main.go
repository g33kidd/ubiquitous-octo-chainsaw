package main

import (
	"net/http"
)

func main() {

	http.HandleFunc("/ws", HandleWsConnection)
	http.ListenAndServe(":8080", nil)

}
