package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	//std necessities
	"net/http"
)

func StartWebsrv() error {
	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, r, buildHTML())
	})

	router.Get("/login/{email}/{password}", func(w http.ResponseWriter, r *http.Request) {
		email = chi.URLParam(r, "email")
		password = chi.URLParam(r, "password")
		doLogin()
	})

	router.Get("/launch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/")
		targetVersion = r.URL.Query().Get("version")
		log.Trace("/launch: version: ", targetVersion)
		doDownloadVersion()
		doDownloadAssets()
		doDownloadLibraries()
		doGameStart()
	})

	walkFunc := func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Trace("method: ", method, ", route: ", route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		return err
	}

	if err := http.ListenAndServe(":25580", router); err != nil {
		return err
	}

	return nil
}
