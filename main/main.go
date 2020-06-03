package main

import (
	"fmt"
	"net/http"

	"github.com/vncsb/urlshort"
)

func main() {
	mux := defaultMux()
	err := urlshort.SetupDB(map[string]string{
		"/bolt":         "https://github.com/boltdb/bolt",
		"/bolt-buckets": "https://github.com/boltdb/bolt#using-buckets",
	})
	if err != nil {
		panic(err)
	}

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	json := `[
		{
			"path": "/testeseufdp",
			"url": "https://google.com/"
		},
		{
			"path": "/meugithub",
			"url": "https://github.com/vncsb"
		}
	]`

	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	dbHandler := urlshort.DBHandler(jsonHandler)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", dbHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
