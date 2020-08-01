package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	//std necessities
	"net/http"
	"net/url"
)

func StartWebsrv() error {
	router := chi.NewRouter()
	router.Use(
		middleware.RedirectSlashes,
		middleware.Recoverer,
		//middleware.Logger,
	)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.HTML(w, r, htmlLogin(nil))
	})

	router.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		email, _ = url.QueryUnescape(r.URL.Query().Get("mcemail"))
		password, _ = url.QueryUnescape(r.URL.Query().Get("mcpwd"))

		doLoadToken()
		if err := doLogin(); err != nil {
			log.Error(err)
			render.HTML(w, r, htmlLogin(err))
			return
		}

		doFetchVersions()

		render.HTML(w, r, htmlLauncher())
	})

	router.Get("/launch", func(w http.ResponseWriter, r *http.Request) {
		targetVersion = r.URL.Query().Get("version")
		server = r.URL.Query().Get("server")

		log.Trace("/launch: version: ", targetVersion)
		render.PlainText(w, r, "Starting Minecraft "+targetVersion+"...")

		go func() {
			doDownloadVersion()
			doDownloadAssets()
			doDownloadLibraries()
			doGameStart()
		}()
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
