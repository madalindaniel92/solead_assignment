// Package es implements ElasticSearch support functions.
package es

import (
	"bytes"
	"context"
	"encoding/json"
	"examples/scrappy/internal/csv"
	"log"
	"net/http"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

const companiesESIndex = "companies"

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

// BulkIndexCompanies will index the companies into the ElasticSearch "companies" index.
//
// https://github.com/elastic/go-elasticsearch/blob/main/esutil/bulk_indexer_example_test.go
func (c *Client) BulkIndexCompanies(companies []csv.Company) (*esutil.BulkIndexerStats, error) {
	ctx := context.Background()

	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: c.client,
		Index:  companiesESIndex,
	})

	if err != nil {
		return nil, err
	}

	for _, company := range companies {
		payload, err := json.Marshal(company)
		if err != nil {
			return nil, err
		}

		err = indexer.Add(ctx, esutil.BulkIndexerItem{
			Action: "index",
			// We use the domain host of the company as a natural key
			DocumentID: company.Domain.Hostname(),
			Body:       bytes.NewReader(payload),
			OnSuccess:  handleBulkIndexSuccess,
			OnFailure:  handleBulkIndexFailure,
		})

		if err != nil {
			log.Printf("es index error for company: %q\n", &company.Domain)
		}
	}

	// Close the indexer channel and flush remaining items
	err = indexer.Close(ctx)
	if err != nil {
		return nil, err
	}

	stats := indexer.Stats()
	return &stats, nil
}

func handleBulkIndexSuccess(
	ctx context.Context,
	item esutil.BulkIndexerItem,
	res esutil.BulkIndexerResponseItem) {

	log.Printf("es successfully indexed %q [%d] %s\n",
		item.DocumentID, res.Status, res.Result)
}

func handleBulkIndexFailure(
	ctx context.Context,
	item esutil.BulkIndexerItem,
	res esutil.BulkIndexerResponseItem,
	err error) {

	if err != nil {
		log.Printf("es index error: %q", err)
	} else {
		log.Printf("es index error: %q - %q", res.Error.Type, res.Error.Reason)
	}
}
