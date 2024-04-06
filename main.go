package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"

	_ "embed"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed data/homepage.txt
var homepageText []byte

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// log requests in dev mode
	if os.Getenv("ENV") != "production" {
		r.Use(middleware.Logger)
	}

	// you can't .Mount() things on the same path
	homeRouter(r)
	pagesRouter(r)

	listenAddr := os.Getenv("LISTEN")

	slog.Info("running on http://" + listenAddr)
	http.ListenAndServe(listenAddr, r)
}

func homeRouter(r chi.Router) {
	if config.DefaultPage == "" {
		if os.Getenv("DISABLE_DEFAULT_HOMEPAGE") != "true" {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write(homepageText)
			})
		}
	} else {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Location", config.DefaultPage)
			w.WriteHeader(302)
		})
	}
}

func pagesRouter(r chi.Router) {
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		isJSON := false
		if strings.HasSuffix(path, ".json") {
			re := regexp.MustCompile(`^(.*)\.json$`)
			p := re.ReplaceAll([]byte(path), []byte("$1"))

			path = string(p)
			isJSON = true
		}

		pageData := resolvePage(path, config)
		if pageData == nil {
			// TODO: pretty page
			w.WriteHeader(404)
			w.Write([]byte("404 not found"))
		}

		if isJSON {
			j, err := json.Marshal(pageData)
			if err != nil {
				panic("failed to unmarshal json: " + err.Error())
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
			return
		}
	})
}
