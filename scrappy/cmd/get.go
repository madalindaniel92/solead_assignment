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
	"examples/scrappy/internal/es"
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:          "get <domain url>",
	Short:        "Get company information from ElasticSearch",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return getCompanyAction(args[0])
	},
}

func init() {
	esCmd.AddCommand(getCmd)
}

func getCompanyAction(url string) error {
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

	// Get company by domain url
	ctx := context.Background()
	company, err := client.GetCompany(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to get company %q: %s", url, err)
	}

	printCompanyResult(company)
	return nil
}
