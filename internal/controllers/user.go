package controllers

import (
	"fmt"
	"net/http"
)

func getUser(request *http.Request) string {
	return "Test user"
}

func InitUserController(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
		case http.MethodGet:
			fmt.Fprintf(writer, getUser(request))
	}
}

const UserMethods = (http.MethodGet)
