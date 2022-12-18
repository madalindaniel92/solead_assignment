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
	"examples/scrappy/internal/csv"
	"fmt"

	"github.com/spf13/cobra"
)

// companiesCmd represents the companies command
var companiesCmd = &cobra.Command{
	Use:          "companies <csv file to load company info from>",
	Short:        "Check companies CSV data",
	Long:         `This command helps validate that we can parse the CSV company info data`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return companiesAction(args[0])
	},
}

func init() {
	checkCmd.AddCommand(companiesCmd)
}

func companiesAction(csvPath string) error {
	if csvPath == "" {
		return fmt.Errorf("missing csv file argument")
	}

	companies, err := csv.LoadCompaniesFromFile(csvPath)
	if err != nil {
		printExtraErrInfo(err)
		return err
	}

	printCompaniesInfo(companies)
	return nil
}

func printCompaniesInfo(companies []csv.Company) {
	for _, company := range companies {
		printCompanyInfo(&company)
		fmt.Println()
	}
}

func printCompanyInfo(company *csv.Company) {
	fmt.Println("Domain:", company.Domain.String())
	printField("Commercial name:", company.CommercialName)
	printField("Legal name:", company.LegalName)
	printCompanyAvailableNames(company.AllAvailableNames)
}

func printCompanyAvailableNames(names []string) {
	if len(names) > 0 {
		fmt.Println("Other names:")
	}

	for _, name := range names {
		fmt.Println("    -", name)
	}
}

func printField(label string, value string) {
	if value != "" {
		fmt.Println(label, value)
	}
}
