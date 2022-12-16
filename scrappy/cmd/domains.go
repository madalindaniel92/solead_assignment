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
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"examples/scrappy/internal/csv"
	"examples/scrappy/internal/web"
)

// domainsCmd represents the domains command
var domainsCmd = &cobra.Command{
	Use:          "domains <csv file to load domain names from>",
	Aliases:      []string{"c"},
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	Short:        "Check CSV file and send head http request to each company domain.",
	Long: `This command helps validate that the domains in the passed in CSV file
	are valid URLS and reachable.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return domainAction(args[0])
	},
}

func init() {
	checkCmd.AddCommand(domainsCmd)
}

func domainAction(csvPath string) error {
	if csvPath == "" {
		return fmt.Errorf("missing csv file argument")
	}

	websites, err := csv.LoadFromFile(csvPath)
	if err != nil {
		printExtraErrInfo(err)
		return err
	}

	// for _, result := range results {
	// 	domain := result.Domain
	// 	fmt.Printf("Domain: %s://%s\n", domain.Scheme, domain.Hostname())
	// }

	websites = websites[:5]
	for _, website := range websites {
		url := website.Domain.String()
		statusCode, err := web.CheckURL(url)
		if err != nil {
			log.Printf("Failed request to domain %q: %q\n", url, err)
			continue
		}

		log.Printf("HEAD %q - %d\n", url, statusCode)
	}

	return nil
}

func printExtraErrInfo(err error) {
	switch value := errors.Unwrap(err).(type) {
	case csv.ErrInvalidCSVLines:
		for _, invalidLine := range value {
			fmt.Fprintf(os.Stderr, "Invalid line %d %q: %s\n",
				invalidLine.Index, invalidLine.Line, invalidLine.Err)
		}
	}
}
