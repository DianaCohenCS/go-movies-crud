/*
Learn Go Programming by Building Projects, by freeCodeCamp.org
Project #2: Build a CRUD API (w/o DB)

Within go-movies-crud directory, run the following commands:

$ go mod init github.com/DianaCohenCS/go-movies-crud
// this initializes the module to manage all its dependencies

$ go get "github.com/gorilla/mux"
// download this package to be used in our project
// its dependency will be registered in our "go.mod" file

To test our API one can use POSTMAN, we use Thunder Client extension for VSCode instead
*/

package main

import (
	"encoding/json" // encode data when sending to POSTMAN
	"fmt"
	"log"
	"math/rand"
	"net/http" // create a server
	"strconv"  // convert integer into string

	"github.com/gorilla/mux"
)

// no DB, just use structs and slices

// every movie has one director
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"` // a pointer to a director struct (object)
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var movies []Movie

func main() {
	r := mux.NewRouter()

	// populate some data for movies
	movies = append(movies, Movie{"1", "438227", "Movie One", &Director{"John", "Doe"}})
	movies = append(movies, Movie{"2", "45455", "Movie Two", &Director{"Steve", "Smith"}})

	// create the routes and map to their functions
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	// start the web server having the mux router
	fmt.Println("Starting server at port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}

// pass a pointer to request sent via POSTMAN
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			// append all the data except for index
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	// optional to return all the remaining movies
	json.NewEncoder(w).Encode(movies)
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
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	// get the movie from the body of the POST request
	_ = json.NewDecoder(r.Body).Decode(&movie)
	// generate new random ID for the movie
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	// return the movie just created
	json.NewEncoder(w).Encode(movie)
}

// two-step update: delete and then insert
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// delete a movie
	for index, item := range movies {
		if item.ID == params["id"] {
			// append all the data except for index
			movies = append(movies[:index], movies[index+1:]...)
			// insert back with new data
			var movie Movie
			// get the movie from the body of the PUT request
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			// return the movie just created
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}
