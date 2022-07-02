package controllers

import (
	"log"
	"fmt"
	"strconv"
	"net/http"
	"net/url"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Body interface{}

type Response struct {
	statusCode int
	body Body
}

type interfaceController interface {
	get(params map[string]int32, query url.Values) Response
	post(params map[string]int32, query url.Values, body Body) Response
	put(params map[string]int32, query url.Values, body Body) Response
	patch(params map[string]int32, query url.Values, body Body) Response
	delete(params map[string]int32, query url.Values) Response
	handleRequest(writer http.ResponseWriter, request *http.Request)
	convertBody(body Body) any
}

type AbstractController struct {
	interfaceController
	ImplementedMethods []string
	Writer http.ResponseWriter
}

func (a *AbstractController) HandleRequest(writer http.ResponseWriter, request *http.Request) {
	a.Writer = writer
	params := getPathParams(writer, request)
	query := request.URL.Query()

	var body Body
	if requestHasBody(request.Method) {
		setJsonBody(writer, request, &body)
	}

	var response Response

	switch request.Method {
		case http.MethodGet:
			response = a.get(params, query)
		case http.MethodPost:
			response = a.post(params, query, body)
		case http.MethodPut:
			response = a.put(params, query, body)
		case http.MethodPatch:
			response = a.patch(params, query, body)
		case http.MethodDelete:
			response = a.delete(params, query)

		setHeaders(writer)
		writeResponse(writer, response)
	}
}

func (a *AbstractController) failRequest(status int, message string) {
	if a.Writer == nil {
		log.Panicln("AbstractController.failRequest method called before AbstractController.handleRequest")
	}

	a.Writer.WriteHeader(status)
	fmt.Fprintln(a.Writer, http.StatusText(status));
}

func getPathParams(writer http.ResponseWriter, request *http.Request) map[string]int32 {
	pathVars := mux.Vars(request)

	parsedVars := make(map[string]int32)
	for key, pathVar := range pathVars {
		id, err := strconv.ParseInt(pathVar, 10, 32)

		if err != nil {
			failRequest(writer, http.StatusNotAcceptable)
		}

		parsedVars[key] = int32(id)
	}

	return parsedVars
}

func requestHasBody(method string) bool {
	switch method {
		case
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch:
			return true
	}
	return false

}

func setJsonBody(writer http.ResponseWriter, request *http.Request, instance *Body) {
	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&instance)

	if err != nil {
		failRequest(writer, http.StatusUnprocessableEntity)
	}
}

func checkRequestContentType(writer http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Type") != "application/json" {
		failRequest(writer, http.StatusUnsupportedMediaType)
	}
}

func setHeaders(writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")
}

func writeResponse(writer http.ResponseWriter, response Response) {
	writer.WriteHeader(response.statusCode)

	encoder := json.NewEncoder(writer)
	err := encoder.Encode([]byte(fmt.Sprintf("%v", response.body)))

	if err != nil {
		log.Default().Println(err)
		failRequest(writer, http.StatusInternalServerError)
	}
}

func failRequest(writer http.ResponseWriter, status int) {
	writer.WriteHeader(status)
	fmt.Fprintln(writer, http.StatusText(status));
}
