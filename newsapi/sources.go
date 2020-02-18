package newsapi

import "context"

// GetSources returns the sources from newsapi see http://newsapi.org/docs for more information on the parameters
func (client *Client) GetSources(ctx context.Context, params *SourceParameters) (*SourceResponse, error) {
	u := sourcesEndpoint

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
		*SourceResponse
	}{}

	_, err = client.do(ctx, req, &response)
	if err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, response.Error
	}

	return response.SourceResponse, nil
}