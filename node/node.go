package node

import (
	"context"
	"log/slog"

	"github.com/xcus33me/chainer/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	proto.UnimplementedNodeServer

	verison string
}

func NewNode() *Node {
	return &Node{
		verison: "chainer-0.1",
	}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	slog.Info("node - HandleTransaction - received tx", "peer", peer)
	return &proto.Ack{}, nil
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.verison,
		Height:  0,
	}

	p, _ := peer.FromContext(ctx)
	slog.Info("received version from", "node_addr", p.Addr)

	return ourVersion, nil
}
