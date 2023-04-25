package web

import "net/http"

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Hello World"))
}
