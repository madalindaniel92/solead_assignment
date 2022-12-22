package es

import (
	"bytes"
	"context"
	"encoding/json"
	"examples/scrappy/internal/csv"
	"fmt"
	"io"
)

var companyNameFields = []string{
	"commercial_name",
	"legal_name",
	"all_available_names",
}

type Company struct {
	csv.Company
	ID           string   `json:"id"`
	PhoneNumbers []string `json:"phone_numbers,omitempty"`
}

type SearchCompaniesResult struct {
	Total     int       `json:"total"`
	Companies []Company `json:"companies"`
}

// Documentation example code for querying ES
// https://github.com/elastic/go-elasticsearch/blob/main/_examples/xkcdsearch/store.go

type searchEnvelope struct {
	Took int
	Hits struct {
		Total struct {
			Value int
		}
		Hits []struct {
			ID         string          `json:"_id"`
			Source     json.RawMessage `json:"_source"`
			Highlights json.RawMessage `json:"_highlight"`
			Sort       []interface{}   `json:"sort"`
		}
	}
}

// esErrorMessage decodes an error response from the ES cluster
type esErrorMessage struct {
	Error esErrorInfo `json:"error"`
}

type esErrorInfo struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

// SearchCompany searches ElasticSearch for a company by name.
func (c *Client) SearchCompany(ctx context.Context, query string) (*SearchCompaniesResult, error) {
	var result SearchCompaniesResult
	searchAPI := c.client.Search

	// Send search request to ES
	res, err := c.client.Search(
		searchAPI.WithIndex(c.companiesIndex),
		searchAPI.WithBody(searchCompanyByNameQuery(query)),
		searchAPI.WithContext(ctx),
	)

	// Check network errors
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrSearchResult, err)
	}
	defer res.Body.Close()

	// Check errors returned by ES
	if res.IsError() {
		return nil, errorFromResponse(res)
	}

	// Decode ES search response
	var envelope searchEnvelope
	err = json.NewDecoder(res.Body).Decode(&envelope)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedResponse, err)
	}

	result.Total = envelope.Hits.Total.Value
	if len(envelope.Hits.Hits) < 1 {
		return &result, nil
	}

	for index, hit := range envelope.Hits.Hits {
		var company Company
		company.ID = hit.ID

		// Decode company
		err := json.Unmarshal(hit.Source, &company)
		if err != nil {
			return nil, fmt.Errorf("%w: %s (index %d)", ErrUnexpectedResponse, err, index)
		}

		result.Companies = append(result.Companies, company)
	}

	return &result, nil
}

func searchCompanyByNameQuery(query string) io.Reader {
	esQuery := h{
		"query": h{
			"multi_match": h{
				"query":  query,
				"fields": companyNameFields,
			},
		},
		"sort": a{
			h{"_score": "desc"},
			h{"_doc": "asc"},
		},
	}

	encoded, _ := json.Marshal(esQuery)
	return bytes.NewReader(encoded)
}
