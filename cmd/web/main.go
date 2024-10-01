package main

import (
	"log"
	"net/http"
)

const PORT = ":4000"

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting server on port%s...", PORT)
	err := http.ListenAndServe(PORT, mux)
	log.Fatal(err)
}
