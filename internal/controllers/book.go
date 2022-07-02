package controllers

import (
	"net/url"
	"net/http"
	"go-example/internal/types"
)

type BookController struct {
	AbstractController
}

func (c *BookController) convertBody(body Body) types.Book {
	book, ok := body.(types.Book)

	if ok {
		return book
	}

	c.AbstractController.failRequest(http.StatusUnprocessableEntity, "Incorrect body parameters given to create a Book")

	// This staste should never be reached since AbstractController.failRequest will panic
	return book
}

func (c *BookController) post(params map[string]int32, query url.Values, body Body) Response {
	book := c.convertBody(body)
	return Response{
		statusCode: http.StatusCreated,
		body: book,
	}
}

func InitBookController() BookController {
	return BookController{
		AbstractController{
			ImplementedMethods: []string{http.MethodPost},
		},
	}
}
