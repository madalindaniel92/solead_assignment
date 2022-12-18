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
	"log"
	"net/http"
	"runtime"

	"github.com/spf13/cobra"

	"examples/scrappy/internal/csv"
	"examples/scrappy/internal/web"
)

// domainsCmd represents the domains command
var domainsCmd = &cobra.Command{
	Use:          "domains <csv file to load domain names from>",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	Short:        "Check CSV file and send head http request to each company domain.",
	Long: `This command helps validate that the domains in the passed in CSV file
	are valid URLS and reachable.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		csvPath := args[0]
		numWorkers, err := cmd.Flags().GetInt("workers")
		if err != nil {
			return err
		}

		return domainAction(csvPath, numWorkers)
	},
}

func init() {
	checkCmd.AddCommand(domainsCmd)

	domainsCmd.Flags().Int("workers", runtime.NumCPU()*20,
		"number of concurrent workers (defaults to 20 * NumCPUs)")
}

func domainAction(csvPath string, numWorkers int) error {
	if csvPath == "" {
		return fmt.Errorf("missing csv file argument")
	}

	// Load website domains from CSV file
	websites, err := csv.LoadDomainsFromFile(csvPath)
	if err != nil {
		printExtraErrInfo(err)
		return err
	}

	// Collect urls
	urls := make([]string, 0, len(websites))
	for _, website := range websites {
		urls = append(urls, website.URL())
	}

	// Run URL checks asynchronously
	results := web.CheckURLs(urls, numWorkers, printDomainResult)

	// Aggregate results
	printDomainAggregateResults(results)
	return nil
}

func printDomainResult(result *web.CheckUrlResult) {
	url := result.URL()

	if result.Err != nil {
		log.Printf("Failed request to domain %q: %q\n", url, result.Err)
		return
	}

	log.Printf("HEAD %q - %d\n", url, result.Status)
}

func printDomainAggregateResults(results []web.CheckUrlResult) {
	successful, badRequests, failed := 0, 0, 0
	statusCount := map[int]int{}
	for _, result := range results {
		switch {
		case result.Status == http.StatusOK:
			successful++
		case result.Err != nil:
			failed++
		default:
			badRequests++
			statusCount[result.Status]++
		}
	}

	fmt.Printf("\nSuccessful requests: %d\n", successful)
	fmt.Printf("Failed to connect: %d\n", failed)
	fmt.Printf("Bad requests: %d\n", badRequests)
	for status, count := range statusCount {
		fmt.Printf("status %d - %d request(s)\n", status, count)
	}
}
