package main

import (
	"net/http"
)

// MapHandler returns a http.HandlerFunc and takes in 2 parameters
// a map of key value pairs of [string] and string
// implements an anonymous function which attempts to map the paths to their corresponding values(URL)
// If the path is not provided in map , fallback will be called
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}
