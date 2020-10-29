package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) HTML(w http.ResponseWriter, statusCode int, content string) error {
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err := fmt.Fprint(w, content)
	return err
}

func (c *Controller) JSONString(w http.ResponseWriter, statusCode int, s string) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err := fmt.Fprint(w, s)
	return err
}

func (c *Controller) JSON(w http.ResponseWriter, statusCode int, i interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(i)
}
