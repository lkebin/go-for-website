package api

import (
	"buhaoyong/app/api/controller"

	"github.com/gorilla/mux"
)

type Config struct {
	Domain string `mapstructure:"domain"`
	Prefix string `mapstructure:"prefix"`
}

type Repository interface{}

type repositoryImpl struct {
	config *Config
	router *mux.Router
}

func New(config *Config, router *mux.Router) Repository {
	impl := &repositoryImpl{
		config: config,
		router: router,
	}

	impl.setupRoutes()

	return impl
}

func (r *repositoryImpl) setupRoutes() {
	var router = r.router

	if r.config.Domain != "" {
		router = router.Host(r.config.Domain).Subrouter()
	}

	if r.config.Prefix != "" {
		router = router.PathPrefix(r.config.Prefix).Subrouter()
	}

	baseController := controller.New()

	helloController := controller.NewHelloController(baseController)
	router.HandleFunc("/hello", helloController.SayHello).Methods("GET")
}
