package controller

import (
	"buhaoyong/pkg/service/post"
	"buhaoyong/pkg/service/post/model"
	"fmt"
	"net/http"
)

type HelloController interface {
	SayHello(w http.ResponseWriter, r *http.Request)
}

type helloControllerImpl struct {
	*Controller

	postService post.Repository
}

func NewHelloController(baseController *Controller, postService post.Repository) HelloController {
	return &helloControllerImpl{Controller: baseController, postService: postService}
}

func (c *helloControllerImpl) SayHello(w http.ResponseWriter, r *http.Request) {
	_, err := c.postService.Create(r.Context(), &model.Post{Title: "标题", Content: "内容"})
	if err != nil {
		fmt.Println(fmt.Errorf("create post error: %w", err))
		c.HTML(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	c.HTML(w, http.StatusOK, "<h1>Hello</h1>")
}
