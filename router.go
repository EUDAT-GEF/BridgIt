package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	// Serving static content (Weblicht service results)
	fs := http.FileServer(http.Dir(Config.StaticContent))
	staticHandler := Logger(fs, "Serving static content")

	router.
		Methods("GET").
		PathPrefix("/static").
		Name("static").
		Handler(http.StripPrefix("/static", staticHandler))
	return router
}
