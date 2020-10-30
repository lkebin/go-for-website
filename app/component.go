package app

import "github.com/gorilla/mux"

type Component struct {
	Router *mux.Router
}

func New(
	router *mux.Router,
) *Component {
	return &Component{
		Router: router,
	}
}
