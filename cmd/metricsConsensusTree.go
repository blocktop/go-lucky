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
	rpcconsensus "github.com/blocktop/go-rpc-client/consensus"

)

// metricsConsensusTreeCmd represents the tree command
var metricsConsensusTreeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Retrieves the current consensus-finding tree from lucky blockchain.",
	Long:  `Usage: lucky metrics consensus tree [OPTIONS]`,
	Run: func(cmd *cobra.Command, args []string) {
		res, err := rpcconsensus.GetTree(getMetricsFormat())
		if err != nil {
			failWithError(err)
		}
		fmt.Println(res.Tree)
		/*
			req := &consensus.GetTreeRequest{}
			reqb, err := json.Marshal(req.GetTree(consensus.GetTreeArgs{"text"}))
			if err != nil {
				failWithError(err)
			}
			url := fmt.Sprintf("http://localhost:%d/rpc", viper.GetInt("rpc.port"))
			res, err := http.Post(url, "application/json", strings.NewReader(string(reqb)))
			if err != nil {
				failWithError(err)
			}
			resS, err := ioutil.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				failWithError(err)
			}
			var data consensus.GetTreeResponse
			err = json.Unmarshal(resS, &data)
			if err != nil {
				failWithError(err)
			}
			fmt.Println(data.Result.Tree)
		*/
	},
}

func init() {
	metricsConsensusCmd.AddCommand(metricsConsensusTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metricsConsensusTreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metricsConsensusTreeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
