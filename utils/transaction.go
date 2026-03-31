package utils

import (
	"context"
	"log/slog"
	"os"

	"github.com/xcus33me/chainer/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateTransaction() {
	client, err := grpc.NewClient(":3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("utils - CreateTransaction - grpc.NewClient", "err", err)
		os.Exit(1)
	}

	c := proto.NewNodeClient(client)

	version := &proto.Version{
		Version: "chainer-0.1",
		Height:  0,
	}

	_, err = c.Handshake(context.TODO(), version)
	if err != nil {
		slog.Error("utils - CreateTransaction - c.Handshake", "err", err)
	}

	// _, err = c.HandleTransaction(context.TODO(), &proto.Transaction{})
	// // if err != nil {
	// // 	slog.Error("utils - CreateTransaction", "err", err)
	// }
}
