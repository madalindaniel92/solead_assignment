package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// UpdateCompanyInfo updates the associated field for the company
func (c *Client) UpdateCompanyInfo(ctx context.Context, url string, info map[string]any) error {
	doc, err := encodeUpdateDoc(info)
	if err != nil {
		return err
	}

	url = strings.TrimPrefix(url, "https://")

	response, err := c.client.Update(c.companiesIndex, url, doc,
		c.client.Update.WithContext(ctx),
	)

	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedRequest, err)
	}

	if response.IsError() {
		return errorFromResponse(response)
	}

	return nil
}

func encodeUpdateDoc(info map[string]any) (io.Reader, error) {
	payload := h{"doc": info}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(encoded), err
}
