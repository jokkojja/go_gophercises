package main

import (
	"fmt"
	"net/http"
  "urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
  maphandler := urlshort.MapHandler(pathsToUrls, mux)
  
	// Build the YAMLHandler using the mapHandler as the
	// fallback
// 	yaml := `
// - path: /urlshort
//   url: https://github.com/gophercises/urlshort
// - path: /urlshort-final
//   url: https://github.com/gophercises/urlshort/tree/solution
// `
// 	yamlhandler, err := urlshort.yamlhandler([]byte(yaml), maphandler)
// 	if err != nil {
// 		panic(err)
// 	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", maphandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
