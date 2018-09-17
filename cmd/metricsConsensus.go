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

	rpcconsensus "github.com/blocktop/go-rpc-client/consensus"
	"github.com/spf13/cobra"
)

// consensusCmd represents the consensus command
var metricsConsensusCmd = &cobra.Command{
	Use:   "consensus",
	Short: "Retrieves metrics from the consensus system of the lucky blockchain.",
	Long:  `Usage: lucky metrics consensus [OPTIONS]`,
	Run: func(cmd *cobra.Command, args []string) {
		// https://gist.github.com/rnix/fc03d74ec128cb6a3099
		res, err := rpcconsensus.GetMetrics(getMetricsFormat())
		if err != nil {
			failWithError(err)
		}
		fmt.Println(res.Metrics)
		/*
			req := &consensus.GetMetricsRequest{}
			reqb, err := json.Marshal(req.GetMetrics(consensus.GetMetricsArgs{metricsFormat}))
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
			var data consensus.GetMetricsResponse
			err = json.Unmarshal(resS, &data)
			if err != nil {
				failWithError(err)
			}
			fmt.Println(data.Result.Metrics)
		*/
	},
}

func init() {
	metricsCmd.AddCommand(metricsConsensusCmd)
}
