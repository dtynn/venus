package main

import (
	"log"
	"path/filepath"

	gen "github.com/whyrusleeping/cbor-gen"

	"github.com/filecoin-project/venus/venus-shared/chain"
	"github.com/filecoin-project/venus/venus-shared/libp2p/exchange"
	"github.com/filecoin-project/venus/venus-shared/libp2p/hello"
)

type genTarget struct {
	dir   string
	pkg   string
	types []interface{}
}

func main() {
	targets := []genTarget{
		{
			dir: "../venus-shared/libp2p/hello/",
			types: []interface{}{
				hello.GreetingMessage{},
				hello.LatencyMessage{},
			},
		},
		{
			dir: "../venus-shared/libp2p/exchange/",
			types: []interface{}{
				exchange.Request{},
				exchange.Response{},
				exchange.CompactedMessages{},
				exchange.BSTipSet{},
			},
		},
		{
			dir: "../venus-shared/chain/",
			types: []interface{}{
				chain.BlockHeader{},
				chain.Ticket{},
				chain.ElectionProof{},
				chain.BeaconEntry{},
				chain.Message{},
				chain.SignedMessage{},
				chain.Actor{},
				chain.MessageRoot{},
				chain.MessageReceipt{},
				chain.BlockMsg{},
				chain.ExpTipSet{},
			},
		},
	}

	for _, target := range targets {
		pkg := target.pkg
		if pkg == "" {
			pkg = filepath.Base(target.dir)
		}

		if err := gen.WriteTupleEncodersToFile(filepath.Join(target.dir, "cbor_gen.go"), pkg, target.types...); err != nil {
			log.Fatalf("gen for %s: %s", target.dir, err)
		}
	}
}