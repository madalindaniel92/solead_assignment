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

// navCmd represents the nav command
var navCmd = &cobra.Command{
	Use:          "nav",
	Short:        "Get website links from HTML nav element",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return navLinksAction(args[0])
	},
}

func init() {
	linksCmd.AddCommand(navCmd)
}

func navLinksAction(url string) error {
	links, err := web.GetLinks(url, "nav")
	if err != nil {
		return err
	}

	printLinks(links)
	return nil
}
