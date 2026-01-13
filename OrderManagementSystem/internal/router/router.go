package router

import "github.com/gorilla/mux"

func NewRounter() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	return r
}
