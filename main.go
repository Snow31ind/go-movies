package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Director struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

var movies []Movie

const (
	MAXIMUM_ID = 1000000000
)

func getRandomId() string {
	return strconv.Itoa(rand.Intn(MAXIMUM_ID))
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	movie.ID = getRandomId()
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movies)
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func updateMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, movie := range movies {
		if movie.ID == params["id"] {
			var newMovie Movie
			if err := json.NewDecoder(r.Body).Decode(&newMovie); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			newMovie.ID = params["id"]
			movies[i] = newMovie
			json.NewEncoder(w).Encode(newMovie)
			return
		}
	}
}

func deleteMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:i], movies[i+1:]...)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}

	http.Error(w, "id not found", http.StatusBadRequest)
}

func seed() {
	movies = append(movies, Movie{
		ID:    "1",
		Isbn:  "150101",
		Title: "The Uncanny Counter",
		Director: &Director{
			FirstName: "You",
			LastName:  "Sun-dong",
		},
	})
	movies = append(movies, Movie{
		ID:    "2",
		Isbn:  "290902",
		Title: "Burn the House Down",
		Director: &Director{
			FirstName: "Mei",
			LastName:  "Negano",
		},
	})
}

func main() {
	const PORT = 8080
	r := mux.NewRouter()
	seed()

	moviesRouter := r.PathPrefix("/movies").Subrouter()
	{
		moviesRouter.HandleFunc("", getMovies).Methods("GET")
		moviesRouter.HandleFunc("", createMovie).Methods("POST")
		moviesRouter.HandleFunc("/{id}", getMovieById).Methods("GET")
		moviesRouter.HandleFunc("/{id}", updateMovieById).Methods("PUT")
		moviesRouter.HandleFunc("/{id}", deleteMovieById).Methods("DELETE")
	}

	fmt.Printf("Starting server at port %v\n", PORT)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", PORT), r); err != nil {
		panic(err)
	}
}
