/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"examples/scrappy/internal/csv"
	"examples/scrappy/internal/es"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var importCmd = &cobra.Command{
	Use:          "import <csv file to load company info from>",
	Short:        "Import company information into ElasticSearch index",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return importCompanies(args[0])
	},
}

func init() {
	esCmd.AddCommand(importCmd)
}

func importCompanies(csvPath string) error {
	if csvPath == "" {
		return fmt.Errorf("missing csv file argument")
	}

	config, err := esConfig()
	if err != nil {
		return err
	}

	client, err := es.NewClient(config)
	if err != nil {
		return err
	}

	// Load company info from CSV
	companies, err := csv.LoadCompaniesFromFile(csvPath)
	if err != nil {
		printExtraErrInfo(err)
		return err
	}

	// Bulk index companies into ElasticSearch
	stats, err := client.BulkIndexCompanies(companies)
	if err != nil {
		return err
	}

	if stats.NumFailed > 0 {
		fmt.Printf("Indexed [%d] documents with [%d] errors\n", stats.NumFlushed, stats.NumFailed)
	} else {
		fmt.Printf("Successfully indexed [%d] documents\n", stats.NumFlushed)
	}

	return nil
}
