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

const phoneFlagKey = "phone"

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:          "search <search query>",
	Short:        "Query ElasticSearch for company by name",
	SilenceUsage: true,
	Args:         cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		phone, err := cmd.Flags().GetString(phoneFlagKey)
		if err != nil {
			return err
		}

		return searchCompany(query, phone)
	},
}

func init() {
	esCmd.AddCommand(searchCmd)

	searchCmd.Flags().String(phoneFlagKey, "", "phone number to search by")
}

func searchCompany(query string, phone string) error {
	// Get ElasticSearch config
	config, err := esConfig()
	if err != nil {
		return err
	}

	// Initialize a new ES client
	client, err := es.NewClient(config)
	if err != nil {
		return err
	}

	// Search for company
	ctx := context.Background()
	result, err := client.SearchCompany(ctx, query, phone)
	if err != nil {
		return err
	}

	printSearchCompaniesResult(query, result)
	return nil
}

func printSearchCompaniesResult(query string, result *es.SearchCompaniesResult) {
	fmt.Printf("%d match the query: %q\n\n", result.Total, query)
	for _, company := range result.Companies {
		printCompanyResult(&company)
		fmt.Println()
	}
}

func printCompanyResult(company *es.Company) {
	printCompanyInfo(&company.Company)
	printCompanyPhoneNumbers(company.PhoneNumbers)
}

func printCompanyPhoneNumbers(phoneNumbers []string) {
	if len(phoneNumbers) > 0 {
		fmt.Println("Phone numbers:")
	}

	for _, phoneNumber := range phoneNumbers {
		fmt.Println("    -", phoneNumber)
	}
}
