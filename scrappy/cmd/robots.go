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
	"io"
	"os"

	"github.com/spf13/cobra"
)

// robotsCmd represents the robots command
var robotsCmd = &cobra.Command{
	Use:          "robots <domain url>",
	Short:        "Download and display robots.txt file for given domain",
	Aliases:      []string{"r"},
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return robotsAction(os.Stdout, args[0])
	},
}

func init() {
	getCmd.AddCommand(robotsCmd)
}

func robotsAction(out io.Writer, url string) error {
	result, err := web.GetRobots(url)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	// Write result body to stdout
	_, err = io.Copy(out, result.Body)
	return err
}
