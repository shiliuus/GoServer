package main

import (
	"context"
	json2 "encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"io/ioutil"
	"time"

	//"fmt"
	"github.com/GoServer/newsapi"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"

	//"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jackc/pgx"
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

	dbhost = "psdev.ctmxeolrv0ba.us-east-1.rds.amazonaws.com"
	dbport = 5432
	dbname = "glucosetracker"
	dbuser = "glucosetracker"
	dbpassword = "password001"
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
	// read the dto,  and send the data to database
	var healthData newsapi.HealthData
	json2.Unmarshal([]byte(reqBody), &healthData)
	fmt.Printf("Result: ", healthData)

	// send a push notification to tell
	opt := option.WithCredentialsFile("personal-science-firebase-adminsdk-56hl6-3a6fff79c1.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatal("error getting Messaging client: %v\n", err)
	}

	registrationToken := healthData.Token

	message := &messaging.Message{
		Data: map[string]string {
			"score" : "850",
			"time" : "2:45",
		},
		Token: registrationToken,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Successfully sent message:", response)

	// Send the data to database
	connConfig, err := pgx.ParseConfig("user=glucosetracker password=password001 host=psdev.ctmxeolrv0ba.us-east-1.rds.amazonaws.com port=5432 dbname=glucosetracker")

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close(context.Background())

	sqlcommand := "INSERT INTO public.metrics_metric_test (created, modified, start_datetime, stop_datetime, metric_type, value, source, user_entered, user_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9);"
	t := time.Now()
	createAt := t.Format("2006-01-02 15:04:05")

	for _, data := range healthData.FitDataList {
		_, err2 := conn.Exec(context.Background(),
			sqlcommand,
			createAt, createAt, data.DateFrom, data.DateTo, data.Source, strconv.Itoa(data.Value), data.Source, data.UserEntered, 24)
		if err2 != nil {
			log.Fatalln(err2)
		}
		fmt.Printf("Successfully created user mwood\n")
	}
	
}

func main() {
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

	log.Fatal(http.ListenAndServe(":8765", request))
	//log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", request))
}
