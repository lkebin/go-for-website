package controller

import (
	"net/http"
)

type HelloController interface {
	SayHello(w http.ResponseWriter, r *http.Request)
}

type helloControllerImpl struct {
	*Controller
}

func NewHelloController(baseController *Controller) HelloController {
	return &helloControllerImpl{Controller: baseController}
}

func (c *helloControllerImpl) SayHello(w http.ResponseWriter, r *http.Request) {
	c.HTML(w, http.StatusOK, "<h1>Hello</h1>")
}
