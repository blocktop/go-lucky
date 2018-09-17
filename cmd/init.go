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
	"crypto/rand"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"

	crypto "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates the configuration file required by lucky.",
	Long: `Usage: lucky init [OPTIONS]`,
	Run: func(cmd *cobra.Command, args []string) {
		if cfgFile != "" && fileExists(cfgFile) {
			fmt.Println("Config file already exists:", cfgFile)
			os.Exit(1)
		}

		r := rand.Reader
		priv, pub, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			failWithError(err)
		}
		privKeyBytes, err := crypto.MarshalPrivateKey(priv)
		if err != nil {
			failWithError(err)
		}
		privKeyStr := crypto.ConfigEncodeKey(privKeyBytes)
		viper.Set("node.privateKey", privKeyStr)

		pubKeyBytes, err := crypto.MarshalPublicKey(pub)
		if err != nil {
			failWithError(err)
		}
		pubKeyStr := crypto.ConfigEncodeKey(pubKeyBytes)
		if err != nil {
			failWithError(err)
		}
		viper.Set("node.publicKey", pubKeyStr)

		peerID, err := peer.IDFromPublicKey(pub)
		if err != nil {
			failWithError(err)
		}
		viper.Set("node.peerID", peerID.Pretty())

		port := viper.GetInt("node.port")
		if port != 29190 {
			addresses := []string{
				fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
				fmt.Sprintf("/ip6/::/tcp/%d", port)}

			viper.Set("node.addresses", addresses)
		}

		if cfgFile == "" {
			cfgFile = path.Join(getHomeDir(), ".lucky", "config.yaml")
		}
		lastSlash := strings.LastIndex(cfgFile, string(os.PathSeparator))
		makeDirAll(cfgFile[:lastSlash])
		err = viper.WriteConfigAs(cfgFile)
		if err != nil {
			failWithError(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	flags := initCmd.PersistentFlags()
	flags.AddFlagSet(blockchainCmd.PersistentFlags())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func failWithError(err error) {
	fmt.Println("An error occurred executing the command:")
	fmt.Println(err)
	os.Exit(1)
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func makeDirAll(path string) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		fmt.Println("An error occurred creating the config directory:")
		fmt.Println(err)
		os.Exit(1)
	}
}
