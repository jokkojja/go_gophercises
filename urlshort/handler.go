package urlshort

import (
	"net/http"
  "gopkg.in/yaml.v3"
  "os"
)

type pathToUrl struct {
    Path string `yaml:"path"`
    Url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func parseYaml(filePath string) (pathToUrls []pathToUrl) {
  
  yamlFile, err := os.ReadFile(filePath)

  if err != nil{
    panic(err)
  }

  err = yaml.Unmarshal(yamlFile, &pathToUrls)

  if err != nil{
    panic(err)
  }
  
  return
}

func buildMap(pathToUrls []pathToUrl) map[string]string {
  builtMap := make(map[string]string)
  for _, pathToUrl := range pathToUrls {
    builtMap[pathToUrl.Path] = pathToUrl.Url
  }
  return builtMap
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request){
    if redirect_url, ok := pathsToUrls[r.URL.Path]; ok{
      http.Redirect(w, r, redirect_url, http.StatusPermanentRedirect)
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
func YAMLHandler(yamlPath string, fallback http.Handler) (http.HandlerFunc, error) {
  pathsToUrls := parseYaml(yamlPath)
  builtMap := buildMap(pathsToUrls)
  return func(w http.ResponseWriter, r *http.Request){
    if redirect_url, ok := builtMap[r.URL.Path]; ok{
      http.Redirect(w, r, redirect_url, http.StatusPermanentRedirect)
    } else {
      fallback.ServeHTTP(w, r)
    }
  }, nil
}
