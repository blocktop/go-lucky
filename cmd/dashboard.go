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
	"path"

	dashboard "github.com/blocktop/go-dashboard"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dashboardCmd represents the dashboard command
var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launches the dashboard UI server.",
	Long: `Usage: lucky dashboard [OPTIONS]

By default the dashboard is accessible at http://localhost:3000/`,
	Run: func(cmd *cobra.Command, args []string) {
		err := dashboard.Start()
		if err != nil {
			failWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(dashboardCmd)

	dashboardCmd.PersistentFlags().String("dashHost", "localhost", `host for dashboard UI server, set to 0.0.0.0 to
expose publicly`)
	dashboardCmd.PersistentFlags().Int("dashPort", 3000, "dashboard UI server port")
	dashboardCmd.PersistentFlags().String("dashViews", "", "the directory of the static client-side files for the dashboard UI")

	viper.BindPFlag("dashboard.host", dashboardCmd.PersistentFlags().Lookup("dashHost"))
	viper.BindPFlag("dashboard.port", dashboardCmd.PersistentFlags().Lookup("dashPort"))
	viper.BindPFlag("dashboard.viewsDir", dashboardCmd.PersistentFlags().Lookup("dashViews"))

	viper.BindEnv("dashboard.viewsDir", "LUCKY_DASHBOARD_VIEWS_DIR")

	viper.SetDefault("dashboard.host", "localhost")
	viper.SetDefault("dashboard.port", 3000)

	homeDir := getHomeDir()
	viper.SetDefault("dashboard.viewsDir", path.Join(homeDir, "go", "src", "github.com", "blocktop", "go-dashboard", "views"))
}
