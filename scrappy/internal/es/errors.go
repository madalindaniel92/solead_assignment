package es

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var (
	ErrMissingCACertificate = errors.New("missing CA certificate")
	ErrSearchResult         = errors.New("error on search")
	ErrUnexpectedResponse   = errors.New("unexpected response from ES")
	ErrFailedRequest        = errors.New("failed request")
	ErrNotFound             = errors.New("not found")
	ErrInvalidIndex         = errors.New("invalid index")
)

// errorFromResponse extracts error information from the ElasticSearch response.
func errorFromResponse(response *esapi.Response) error {
	var e esErrorMessage
	err := json.NewDecoder(response.Body).Decode(&e)
	if err != nil {
		return fmt.Errorf("%w: [%d] %s",
			ErrUnexpectedResponse, response.StatusCode, err)
	}

	// Build detailed error message
	return fmt.Errorf("%w: [%s] %s: %s", ErrFailedRequest,
		response.Status(), e.Error.Type, e.Error.Reason)
}
