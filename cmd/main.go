package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/node"
	"github.com/xcus33me/chainer/proto"
	"github.com/xcus33me/chainer/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	makeNode(":3000", []string{}, true)
	time.Sleep(1 * time.Second)
	makeNode(":4000", []string{":3000"}, false)
	time.Sleep(1 * time.Second)
	makeNode(":5000", []string{":4000"}, false)
	time.Sleep(time.Second)

	for {
		time.Sleep(time.Millisecond * 100)
		makeTransaction()
	}
	// select {}
}

func makeNode(listenAddr string, bootstrapNodes []string, isValidator bool) *node.Node {
	cfg := node.ServerConfig{
		Version:    "chainer-0.1",
		ListenAddr: listenAddr,
	}

	if isValidator {
		cfg.PrivateKey = crypto.GeneratePrivateKey()
	}

	n := node.NewNode(cfg)
	go n.Start(listenAddr, bootstrapNodes)
	return n
}

func makeTransaction() {
	client, err := grpc.Dial(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("utils - CreateTransaction - grpc.NewClient", "err", err)
		os.Exit(1)
	}

	c := proto.NewNodeClient(client)
	privKey := crypto.GeneratePrivateKey()

	tx := &proto.Transaction{
		Version: 1,
		Inputs: []*proto.TxInput{
			{
				PrevTxHash: utils.RandomHash(),
				PrevOutIdx: 0,
				PublicKey:  privKey.Public().Bytes(),
			},
		},
		Outputs: []*proto.TxOutput{
			{
				Amount:  99,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}

	_, err = c.HandleTransaction(context.TODO(), tx)
	if err != nil {
		slog.Error("utils - CreateTransaction", "err", err)
	}
}
