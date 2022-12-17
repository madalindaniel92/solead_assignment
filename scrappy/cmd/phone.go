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

	"github.com/spf13/cobra"
)

// phoneCmd represents the phone command
var phoneCmd = &cobra.Command{
	Use:          "phone <website url>",
	Short:        "Scrape website for phone numbers",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		raw, err := cmd.Flags().GetBool("raw")
		if err != nil {
			return err
		}

		return phoneAction(args[0], raw)
	},
}

func init() {
	scrapeCmd.AddCommand(phoneCmd)
	phoneCmd.Flags().Bool("raw", false, "show raw scraped phone numbers, without validation or deduplication")
}

func phoneAction(url string, raw bool) error {
	phoneNumbers, err := web.GetPhoneNums(url)
	if err != nil {
		return err
	}

	// Validate and deduplicate phone numbers when raw flag is not set.
	if !raw {
		var invalid []web.FailedValidation
		phoneNumbers, invalid = web.ValidatePhoneNumbers(phoneNumbers)
		printInvalidPhoneNumbers(invalid)

		phoneNumbers = web.DedupPhoneNumbers(phoneNumbers)
	}

	fmt.Printf("Domain: %q\n", url)
	printPhoneNumbers(phoneNumbers)
	return nil
}

func printPhoneNumbers(phoneNumbers []web.Phone) {
	for index, phone := range phoneNumbers {
		fmt.Printf("%2d. %q (%s)\n", index, phone.Number, phone.Confidence)
	}
}

func printInvalidPhoneNumbers(invalid []web.FailedValidation) {
	if len(invalid) == 0 {
		return
	}

	fmt.Printf("%d invalid phone number(s)\n", len(invalid))
	for _, entry := range invalid {
		fmt.Printf("Invalid phone number: %q (%q)\n", entry.Number, entry.Err)
	}
	fmt.Println()
}
