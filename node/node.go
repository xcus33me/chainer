package node

import (
	"context"
	"net"
	"sync"

	"github.com/xcus33me/chainer/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type Node struct {
	proto.UnimplementedNodeServer

	peers map[proto.NodeClient]*proto.Version
	m     sync.RWMutex

	listenAddr string
	version    string

	logger *zap.SugaredLogger
}

func NewNode() *Node {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()

	return &Node{
		peers:   make(map[proto.NodeClient]*proto.Version),
		version: "chainer-0.1",
		logger:  logger.Sugar(),
	}
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.m.Lock()
	defer n.m.Unlock()

	n.peers[c] = v

	for _, addr := range v.PeerList {
		if addr != n.listenAddr {
			n.logger.Debugf("[%s] need to connect with %s", n.listenAddr, addr)
		}
	}

	n.logger.Debugw("new peer connected", "addr", v.ListenAddr, "height", v.Height)
}

func (n *Node) removePeer(c proto.NodeClient) {
	n.m.Lock()
	defer n.m.Unlock()
	delete(n.peers, c)
}

func (n *Node) BootstrapNetwork(ctx context.Context, addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}

		v, err := c.Handshake(ctx, n.getVersion())
		if err != nil {
			n.logger.Error("handshake error:", err)
			continue
		}

		n.addPeer(c, v)
	}

	return nil
}

func (n *Node) Start(listenAddr string) error {
	n.listenAddr = listenAddr

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Infow("node running...", "port", n.listenAddr)

	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(c, v)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	n.logger.Debug("received tx from", peer.Addr)
	return &proto.Ack{}, nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "chainer-0.1",
		Height:     0,
		ListenAddr: n.listenAddr,
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) getPeerList() []string {
	n.m.RLock()
	defer n.m.RUnlock()

	peers := []string{}
	for _, version := range n.peers {
		peers = append(peers, version.ListenAddr)
	}

	return peers
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	c, err := grpc.Dial(listenAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(c), nil
}
