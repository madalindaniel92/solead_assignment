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
	"context"
	"fmt"

	"examples/scrappy/internal/es"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:          "search <search query>",
	Short:        "Query ElasticSearch for company by name",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return searchCompany(args[0])
	},
}

func init() {
	esCmd.AddCommand(searchCmd)
}

func searchCompany(query string) error {
	if query == "" {
		return fmt.Errorf("missing query argument")
	}

	config, err := esConfig()
	if err != nil {
		return err
	}

	client, err := es.NewClient(config)
	if err != nil {
		return err
	}

	ctx := context.Background()
	result, err := client.SearchCompany(&ctx, query)
	if err != nil {
		return err
	}

	printSearchCompaniesResult(query, result)
	return nil
}

func printSearchCompaniesResult(query string, result *es.SearchCompaniesResult) {
	fmt.Printf("%d match the query: %q\n\n", result.Total, query)
	for _, company := range result.Companies {
		printCompanyInfo(&company.Company)
	}
}
