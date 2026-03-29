package node

import (
	"context"
	"log/slog"

	"github.com/xcus33me/chainer/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.None, error) {
	peer, _ := peer.FromContext(ctx)
	slog.Info("node - HandleTransaction - received tx", "peer", peer)
	return nil, nil
}
