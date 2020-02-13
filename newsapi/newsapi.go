package newsapi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-querystring/query"
)

type Client struct {
	client *http.Client
	baseURL *url.URL
	apiKey string
	userAgent string
}

type OptionFunc func(*Client)

func WithHTTPClient(httpClient *http.Client) OptionFunc {
	return func(client *Client) {
		client.client = httpClient
	}
}

func WithUserAgent(userAgent string) OptionFunc {
	return func(client *Client) {
		client.userAgent = userAgent
	}
}

func NewClient(apiKey string, options ...OptionFunc) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	client := &Client{
		client:    &http.Client{ Timeout: time.Second * 10 },
		baseURL:   baseURL,
		apiKey:    apiKey,
		userAgent: "github.com/GoServer/newsapi",
	}

	for _, opt := range options { opt(client) }

	return client
}

// newGetRequest returns a new Get request for the given url URLStr
// It returns a pointer to a http request which can be executed by a http.client
func (client *Client) newGetRequest(urlString string) (*http.Request, error) {
	relevantUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	finalUrl := client.baseURL.ResolveReference(relevantUrl)

	req, err := http.NewRequest("GET", finalUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(apiKeyHeader, client.apiKey)

	if client.userAgent != "" {
		req.Header.Set(userAgentHeader, client.userAgent)
	}

	return req, nil
}

// do executes the http.Request and marshal's the response into v
// v must be a pointer to a value instead of a regular value
// It returns the actual response from the request and also an error
func (client *Client) do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := client.client.Do(req)

	if err != nil {
		return resp, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err == io.EOF {
		err = nil
	}

	return resp, err
}

// setOptions set the options for the url
// It will set the query parameters and encodes them with the url
func setOptions(endpoint string, options interface{}) (string, error) {
	urlWithOptions, err := url.Parse(endpoint)
	if err != nil {
		return endpoint, err
	}

	qs, err := query.Values(options)
	if err != nil {
		return endpoint, err
	}

	urlWithOptions.RawQuery = qs.Encode()
	return urlWithOptions.String(), nil
}