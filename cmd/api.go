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
	"github.com/blocktop/go-api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// apiCmd represents the dashboard command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Launches the API server.",
	Long: `Usage: lucky api [OPTIONS]

By default the API is accessible at localhost:3000
Monitor API health with
	curl http://localhost:3000/api/`,
	Run: func(cmd *cobra.Command, args []string) {
		err := api.Start()
		if err != nil {
			failWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	apiCmd.PersistentFlags().String("apiHost", "localhost", `host for API server, set to 0.0.0.0 to
expose publicly`)
	apiCmd.PersistentFlags().Int("apiPort", 3000, "API server port")

	viper.BindPFlag("api.host", apiCmd.PersistentFlags().Lookup("apiHost"))
	viper.BindPFlag("api.port", apiCmd.PersistentFlags().Lookup("apiPort"))

	viper.SetDefault("api.host", "localhost")
	viper.SetDefault("api.port", 3000)
}
