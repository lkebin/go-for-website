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
	data := make(map[string]interface{})
	data["message"] = "Hello"

	c.JSON(w, http.StatusOK, data)
}
