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
	"examples/scrappy/internal/web"
	"fmt"
	"log"
	"runtime"

	"github.com/spf13/cobra"
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:          "scrape <csv file to load domain names from>",
	Short:        "Scrape domains for phone numbers",
	Aliases:      []string{"s"},
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		csvPath := args[0]
		numWorkers, err := cmd.Flags().GetInt("workers")
		if err != nil {
			return err
		}

		return scrapeDomainsAction(csvPath, numWorkers)
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	scrapeCmd.Flags().Int("workers", runtime.NumCPU()*20,
		"number of concurrent workers (defaults to 20 * NumCPUs)")

}

type scrapeResult struct {
	// Number of domains for which we have collected phone numbers
	phoneNumbersCollected int
}

func scrapeDomainsAction(csvPath string, numWorkers int) error {
	// Load website URLs from CSV file
	urls, err := loadDomainUrls(csvPath)
	if err != nil {
		return err
	}

	var stats scrapeResult

	// Scrape domains and handle each job result.
	web.ScrapeDomains(urls, numWorkers, func(result *web.ScrapeJobResult) {
		if result.Err != nil {
			log.Printf("Failed request to domain %q: %q\n", result.Url, result.Err)
			return
		}

		url, info := result.Url, result.Info
		if len(info.PhoneNumbers) > 0 {
			stats.phoneNumbersCollected++
		}

		fmt.Printf("Info for %q: %#v\n", url, info)
	})

	printScrapeResultStats(&stats)
	return nil
}

func printScrapeResultStats(stats *scrapeResult) {
	if stats.phoneNumbersCollected > 0 {
		fmt.Printf("Collected phone numbers for %d domain(s)\n",
			stats.phoneNumbersCollected)
	}
}
