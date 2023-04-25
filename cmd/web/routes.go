package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Get("/", app.HomePage)
	return mux

}

func (app *Config) Serve(port int) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app.Routes(),
	}
	app.InfoLog.Println("starting web server")
	if err := srv.ListenAndServe(); err != nil {
		app.ErrorLog.Println("error starting server", err)
	}
}
