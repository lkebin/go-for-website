package website

import (
	"buhaoyong/app/website/controller"
	"buhaoyong/pkg/service/post"

	"github.com/gorilla/mux"
)

type Config struct {
	Domain string `mapstructure:"domain"`
	Prefix string `mapstructure:"prefix"`
}

type Repository interface{}

type repositoryImpl struct {
	config      *Config
	router      *mux.Router
	postService post.Repository
}

func New(config *Config, router *mux.Router, postService post.Repository) Repository {
	impl := &repositoryImpl{
		config:      config,
		router:      router,
		postService: postService,
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

	helloController := controller.NewHelloController(baseController, r.postService)
	router.HandleFunc("/hello", helloController.SayHello).Methods("GET")
}
