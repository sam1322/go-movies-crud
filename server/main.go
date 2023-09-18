package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello")
	jsonData, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
	// json.NewEncoder(w).Encode(movies)
}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	message := fmt.Sprintf("No such item with id %v found", params["id"])
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{message})
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = params["id"]
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
	message := fmt.Sprintf("No such item with id %v found", params["id"])
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{message})
}

func deleteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	flag := false
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			flag = true
			break
		}
	}
	var message string
	if flag == true {
		message = fmt.Sprintf("Successfull deletion of id %v", params["id"])
	} else {
		message = fmt.Sprintf("No such id found")
	}
	// fmt.Printf("message %v\n", message)
	log.Printf("message %v\n", message)
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{message})
}
func main() {
	router := mux.NewRouter()
	movies = append(movies, Movie{ID: "1", Isbn: "235235", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "345235", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Jobs"}})
	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies", createMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", updateMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", deleteMovies).Methods("DELETE")
	http.Handle("/", router)

	fmt.Println("Starting server at port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
