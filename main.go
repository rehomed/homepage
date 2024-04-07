package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	_ "embed"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed data/homepage.txt
var homepageText []byte

//go:embed webapp/dist
var webDist embed.FS

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	// log requests in dev mode
	if os.Getenv("ENV") != "production" {
		r.Use(middleware.Logger)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET"},
		}))
	}

	// you can't .Mount() things on the same path
	homeRouter(r)
	pagesRouter(r)

	listenAddr := os.Getenv("LISTEN")

	slog.Info("running on http://" + listenAddr)
	http.ListenAndServe(listenAddr, r)
}

func serveIndex(w http.ResponseWriter) {
	indexHtml, err := webDist.ReadFile("webapp/dist/index.html")
	fmt.Println("eee", err)
	if err != nil {
		slog.Error("failed reading index.html", "err", err)
		w.Write([]byte("internal error: " + err.Error()))
		w.WriteHeader(500)
		return
	}
	if indexHtml != nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexHtml)
		w.WriteHeader(200)
		return
	}
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
		// test as static file
		f, err := webDist.ReadFile("webapp/dist" + path)
		if f != nil && err == nil {
			switch filepath.Ext(path) {
			case ".js":
				w.Header().Set("Content-Type", "application/javascript")
			case ".css":
				w.Header().Set("Content-Type", "text/css")
			case ".woff":
				w.Header().Set("Content-Type", "font/woff")
			case ".woff2":
				w.Header().Set("Content-Type", "font/woff2")
			}
			w.Write(f)
			w.WriteHeader(200)
			return
		}

		isJSON := false
		if strings.HasSuffix(path, ".json") {
			re := regexp.MustCompile(`^(.*)\.json$`)
			p := re.ReplaceAll([]byte(path), []byte("$1"))

			path = string(p)
			isJSON = true
		}

		pageData := resolvePage(path, config)
		if pageData == nil {
			serveIndex(w)
			return
		}

		if isJSON {
			j, err := json.Marshal(pageData)
			if err != nil {
				panic("failed to unmarshal json: " + err.Error())
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
			return
		} else {
			serveIndex(w)

		}
	})
}
