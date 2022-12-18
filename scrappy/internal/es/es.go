// Package es implements ElasticSearch support functions.
package es

import (
	"net/http"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// Config represents the credentials needed to connect to the ElasticSearch cluster.
type Config struct {
	Username string
	Password string
	// CA certificate for Elastic Search cluster
	CACert []byte
	// URL addresses of cluster replicas
	Addresses []string
}

// Client is a wrapper around the go-elasticsearch client.
type Client struct {
	client *elastic.Client
	config *Config
}

func NewClient(config *Config) (*Client, error) {
	if len(config.CACert) == 0 {
		return nil, ErrMissingCACertificate
	}

	// Set extended client options
	esConfig := elastic.Config{
		Username:  config.Username,
		Password:  config.Password,
		CACert:    config.CACert,
		Addresses: config.Addresses,
		// Options from:
		// 		https://github.com/elastic/go-elasticsearch/blob/main/esutil/bulk_indexer_example_test.go
		//
		// Retry on 429 TooManyRequests status
		RetryOnStatus: []int{502, 503, 504, http.StatusTooManyRequests},
		// Simple incremental backoff using 100 milliseconds * index
		RetryBackoff: func(i int) time.Duration {
			return time.Duration(i) * 100 * time.Millisecond
		},
		// Retry up to 5 times
		MaxRetries: 5,
	}

	client, err := elastic.NewClient(esConfig)
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) Info() (*esapi.Response, error) {
	return c.client.Info()
}

// type BulkIndexResult struct{}

// // BulkIndexCompanies will index the companies into the ElasticSearch "companies" index.
// //
// // https://github.com/elastic/go-elasticsearch/blob/main/esutil/bulk_indexer_example_test.go
// func BulkIndexCompanies(companies []csv.Company) (*BulkIndexResult, error) {
// 	config, err := defaultConfig()
// 	if err != nil {
// 		return nil, err
// 	}

// 	es, err := elastic.NewClient(*config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(es.Info())
// 	return nil, nil
// }
