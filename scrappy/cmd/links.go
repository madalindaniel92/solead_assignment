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

	"github.com/spf13/cobra"
)

// linksCmd represents the links command
var linksCmd = &cobra.Command{
	Use:     "links",
	Aliases: []string{"l"},
	Short:   "Strategies for getting the links from a website",
}

func init() {
	rootCmd.AddCommand(linksCmd)
}

func printLinks(links []string) {
	for index, link := range links {
		fmt.Printf("	%d %q\n", index, link)
	}
}
