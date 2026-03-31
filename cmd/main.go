package main

import (
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/xcus33me/chainer/node"
	"github.com/xcus33me/chainer/proto"
	"github.com/xcus33me/chainer/utils"
	"google.golang.org/grpc"
)

func main() {
	node := node.NewNode()

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		slog.Error("Failed to listen on :3000", "err", err)
		os.Exit(1)
	}
	proto.RegisterNodeServer(grpcServer, node)

	slog.Info("node running on port :3000")

	go func() {
		for {
			time.Sleep(time.Second)
			utils.CreateTransaction()
		}
	}()
	grpcServer.Serve(ln)
}
