package main

import (
	"context"
	"fmt"
	"github.com/GoServer/newsapi"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var headerMap = map[string]string {"Content-Type": "application/json"}

func constructResponse(httpResponseWriter http.ResponseWriter, headerMap map[string]string, statusCode int, content []byte) {
	for key, val := range headerMap {
		httpResponseWriter.Header().Set(key, val)
	}
	httpResponseWriter.WriteHeader(statusCode)
	httpResponseWriter.Write(content)
}

func get(httpResponseWriter http.ResponseWriter, request *http.Request) {
	constructResponse(httpResponseWriter, headerMap, http.StatusOK, []byte(`{"message": "get called"}`))
}

func post(httpResponseWriter http.ResponseWriter, request *http.Request) {
	constructResponse(httpResponseWriter, headerMap, http.StatusCreated, []byte(`{"message": "post called"}`))
}

func put(httpResponseWriter http.ResponseWriter, request *http.Request) {
	constructResponse(httpResponseWriter, headerMap, http.StatusAccepted, []byte(`{"message": "put called"}`))
}

func delete(httpResponseWriter http.ResponseWriter, request *http.Request) {
	constructResponse(httpResponseWriter, headerMap, http.StatusOK, []byte(`{"message": "delete called"}`))
}

func notFound(httpResponseWriter http.ResponseWriter, request *http.Request) {
	constructResponse(httpResponseWriter, headerMap, http.StatusNotFound, []byte(`{"message": "not found"}`))
}

func params(httpResponseWriter http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusOK)

	if _, ok := pathParams["keyword"]; ok {
		c := newsapi.NewClient("6c5c888290f647818122022f271a88f0", newsapi.WithHTTPClient(http.DefaultClient))

		sources, err := c.GetSources(context.Background(), nil)

		if err != nil {
			constructResponse(httpResponseWriter, headerMap, http.StatusInternalServerError, []byte(`{"message": "something is wrong"}`))
			panic(err)
		}

		for _, s := range sources.Sources {
			httpResponseWriter.Write([]byte(fmt.Sprintf(s.Description)))
		}
	}
}

func main() {
	request := mux.NewRouter()

	apiTest := request.PathPrefix("/api/test").Subrouter()
	apiTest.HandleFunc("", get).Methods(http.MethodGet)
	apiTest.HandleFunc("", post).Methods(http.MethodPost)
	apiTest.HandleFunc("", put).Methods(http.MethodPut)
	apiTest.HandleFunc("", delete).Methods(http.MethodDelete)

	api := request.PathPrefix("/newsapi").Subrouter()
	api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)
	api.HandleFunc("/keyword={keyword}", params).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8765", request))
}