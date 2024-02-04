package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
	Img      string    `json:"img"`
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
	var movie Movie
	params := mux.Vars(r)
	json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = params["id"]
	for index, item := range movies {
		if item.ID == params["id"] {
			movies[index] = movie
			// movies = append(movies[:index], movies[index+1:]...)
			// movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
	message := fmt.Sprintf("No such item with id %v found", params["id"])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{message})
}

func deleteMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "authentication, content-type")
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
	w.Header().Set("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Headers", "Authorization, content-type")

	json.NewEncoder(w).Encode(struct {
		Message string `json:"message"`
	}{message})
}

func main() {
	router := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Isbn: "235235", Title: "Jurassic Park", Img: "https://play-lh.googleusercontent.com/BVSejbKFir0thw8OmJKsWL-uDexGT9LDwSOcDuGE7vTC13b2JxjBHGzby7suSzvzziI", Director: &Director{Firstname: "Steven", Lastname: "Spielberg"}})

	movies = append(movies, Movie{ID: "2", Isbn: "345235", Title: "The Lost World: Jurassic Park", Img: "https://raisingchildren.net.au/__data/assets/image/0019/51355/jurassic-park-ii-the-lost-world.jpg", Director: &Director{Firstname: "Steven", Lastname: "Spielberg"}})

	router.HandleFunc("/movies", getMovies).Methods("GET")
	router.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	router.HandleFunc("/movies/create", createMovie).Methods("POST")
	router.HandleFunc("/movies/update/{id}", updateMovie).Methods("POST")
	router.HandleFunc("/movies/delete/{id}", deleteMovies).Methods("DELETE")

	cs := NewChatServer()

	router.HandleFunc("/subscribe", cs.subscribeHandler)
	router.HandleFunc("/send", cs.publishHandler).Methods("POST")

	// http.Handle("/", router)
	fmt.Println("Starting server at port 8080")

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})

	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))

}
