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
	"flag"
	"fmt"
	"os"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var defaultCfgFile string
var homeDir string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.0.1",
	Use:     "lucky",
	Short:   "The 'Lucky' example blockchain application",
	Long: `Lucky is an example blockchain application that showcases the blocktop
	blockchain development kit.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	homeDir := getHomeDir()
	defaultCfgFile = path.Join(homeDir, ".lucky", "config.yaml")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", fmt.Sprintf("config file (default is %s)", defaultCfgFile))

	viper.BindEnv("rpc.port", "LUCKY_RPC_PORT")
	rootCmd.PersistentFlags().IntP("rpcport", "r", 28180, "port for RPC server")
	viper.BindPFlag("rpc.port", rootCmd.PersistentFlags().Lookup("rpcport"))
	viper.SetDefault("rpc.port", 28180)

	flag.Parse()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".lucky" (without extension).
		viper.AddConfigPath(path.Join(homeDir, ".lucky"))
		viper.SetConfigName("config")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	viper.SetEnvPrefix("LUCKY_")
}

func getHomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home
}
