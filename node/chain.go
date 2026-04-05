package node

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/xcus33me/chainer/core"
	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/proto"
)

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*proto.Header{},
	}
}

func (l *HeaderList) Add(h *proto.Header) {
	l.headers = append(l.headers, h)
}

func (l *HeaderList) Get(idx int) *proto.Header {
	if idx > l.Height() {
		panic("index too high")
	}

	return l.headers[idx]
}

func (l *HeaderList) Len() int {
	return len(l.headers)
}

func (l *HeaderList) Height() int {
	return l.Len() - 1
}

type Chain struct {
	blockStore BlockStorer
	headers    *HeaderList
}

func NewChain(bs BlockStorer) *Chain {
	chain := &Chain{
		blockStore: bs,
		headers:    NewHeaderList(),
	}
	chain.addBlock(createGenesisBlock())

	return chain
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

func (c *Chain) AddBlock(b *proto.Block) error {
	if err := c.ValidateBlock(b); err != nil {
		return err
	}

	return c.addBlock(b)
}

func (c *Chain) addBlock(b *proto.Block) error {
	// Add the header to the list of headers
	c.headers.Add(b.Header)
	// validation
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	if c.Height() < height {
		return nil, fmt.Errorf("given height (%d) is too high - height (%d)", height, c.Height())
	}

	header := c.headers.Get(height)
	hash := core.HashHeader(header)

	return c.GetBlockByHash(hash)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) ValidateBlock(b *proto.Block) error {
	// Validate the signature of block
	if !core.VerifyBlock(b) {
		return fmt.Errorf("invalid block signature")
	}

	// Validate if the prevHash is the actual hash of the curr block
	currBlock, err := c.GetBlockByHeight(c.Height())
	if err != nil {
		return err
	}

	hash := core.HashBlock(currBlock)
	if !bytes.Equal(hash, b.Header.PrevHash) {
		return fmt.Errorf("invalid prev block hash")
	}

	return nil
}

func createGenesisBlock() *proto.Block {
	privKey := crypto.GeneratePrivateKey()
	block := &proto.Block{
		Header: &proto.Header{
			Version: 1,
		},
	}
	core.SignBlock(privKey, block)

	return block
}
