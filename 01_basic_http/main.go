package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandlerFunc)
	http.ListenAndServe(":3000", nil)
}

func indexHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
