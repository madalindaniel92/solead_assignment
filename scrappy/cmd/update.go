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
	"encoding/json"
	"examples/scrappy/internal/es"
	"fmt"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:          "update <domain url / id> <key> <JSON encoded value>",
	Short:        "Updates a field of company information in ElasticSearch",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateCompanyAction(args[0], args[1], args[2])
	},
}

func init() {
	esCmd.AddCommand(updateCmd)
}

func updateCompanyAction(id string, key string, encodedValue string) error {
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

	// JSON decode value
	var value interface{}
	err = json.Unmarshal([]byte(encodedValue), &value)
	if err != nil {
		return fmt.Errorf("failed to JSON decode value: %q", err)
	}

	// Updateable document fragment
	doc := map[string]interface{}{
		key: value,
	}

	// Update company info
	ctx := context.Background()
	err = client.UpdateCompanyInfo(ctx, id, doc)
	if err != nil {
		return fmt.Errorf("failed to update company info: %s", err)
	}

	return nil
}
