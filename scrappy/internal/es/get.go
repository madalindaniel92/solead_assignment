package es

import (
	"context"
	"encoding/json"
	"fmt"
)

type getCompanyResult struct {
	Found   bool    `json:"found"`
	Company Company `json:"_source"`
}

// GetCompany gets a company from ElasticSearch by url.
//
// If the company can't be found, ErrNotFound is returned.
func (c *Client) GetCompany(ctx context.Context, url string) (*Company, error) {
	id, err := urlToId(url)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Get(c.companiesIndex, id,
		c.client.Get.WithContext(ctx),
	)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check errors returned by ES
	if res.IsError() {
		return nil, errorFromResponse(res)
	}

	// Decode ES search response
	var envelope getCompanyResult
	err = json.NewDecoder(res.Body).Decode(&envelope)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedResponse, err)
	}

	if !envelope.Found {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, err)
	}

	return &envelope.Company, nil
}
