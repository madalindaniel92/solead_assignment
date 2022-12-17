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
		return phoneAction(args[0])
	},
}

func init() {
	scrapeCmd.AddCommand(phoneCmd)
}

func phoneAction(url string) error {
	phoneNumbers, err := web.GetPhoneNums(url)
	if err != nil {
		return err
	}

	fmt.Printf("Domain: %q\n", url)
	for index, phone := range phoneNumbers {
		fmt.Printf("%2d. %q (%s)\n", index, phone.Number, phone.Confidence)
	}

	return nil
}
