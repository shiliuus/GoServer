package newsapi

import "time"

const (
	apiKeyHeader    = "X-Api-Key"
	userAgentHeader = "User-Agent"
	defaultBaseURL  = "https://newsapi.org/v2/"

	sourcesEndpoint      = "sources"
	topHeadlinesEndpoint = "top-headlines"
	everythingEndpoint   = "everything"
)

// url request params
type UrlParams struct {
	Category string `url:"category, omitempty"`
	Language string `url:"language,omitempty"`
	Country  string `url:"country,omitempty"`
}

// SourceResponse is the response from the source request
type SourceResponse struct {
	Status  string   `json:"status"`
	Sources []Source `json:"sources"`
}

type Source struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Category    string `json:"category"`
	Language    string `json:"language"`
	Country     string `json:"country"`
	UrlsToLogos struct {
		Small  string `json:"small"`
		Medium string `json:"medium"`
		Large  string `json:"large"`
	} `json:"urlsToLogos"`
}

// TopHeadlineParameters are the parameters which can be used to tweak to request for the top headlines.
type TopHeadlineParameters struct {
	Country  string   `url:"country,omitempty"`
	Category string   `url:"category,omitempty"`
	Sources  []string `url:"sources,omitempty,comma"`
	Keywords string   `url:"q,omitempty"`
	Page     int      `url:"page,omitempty"`
	PageSize int      `url:"pageSize,omitempty"`
}

// EverythingParameters are the parameters used for the newsapi everything endpoint.
type EverythingParameters struct {
	Keywords       string   `url:"q,omitempty"`
	Sources        []string `url:"sources,omitempty,comma"`
	Domains        []string `url:"domains,omitempty,comma"`
	ExcludeDomains []string `url:"excludeDomains,omitempty"`

	From time.Time `url:"from,omitempty"`
	To   time.Time `url:"to,omitempty"`

	Language string `url:"language,omitempty"`
	SortBy   string `url:"sortBy,omitempty"`

	Page     int `url:"page,omitempty"`
	PageSize int `url:"pageSize,omitempty"`
}

// Article is a single article from the newsapi article response
// See http://newsapi.org/docs for more details on the property's
type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

// ArticleResponse is the response from the newsapi article endpoint.
// Code and Message property will be filled when an error happened.
// See http://newsapi.org/docs for more details on the property's.
type ArticleResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Error struct {
	Status  string `json:"status, omitempty"`
	Code    string `json:"code, omitempty"`
	Message string `json:"message, omitempty"`
}