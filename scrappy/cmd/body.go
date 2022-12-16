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

	"github.com/spf13/cobra"
)

// bodyCmd represents the nav command
var bodyCmd = &cobra.Command{
	Use:          "body",
	Short:        "Get all website links from HTML a elements",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return bodyLinksAction(args[0])
	},
}

func init() {
	linksCmd.AddCommand(bodyCmd)
}

func bodyLinksAction(url string) error {
	links, err := web.GetLinks(url, "body")
	if err != nil {
		return err
	}

	printLinks(links)
	return nil
}
