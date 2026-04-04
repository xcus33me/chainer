package main

import (
	"context"
	"log"
	"time"

	"github.com/xcus33me/chainer/node"
)

func main() {
	makeNode(":3000", []string{})
	time.Sleep(1 * time.Second)
	makeNode(":4000", []string{":3000"})
	time.Sleep(4 * time.Second)
	makeNode(":5000", []string{":4000"})
}

func makeNode(listenAddr string, bootstrapNodes []string) *node.Node {
	n := node.NewNode()
	go n.Start(listenAddr)

	if len(bootstrapNodes) >= 1 {
		if err := n.BootstrapNetwork(context.Background(), bootstrapNodes); err != nil {
			log.Fatal(err)
		}
	}

	return n
}
