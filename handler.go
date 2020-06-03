package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

type pathUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type pathFormat string

const (
	YAML pathFormat = "yaml"
	JSON pathFormat = "json"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url, ok := pathsToUrls[r.URL.Path]

		if ok {
			http.Redirect(w, r, url, http.StatusMovedPermanently)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yamlBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return unmarshalHandler(yamlBytes, fallback, YAML)
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	return unmarshalHandler(jsonBytes, fallback, JSON)
}

func unmarshalHandler(pathBytes []byte, fallback http.Handler, format pathFormat) (http.HandlerFunc, error) {
	pathUrls, err := parse(pathBytes, format)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(pathUrls)
	return MapHandler(pathsToUrls, fallback), nil
}

func parse(pathBytes []byte, format pathFormat) ([]pathUrl, error) {
	var pathUrls []pathUrl
	var err error
	if format == "json" {
		err = json.Unmarshal(pathBytes, &pathUrls)
	} else {
		err = yaml.Unmarshal(pathBytes, &pathUrls)
	}
	if err != nil {
		return nil, err
	}
	return pathUrls, nil
}

func buildMap(pathUrls []pathUrl) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, url := range pathUrls {
		pathsToUrls[url.Path] = url.URL
	}
	return pathsToUrls
}
