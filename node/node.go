package node

import (
	"context"
	"encoding/hex"
	"net"
	"sync"
	"time"

	"github.com/xcus33me/chainer/core"
	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

const (
	blockTime = time.Second * 5
)

type Mempool struct {
	txx map[string]*proto.Transaction
}

func NewMempool() *Mempool {
	return &Mempool{
		txx: make(map[string]*proto.Transaction),
	}
}

func (p *Mempool) Has(tx *proto.Transaction) bool {
	hash := hex.EncodeToString(core.HashTransaction(tx))
	_, ok := p.txx[hash]
	return ok
}

func (p *Mempool) Add(tx *proto.Transaction) bool {
	if p.Has(tx) {
		return false
	}

	hash := hex.EncodeToString(core.HashTransaction(tx))
	p.txx[hash] = tx
	return true
}

type ServerConfig struct {
	Version    string
	ListenAddr string
	PrivateKey *crypto.PrivateKey
}

type Node struct {
	proto.UnimplementedNodeServer

	ServerConfig

	mempool *Mempool
	peers   map[proto.NodeClient]*proto.Version
	m       sync.RWMutex

	logger *zap.SugaredLogger
}

func NewNode(cfg ServerConfig) *Node {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = ""
	logger, _ := loggerConfig.Build()

	return &Node{
		peers:        make(map[proto.NodeClient]*proto.Version),
		logger:       logger.Sugar(),
		mempool:      NewMempool(),
		ServerConfig: cfg,
	}
}

func (n *Node) Start(listenAddr string, bootstrapNodes []string) error {
	n.ListenAddr = listenAddr

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	proto.RegisterNodeServer(grpcServer, n)

	n.logger.Infow("node running...", "port", n.ListenAddr)

	// bootstrap the network with a list of already known nodes
	// in the network.
	if len(bootstrapNodes) > 0 {
		go n.bootstrapNetwork(context.Background(), bootstrapNodes)
	}

	if n.PrivateKey != nil {
		go n.validatorLoop()
	}

	return grpcServer.Serve(ln)
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}

	n.addPeer(ctx, c, v)

	return n.getVersion(), nil
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, _ := peer.FromContext(ctx)
	hash := hex.EncodeToString(core.HashTransaction(tx))

	if n.mempool.Add(tx) {
		n.logger.Debugw("received tx", "we", n.ListenAddr, "from", peer.Addr, "hash", hash)
		go func() {
			if err := n.broadcast(tx); err != nil {
				n.logger.Errorw("broadcast error", "err", err)
			}
		}()
	}

	return &proto.Ack{}, nil
}

func (n *Node) validatorLoop() {
	n.logger.Infow("starting validator loop", "pubKey", n.PrivateKey.Public(), "blockTime", blockTime)
	ticker := time.NewTicker(blockTime)

	for {
		<-ticker.C
		n.logger.Debugw("time to create a new block", "txLen", len(n.mempool.txx))
	}
}

func (n *Node) broadcast(msg any) error {
	for peer := range n.peers {
		switch v := msg.(type) {
		case *proto.Transaction:
			_, err := peer.HandleTransaction(context.Background(), v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *Node) addPeer(ctx context.Context, c proto.NodeClient, v *proto.Version) {
	n.m.Lock()
	defer n.m.Unlock()

	// TODO: Handle logic where we decide to accept or drop the incoming node connection

	n.peers[c] = v

	// Connect to all peers in the received list of peers
	if len(v.PeerList) > 0 {
		go n.bootstrapNetwork(ctx, v.PeerList)
	}

	n.logger.Debugw("new peer successfully connected",
		"we", n.ListenAddr, "remote", v.ListenAddr, "height", v.Height)
}

func (n *Node) removePeer(c proto.NodeClient) {
	n.m.Lock()
	defer n.m.Unlock()
	delete(n.peers, c)
}

func (n *Node) bootstrapNetwork(ctx context.Context, addrs []string) error {
	for _, addr := range addrs {
		if !n.canConnectWith(addr) {
			continue
		}

		n.logger.Debugw("dialing remote node", "we", n.ListenAddr, "remote", addr)

		c, v, err := n.dialRemoteNode(ctx, addr)
		if err != nil {
			return err
		}
		n.addPeer(ctx, c, v)
	}

	return nil
}

func (n *Node) dialRemoteNode(ctx context.Context, addr string) (proto.NodeClient, *proto.Version, error) {
	c, err := makeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}

	v, err := c.Handshake(ctx, n.getVersion())

	if err != nil {
		return nil, nil, err
	}

	return c, v, nil
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "chainer-0.1",
		Height:     0,
		ListenAddr: n.ListenAddr,
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) canConnectWith(addr string) bool {
	if n.ListenAddr == addr {
		return false
	}

	connectedPeers := n.getPeerList()
	for _, connectedAddr := range connectedPeers {
		if addr == connectedAddr {
			return false
		}
	}

	return true
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
