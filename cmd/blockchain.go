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

	"github.com/blocktop/go-kernel"

	blockchain "github.com/blocktop/go-blockchain"
	consensus "github.com/blocktop/go-consensus"
	luckyblock "github.com/blocktop/go-luckyblock"
	p2p "github.com/blocktop/go-network-libp2p"
	rpc "github.com/blocktop/go-rpc-server"
	spec "github.com/blocktop/go-spec"

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
			failWithError(errors.New("config file not found. Use the init command to create it"))
		}

		node := buildNode()
		cons := buildConsensus()
		bg := buildBlockGenerator()
		bc := buildBlockchain(cons, bg)

		cfg := &kernel.KernelConfig{
			Blockchain:     bc,
			Consensus:      cons,
			BlockFrequency: viper.GetFloat64("blockchain.blockFrequency"),
			BlockPrototype: bg.BlockPrototype(),
			NetworkNode:    node}

		kernel.Init(cfg)

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
		bc.Start(ctx)
		rpc.Start()
		kernel.Start(ctx)

		sig := make(chan os.Signal, 1)
		signal.Notify(sig,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		<-sig

		kernel.Stop()
		bc.Stop()
		node.Close()
	},
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rootCmd.AddCommand(blockchainCmd)

	homeDir := getHomeDir()

	flags := blockchainCmd.Flags()
	flags.IntP("p2pport", "p", 29190, "port for P2P network listener")
	flags.String("dataDir", "", "directory for blockchain database")
	flags.BoolP("genesis", "g", false, "produce the genesis block")
	flags.StringArrayP("bootstrapPeer", "b", []string{}, `address of peer to bootstrap with, may be specified
more than once`)
	flags.Bool("nobootstrap", false, "disable bootstrapping")
	flags.Bool("nodiscovery", false, "disable peer discovery")
	flags.Bool("trackall", false, `include immediately disqualified blocks in consensus metrics`)
	flags.String("cpuprofile", "", "output file for CPU profile info")
	flags.Float64P("blockFrequency", "f", 1.0, "Number of blocks per second. Can be a decimal number.")
	flags.DurationP("consensusTime", "t", 30*time.Second, "The duration blocks are tracked before consensus is reached.")
	viper.BindEnv("blockchain.dataDir", "LUCKY_DATA_DIR")
	viper.BindEnv("node.port", "LUCKY_P2P_PORT")

	viper.BindPFlag("node.port", flags.Lookup("p2pport"))
	viper.BindPFlag("node.bootstrapper.peers", flags.Lookup("bootstrapPeer"))
	viper.BindPFlag("node.bootstrapper.disable", flags.Lookup("nobootstrap"))
	viper.BindPFlag("store.dataDir", flags.Lookup("dataDir"))
	viper.BindPFlag("blockchain.genesis", flags.Lookup("genesis"))
	viper.BindPFlag("blockchain.metrics.trackall", flags.Lookup("trackall"))
	viper.BindPFlag("blockchain.blockFrequency", flags.Lookup("blockFrequency"))
	viper.BindPFlag("blockchain.consensus.time", flags.Lookup("consensusTime"))
	viper.BindPFlag("diagnostics.cpuprofile", flags.Lookup("cpuprofile"))

	viper.SetDefault("blockchain.dataDir", path.Join(homeDir, ".lucky", "data"))
	viper.SetDefault("blockchain.genesis", false)
	viper.SetDefault("blockchain.blockFrequency", float64(1)) // blocks/second
	viper.SetDefault("blockchain.name", "luckychain")
	viper.SetDefault("blockchain.block.name", "luckyblock")
	viper.SetDefault("blockchain.block.namespace", "io.blocktop.lucky")
	viper.SetDefault("blockchain.block.version", "v1")
	viper.SetDefault("blockchain.consensus.time", 30*time.Second)
	viper.SetDefault("blockchain.receiveconcurrency", 2)

	viper.SetDefault("node.bootstrapper.disable", false)
	viper.SetDefault("node.bootstrapper.checkInterval", 5)         // seconds
	viper.SetDefault("node.bootstrapper.rebootstrapInterval", 300) // seconds
	viper.SetDefault("node.bootstrapper.minPeers", 1)
	viper.SetDefault("node.bootstrapper.peers", []string{
		"/ip4/104.196.155.69/tcp/29190/ipfs/QmTTDpNa8ErE23Fs3YZFLnprv6UaXTWFsm11Tt2zcWgKBJ",
		"/ip4/35.204.208.27/tcp/29190/ipfs/QmdKoGtMGzeqeZ9M1zt4zE5YsRRdH3h2b6oPKW9pmvb3Xc",
		"/ip4/35.200.229.227/tcp/29190/ipfs/QmUCx8w8YjnhMLARdjHfTjf4S1DMqB5PhW2CUGHcDeMD4S"})
	viper.SetDefault("node.port", 29190)
	viper.SetDefault("node.addresses", []string{"/ip4/0.0.0.0/tcp/29190", "/ip6/::/tcp/29190"})
	viper.SetDefault("node.discovery.disable", false)
	viper.SetDefault("node.discovery.interval", 5) // second
	viper.SetDefault("node.broadcastconcurrency", 4)
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
	cons := consensus.NewConsensus(blockComparator)

	return cons
}

func buildBlockGenerator() spec.BlockGenerator {
	blockGenerator := luckyblock.NewBlockGenerator()
	return blockGenerator
}

func buildBlockchain(cons spec.Consensus, bg spec.BlockGenerator) spec.Blockchain {
	bc := blockchain.NewBlockchain(bg, cons)
	return bc
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

	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", node.PeerID()))
	addr := node.Host.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	fmt.Fprintf(os.Stderr, "P2P address: %s\n", fullAddr)

	return node
}
