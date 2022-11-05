// Copyright (c) 2018-2022 Dave Collins
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/decred/dcrd/blockchain/stake/v4"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/chaincfg/v3"
	"github.com/decred/dcrd/dcrutil/v4"
	"github.com/decred/dcrd/rpcclient/v7"
)

// ticketData houses information about a purchased ticket.
type ticketData struct {
	minedHeight int64
	ticketPrice dcrutil.Amount
}

// reportProgress periodically prints out the current height to stdout.
func reportProgress(height int64) {
	if height%10000 == 0 && height != 0 {
		fmt.Println()
	}
	if height%1000 == 0 && height != 0 {
		fmt.Printf("..%d", height)
	}
}

func main() {
	// Define and parse command line flags.
	dcrdHomeDir := dcrutil.AppDataDir("dcrd", false)
	var rpcServer = flag.String("rpcserver", "localhost:9109",
		"RPC server address")
	var rpcUser = flag.String("rpcuser", "", "RPC server username")
	var rpcPass = flag.String("rpcpass", "", "RPC server passphrase")
	var rpcCert = flag.String("rpccert", filepath.Join(dcrdHomeDir, "rpc.cert"),
		"RPC server TLS certificate")
	var verbose = flag.Bool("verbose", false, "Print details about every vote")
	var version = flag.Bool("version", false, "Display version information and exit")
	flag.Parse()

	// Show the version and exit if the version flag was specified.
	if *version {
		appName := filepath.Base(os.Args[0])
		appName = strings.TrimSuffix(appName, filepath.Ext(appName))
		fmt.Printf("%s version %s (Go version %s %s/%s)\n", appName,
			Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// Connect to dcrd RPC server using websockets.
	certs, err := ioutil.ReadFile(*rpcCert)
	if err != nil {
		fmt.Println("Unable to load RPC TLS cert:", err)
		return
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         *rpcServer,
		Endpoint:     "ws",
		User:         *rpcUser,
		Pass:         *rpcPass,
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	params := chaincfg.MainNetParams()

	// Find the best block height so the data is calculated for all blocks.
	ctx := context.Background()
	_, blockHeight, err := client.GetBestBlock(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Calculating average vote time through block height %d...\n",
		blockHeight)
	if !*verbose {
		fmt.Printf("Height")
	}

	var totalVotes, totalWaitBlocks int64
	var totalWaitSeconds float64
	blockTimes := make(map[int64]time.Time)
	tickets := make(map[chainhash.Hash]ticketData)
	for i := int64(1); i <= blockHeight; i++ {
		if !*verbose {
			reportProgress(i)
		}

		// Load the block for the height and store its timestamp for later use.
		hash, err := client.GetBlockHash(ctx, i)
		if err != nil {
			fmt.Println(err)
			return
		}
		blk, err := client.GetBlock(ctx, hash)
		if err != nil {
			fmt.Println(err)
			return
		}
		blockTimes[i] = blk.Header.Timestamp

		for _, stx := range blk.STransactions {
			switch stake.DetermineTxType(stx, true, true) {

			// Track ticket purchases.
			case stake.TxTypeSStx:
				tickets[stx.TxHash()] = ticketData{
					minedHeight: i,
					ticketPrice: dcrutil.Amount(stx.TxOut[0].Value),
				}

			// Calculate and track wait times for votes.
			case stake.TxTypeSSGen:
				// Lookup the ticket associated with the vote.
				ticketHash := &stx.TxIn[1].PreviousOutPoint.Hash
				ticket, ok := tickets[*ticketHash]
				if !ok {
					fmt.Printf("Ticket %s not found\n", ticketHash)
					return
				}

				// Calculate the wait time based on when the ticket matured and
				// when it voted.
				maturityHeight := ticket.minedHeight + int64(params.TicketMaturity) + 1
				voteWaitBlocks := (i - maturityHeight) + 1
				voteWaitTime := blk.Header.Timestamp.Sub(blockTimes[maturityHeight])

				if *verbose {
					voteWaitDays := voteWaitTime.Hours() / 24.0
					fmt.Printf("Ticket %s... (%v) mined in block %d, voted %d "+
						"blocks (%.2f days) after maturity\n",
						ticketHash.String()[:8], ticket.ticketPrice,
						ticket.minedHeight, voteWaitBlocks, voteWaitDays)
				}

				totalVotes++
				totalWaitBlocks += voteWaitBlocks
				totalWaitSeconds += voteWaitTime.Seconds()

				delete(tickets, *ticketHash)
			}
		}
	}

	if !*verbose {
		fmt.Println("..done")
	}

	avgWaitBlocks := float64(totalWaitBlocks) / float64(totalVotes)
	avgWaitSeconds := totalWaitSeconds / float64(totalVotes)
	fmt.Printf("Mean wait for %d votes: %.1f blocks, %.2f days\n", totalVotes,
		avgWaitBlocks, avgWaitSeconds/86400.0)

	client.Shutdown()
	client.WaitForShutdown()
}
