package main

import (
	"context"
	"fmt"
	"github.com/GoServer/newsapi"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func get(httpResponseWriter http.ResponseWriter, request *http.Request) {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusOK)
	httpResponseWriter.Write([]byte(`{"message": "get called"}`))
}

func post(httpResponseWriter http.ResponseWriter, request *http.Request) {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusCreated)
	httpResponseWriter.Write([]byte(`{"message": "post called"}`))
}

func put(httpResponseWriter http.ResponseWriter, request *http.Request) {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusAccepted)
	httpResponseWriter.Write([]byte(`{"message": "put called"}`))
}

func delete(httpResponseWriter http.ResponseWriter, request *http.Request) {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusOK)
	httpResponseWriter.Write([]byte(`{"message": "delete called"}`))
}

func notFound(httpResponseWriter http.ResponseWriter, request *http.Request) {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusNotFound)
	httpResponseWriter.Write([]byte(`{"message": "not found"}`))
}

func params(httpResponseWriter http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	httpResponseWriter.Header().Set("Content-Type", "application/json")

	if _, ok := pathParams["keyword"]; ok {
		c := newsapi.NewClient("6c5c888290f647818122022f271a88f0", newsapi.WithHTTPClient(http.DefaultClient))

		sources, err := c.GetSources(context.Background(), nil)

		if err != nil {
			panic(err)
		}

		for _, s := range sources.Sources {
			httpResponseWriter.Write([]byte(fmt.Sprintf(s.Description)))
		}
	}

	userID := -1
	var err error
	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			httpResponseWriter.WriteHeader(http.StatusInternalServerError)
			httpResponseWriter.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	commentID := -1
	if val, ok := pathParams["commontID"]; ok {
		commentID, err = strconv.Atoi(val)
		if err != nil {
			httpResponseWriter.WriteHeader(http.StatusInternalServerError)
			httpResponseWriter.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}
	httpResponseWriter.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d}`, userID, commentID)))
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