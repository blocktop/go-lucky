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
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"syscall"
	"time"

	blockchain "github.com/blocktop/go-blockchain"
	consensus "github.com/blocktop/go-consensus"
	ctrl "github.com/blocktop/go-controller"
	luckyblock "github.com/blocktop/go-luckyblock"
	p2p "github.com/blocktop/go-network-libp2p"
	rpc "github.com/blocktop/go-rpc-server"
	"github.com/blocktop/go-spec"
	"github.com/golang/glog"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// blockchainCmd represents the blockchain command
var blockchainCmd = &cobra.Command{
	Use:   "blockchain",
	Short: "Starts the lucky blockchain",
	Long:  `Usage: lucky blockchain [OPTIONS]`,
	Run: func(cmd *cobra.Command, args []string) {
		defer glog.Flush()

		if !fileExists(viper.ConfigFileUsed()) {
			failWithError(errors.New("Config file not found. Use the init command to create it."))
		}

		node := buildNode()
		consensus := buildConsensus()
		controller := buildController(node, consensus)

		// TODO temp
		flag.Set("logtostderr", "true")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		pfile := viper.GetString("diagnostics.cpuprofile")
		if pfile != "" {
			f, err := os.Create(pfile)
			if err != nil {
				failWithError(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}

		err := node.Bootstrap(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		node.Listen(ctx)
		controller.Start(ctx)
		rpc.Start()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		<-sig

		controller.Stop()
		node.Close()
	},
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rootCmd.AddCommand(blockchainCmd)

	homeDir := getHomeDir()

	flags := blockchainCmd.PersistentFlags()
	flags.IntP("p2pport", "p", 29190, "port for P2P network listener")
	flags.String("dataDir", "", "directory for blockchain database")
	flags.BoolP("genesis", "g", false, "produce the genesis block")
	flags.StringArrayP("bootstrapPeer", "b", []string{}, `address of peer to bootstrap with, may be specified
more than once`)
	flags.Bool("nobootstrap", false, "disable bootstrapping")
	flags.Bool("nodiscovery", false, "disable peer discovery")
	flags.Bool("trackall", false, `include immediately disqualified blocks in consensus metrics`)
	flags.String("cpuprofile", "", "output file for CPU profile info")

	viper.BindEnv("blockchain.dataDir", "LUCKY_DATA_DIR")
	viper.BindEnv("node.port", "LUCKY_P2P_PORT")

	viper.BindPFlag("node.port", flags.Lookup("p2pport"))
	viper.BindPFlag("node.bootstrapper.peers", flags.Lookup("bootstrapPeer"))
	viper.BindPFlag("node.bootstrapper.disable", flags.Lookup("nobootstrap"))
	viper.BindPFlag("store.dataDir", flags.Lookup("dataDir"))
	viper.BindPFlag("blockchain.genesis", flags.Lookup("genesis"))
	viper.BindPFlag("blockchain.metrics.trackall", flags.Lookup("trackall"))
	viper.BindPFlag("diagnostics.cpuprofile", flags.Lookup("cpuprofile"))

	viper.SetDefault("blockchain.dataDir", path.Join(homeDir, ".lucky", "data"))
	viper.SetDefault("blockchain.genesis", false)
	viper.SetDefault("blockchain.blockInterval", 3*time.Second)
	viper.SetDefault("blockchain.type", "luckychain")
	viper.SetDefault("blockchain.block.type", "luckyblock")
	viper.SetDefault("blockchain.block.version", "v1")
	viper.SetDefault("blockchain.consensus.depth", 10)
	viper.SetDefault("blockchain.consensus.depthBuffer", 2)

	viper.SetDefault("node.bootstrapper.disable", false)
	viper.SetDefault("node.bootstrapper.checkInterval", 5)         // seconds
	viper.SetDefault("node.bootstrapper.rebootstrapInterval", 300) // seconds
	viper.SetDefault("node.bootstrapper.minPeers", 1)
	viper.SetDefault("node.bootstrapper.peers", []string{})
	viper.SetDefault("node.port", 29190)
	viper.SetDefault("node.addresses", []string{"/ip4/0.0.0.0/tcp/29190", "/ip6/::/tcp/29190"})
	viper.SetDefault("node.discovery.disable", false)
	viper.SetDefault("node.discovery.interval", 5) // second
	viper.SetDefault("store.ipfs.apiport", 5001)
	viper.SetDefault("store.ipfs.gatewayport", 8081)
	viper.SetDefault("store.ipfs.swarmport", 4001)
	viper.SetDefault("store.ipfs.swarmhosts", []string{"/ip4/0.0.0.0/tcp", "/ip6/::/tcp"})
	viper.SetDefault("store.ipfs.bootstraplist", []string{}) //TODO
	viper.SetDefault("store.ipfs.pin", false)
	viper.SetDefault("store.ipfs.disablenat", false)
}

func buildConsensus() spec.Consensus {
	blockComparator := luckyblock.BlockComparator
	consensus := consensus.New(10, blockComparator)

	return consensus
}

func buildController(node spec.NetworkNode, consensus spec.Consensus) spec.Controller {
	controller := ctrl.NewController(node)

	blockGenerator := luckyblock.NewBlockGenerator(node.GetPeerID())
	bc := blockchain.NewBlockchain(blockGenerator, consensus, node.GetPeerID())

	controller.AddBlockchain(bc)

	return controller
}

func buildNode() *p2p.NetworkNode {
	addresses := viper.GetStringSlice("node.addresses")
	port := viper.GetInt("port")
	for i, a := range addresses {
		lastSlash := strings.LastIndex(a, "/")
		if lastSlash > -1 {
			aportS := a[lastSlash+1:]
			aport, err := strconv.ParseInt(aportS, 10, 32)
			if err != nil || aport == int64(port) {
				continue
			}
			addresses[i] = a[:lastSlash+1] + strconv.FormatInt(aport, 10)
		}
	}
	viper.Set("node.addresses", addresses)

	node, err := p2p.NewNode()
	if err != nil {
		failWithError(err)
	}

	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", node.GetPeerID()))
	addr := node.Host.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	fmt.Fprintf(os.Stderr, "P2P address: %s\n", fullAddr)

	return node
}
