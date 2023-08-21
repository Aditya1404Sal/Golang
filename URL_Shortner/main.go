package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/me":     "https://linkfree.io/Aditya1404Sal",
		"/mywork": "https://github.com/Aditya1404Sal",
	}
	mapHandler := MapHandler(pathsToUrls, mux)

	fmt.Println("Starting the server on :8080")
	// start listening to requests on port 8080 and act as an entry point for this web program
	http.ListenAndServe(":8080", mapHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	// on deault path "/" serve func hello
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Feel free to enter the key values A.K.A the Path in the search bar infront of localhost:8080")
	fmt.Fprintln(w, "1) /me")
	fmt.Fprintln(w, "2) /mywork")
}
