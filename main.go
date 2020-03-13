package main

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"io/ioutil"

	//"fmt"
	"github.com/GoServer/newsapi"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"

	//firebase "firebase.google.com/go"
	//"google.golang.org/api/option"
)



const (
	key_q              = "q"
	key_qInTitle       = "qInTitle"
	key_category       = "category"
	key_country        = "country"
	key_sources        = "sources"
	key_domains        = "domains"
	key_excludeDomains = "excludeDomains"
	key_from           = "from"
	key_to             = "to"
	key_language       = "language"
	key_sortBy         = "sortBy"
	key_pageSize       = "pageSize"
	key_page           = "page"

	INT_MAX = int(^uint(0) >> 1)
)

var headerMap = map[string]string{"Content-Type": "application/json"}

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

func trendingParams(httpResponseWriter http.ResponseWriter, request *http.Request) {
	setUpHttpResponseHeader(httpResponseWriter)
	c := newsapi.NewClient("6c5c888290f647818122022f271a88f0", newsapi.WithHTTPClient(http.DefaultClient))

	queryParams := request.URL.Query()
	topHeadlineParams := newsapi.TopHeadlineParameters{}
	countries, ok := queryParams[key_country]
	if ok && len(countries) > 0 {
		topHeadlineParams.Country = countries[0]
	}

	categories, ok := queryParams[key_category]
	if ok && len(categories) > 0 {
		topHeadlineParams.Category = categories[0]
	}

	sourcesStr, ok := queryParams[key_sources]
	if ok && len(sourcesStr) > 0 {
		topHeadlineParams.Sources = strings.Split(sourcesStr[0], ",")
	}

	queries, ok := queryParams[key_q]
	if ok && len(queries) > 0 {
		topHeadlineParams.Keywords = queries[0]
	}

	pageSizes, ok := queryParams[key_pageSize]
	if ok && len(pageSizes) > 0 {
		pageSize, err := strconv.Atoi(pageSizes[0])
		if err != nil {
			topHeadlineParams.PageSize = pageSize
		}
	}

	pages, ok := queryParams[key_page]
	if ok && len(pages) > 0 {
		page, err := strconv.Atoi(pages[0])
		if err != nil {
			topHeadlineParams.Page = page
		}
	}

	articles, err := c.GetTopHeadlines(context.Background(), &topHeadlineParams)
	handleApiError(err, httpResponseWriter)
	response, err := json2.Marshal(articles)
	handleJsonResponse(response, err, httpResponseWriter)
}

func searchParams(httpResponseWriter http.ResponseWriter, request *http.Request) {
	setUpHttpResponseHeader(httpResponseWriter)
	c := newsapi.NewClient("6c5c888290f647818122022f271a88f0", newsapi.WithHTTPClient(http.DefaultClient))

	queryParams := request.URL.Query()
	searchParams := newsapi.EverythingParameters{}
	qs, ok := queryParams[key_q]
	if ok && len(qs) > 0 {
		searchParams.Keywords = qs[0]
	}

	qsInTitle, ok := queryParams[key_qInTitle]
	if ok && len(qsInTitle) > 0 {
		searchParams.KeywordsInTitle = qsInTitle[0]
	}

	sourcesStr, ok := queryParams[key_sources]
	if ok && len(sourcesStr) > 0 {
		searchParams.Sources = strings.Split(sourcesStr[0], ",")
	}

	languages, ok := queryParams[key_language]
	if ok && len(languages) > 0 {
		searchParams.Language = languages[0]
	}

	sortBys, ok := queryParams[key_sortBy]
	if ok && len(sortBys) > 0 {
		searchParams.SortBy = sortBys[0]
	} else {
		searchParams.SortBy = "publishedAt"
	}

	pageSizes, ok := queryParams[key_pageSize]
	if ok && len(pageSizes) > 0 {
		pageSize, err := strconv.Atoi(pageSizes[0])
		if err != nil {
			searchParams.PageSize = pageSize
		}
	}

	pages, ok := queryParams[key_page]
	if ok && len(pages) > 0 {
		page, err := strconv.Atoi(pages[0])
		if err != nil {
			searchParams.Page = page
		}
	}

	articles, err := c.GetEverything(context.Background(), &searchParams)
	handleApiError(err, httpResponseWriter)
	response, err := json2.Marshal(articles)
	handleJsonResponse(response, err, httpResponseWriter)
}

func sourcesParams(httpResponseWriter http.ResponseWriter, request *http.Request) {
	setUpHttpResponseHeader(httpResponseWriter)
	c := newsapi.NewClient("6c5c888290f647818122022f271a88f0", newsapi.WithHTTPClient(http.DefaultClient))

	queryParams := request.URL.Query()
	sourcesParams := newsapi.SourceParameters{}
	countries, ok := queryParams[key_country]
	if ok && len(countries) > 0 {
		sourcesParams.Country = countries[0]
	}

	categories, ok := queryParams[key_category]
	if ok && len(categories) > 0 {
		sourcesParams.Category = categories[0]
	}

	languages, ok := queryParams[key_language]
	if ok && len(languages) > 0 {
		sourcesParams.Language = languages[0]
	}

	sources, err := c.GetSources(context.Background(), &sourcesParams)
	handleApiError(err, httpResponseWriter)
	response, err := json2.Marshal(sources)
	handleJsonResponse(response, err, httpResponseWriter)
}

func setUpHttpResponseHeader(httpResponseWriter http.ResponseWriter) {
	httpResponseWriter.Header().Set("Content-Type", "application/json")
	httpResponseWriter.WriteHeader(http.StatusOK)
}

func handleApiError(err error, httpResponseWriter http.ResponseWriter) {
	if err != nil {
		if newsapi.ApiError(err) {
			errContent, err := json2.Marshal(err.(*newsapi.Error))
			handleJsonResponse(errContent, err, httpResponseWriter)
		} else {
			constructResponse(httpResponseWriter, headerMap, http.StatusInternalServerError, []byte(`{"message": "something is wrong"}`))
		}
		return
	}
}

func handleJsonResponse(response []byte, err error, httpResponseWriter http.ResponseWriter) {
	if err != nil {
		http.Error(httpResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	httpResponseWriter.Write(response)
}

func postHealthDataParams(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", reqBody)
}

func main() {
	// FCM push notification
	//opt := option.WithCredentialsFile("github.com/GoServer/instant-news-7840b-firebase-adminsdk-d7ns0-aa18d49621.json")
	//app, err := firebase.NewApp(context.Background(), nil, opt)
	//if err != nil {
	//	log.Fatalf("error initializing app: %v\n", err)
	//}
	

	request := mux.NewRouter()

	apiTest := request.PathPrefix("/api/test").Subrouter()
	apiTest.HandleFunc("", get).Methods(http.MethodGet)
	apiTest.HandleFunc("", post).Methods(http.MethodPost)
	apiTest.HandleFunc("", put).Methods(http.MethodPut)
	apiTest.HandleFunc("", delete).Methods(http.MethodDelete)

	newsApi := request.PathPrefix("/api/news").Subrouter()
	newsApi.HandleFunc("/top-headlines", trendingParams).Methods(http.MethodGet)
	newsApi.HandleFunc("/everything", searchParams).Methods(http.MethodGet)
	newsApi.HandleFunc("/sources", sourcesParams).Methods(http.MethodGet)

	psApi := request.PathPrefix("/api/ps").Subrouter()
	psApi.HandleFunc("/healthData", postHealthDataParams).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", request))
	//log.Fatal(http.ListenAndServe(":8080", request))
}
