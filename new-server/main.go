package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func PrintHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hellow")
	fmt.Fprintf(w, "Hello you've requested: %s\n", r.URL.Path)
}

func PrintHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World")
	vars := mux.Vars(r)
	title := vars["title"]
	fmt.Fprintf(w, "You have requested the following title => %s\n", title)
	// fmt.Fprintf(w, "Hello you've requested: %s\n", r.URL.Path)
}

func main() {
	// http.HandleFunc("/hello", PrintHello)
	// http.ListenAndServe(":8080", nil)

	r := mux.NewRouter()
	r.HandleFunc("/hello", PrintHello)
	r.HandleFunc("/hello/{title}", PrintHelloWorld)
	http.ListenAndServe(":8080", r)
}
