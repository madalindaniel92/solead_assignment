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
	"examples/scrappy/internal/csv"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check input files",
}

func init() {
	rootCmd.AddCommand(checkCmd)
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
