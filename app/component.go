package app

import (
	"buhaoyong/pkg/db"
	"github.com/gorilla/mux"
)

type Component struct {
	Router *mux.Router
	DB     *db.DB
}

func New(
	router *mux.Router,
	dbc *db.DB,
) *Component {
	return &Component{
		Router: router,
		DB:     dbc,
	}
}
