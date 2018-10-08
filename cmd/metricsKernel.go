// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	rpckernel "github.com/blocktop/go-rpc-client/kernel"

)

// metricsKernelCmd represents the kernel command
var metricsKernelCmd = &cobra.Command{
	Use:   "kernel",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := rpckernel.GetMetrics(getMetricsFormat())
		if err != nil {
			failWithError(err)
		}
		fmt.Println(res.Metrics)
	},
}

func init() {
	metricsCmd.AddCommand(metricsKernelCmd)
}
