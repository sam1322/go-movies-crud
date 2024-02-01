// package main

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// func PrintHello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Hellow")
// 	fmt.Fprintf(w, "Hello you've requested: %s\n", r.URL.Path)
// }

// func PrintHelloWorld(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Hello World")
// 	vars := mux.Vars(r)
// 	title := vars["title"]
// 	fmt.Fprintf(w, "You have requested the following title => %s\n", title)
// 	// fmt.Fprintf(w, "Hello you've requested: %s\n", r.URL.Path)
// }

// func main() {
// 	// http.HandleFunc("/hello", PrintHello)
// 	// http.ListenAndServe(":8080", nil)

//		r := mux.NewRouter()
//		r.HandleFunc("/hello", PrintHello)
//		r.HandleFunc("/hello/{title}", PrintHelloWorld)
//		http.ListenAndServe(":8080", r)
//	}
package main

import (
	"fmt"
	"time"
)

func main() {
	pingChan := make(chan int, 1)
	pongChan := make(chan int, 1)

	go ping(pingChan, pongChan)
	go pong(pongChan, pingChan)

	pingChan <- 1

	select {}
}

func ping(pingChan <-chan int, pongChan chan<- int) {
	for {
		ball := <-pingChan

		fmt.Println("Ping", ball)
		time.Sleep(1 * time.Second)

		pongChan <- ball + 1
	}
}

func pong(pongChan <-chan int, pingChan chan<- int) {
	for {
		ball := <-pongChan

		fmt.Println("Pong", ball)
		time.Sleep(1 * time.Second)

		pingChan <- ball + 1
	}
}
