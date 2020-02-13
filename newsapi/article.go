package newsapi

import (
	"context"
)

// GetTopHeadlines returns the articles from newsapi
// See http://newsapi.org/docs for more information
// It will return the error from newsapi if there is an error
func (client *Client) GetTopHeadlines(ctx context.Context, params *TopHeadlineParameters) (*ArticleResponse, error) {
	return client.getArticles(ctx, topHeadlinesEndpoint, params)
}

// GetEverything returns the articles from newsapi
// See http://newsapi.org/docs for more information
// It will return the error from newsapi if there is an error
func (client *Client) GetEverything(ctx context.Context, params *EverythingParameters) (*ArticleResponse, error) {
	return client.getArticles(ctx, everythingEndpoint, params)
}

// GetArticles returns the articles from newsapi
// See http://newsapi.org/docs for more information
// It will return the error from newsapi if there is an error
func (client *Client) getArticles(ctx context.Context, u string, params interface{}) (*ArticleResponse, error) {
	if params != nil {
		var err error
		u, err = setOptions(u, params)

		if err != nil {
			return nil, err
		}
	}

	req, err := client.newGetRequest(u)
	if err != nil {
		return nil, err
	}

	var response = struct {
		*Error
		*ArticleResponse
	}{}

	_, err = client.do(ctx, req, &response)

	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return response.ArticleResponse, nil
}
