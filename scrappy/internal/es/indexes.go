package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

// Index management functions examples from:
// https://github.com/elastic/go-elasticsearch/blob/main/_examples/xkcdsearch/store.go

func (c *Client) CreateCompanyIndex(ctx context.Context, indexName string) error {
	if indexName != c.companiesIndex {
		return fmt.Errorf("%w: only %s index is currently supported", ErrInvalidIndex, c.companiesIndex)
	}

	indexAPI := c.client.Indices
	createOp := indexAPI.Create

	res, err := indexAPI.Create(c.companiesIndex,
		createOp.WithBody(companyIndexMapping()),
		createOp.WithContext(ctx),
	)

	// Check network errors
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedRequest, err)
	}
	res.Body.Close()

	// Check errors returned by ES
	if res.IsError() {
		return errorFromResponse(res)
	}

	return nil
}

func companyIndexMapping() io.Reader {
	mapping := h{
		"mappings": h{
			"properties": h{
				"domain":              h{"type": "keyword"},
				"phone_numbers":       h{"type": "keyword"},
				"commercial_name":     h{"type": "text"},
				"legal_name":          h{"type": "text"},
				"all_available_names": h{"type": "text"},
			},
		},
	}

	encoded, _ := json.Marshal(mapping)
	return bytes.NewReader(encoded)
}
