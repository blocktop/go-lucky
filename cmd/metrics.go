// Copyright © 2018 J. Strobus White.
// This file is part of the blocktop blockchain development kit.
//
// Blocktop is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Blocktop is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with blocktop. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Retrieves metrics from a running lucky blockchain.",
	Long: `Usage: lucky metrics [OPTIONS]  (all metrics)
			 lucky metrics [SUBCOMMAND] [OPTIONS]  (specific metrics)
			 
Metrics can either be output in plain text or JSON format by using
the --format option. The default is plain text.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, c := range cmd.Commands() {
			_, err := c.ExecuteC()
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	},
}

var metricsInJson bool

func init() {
	rootCmd.AddCommand(metricsCmd)

	metricsCmd.PersistentFlags().BoolVarP(&metricsInJson, "json", "j", false, "output in json")
}

func getMetricsFormat() string {
	if metricsInJson {
		return "json"
	}
	return "text"
}