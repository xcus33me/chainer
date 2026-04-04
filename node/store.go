package node

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/xcus33me/chainer/core"
	"github.com/xcus33me/chainer/proto"
)

type BlockStorer interface {
	Put(*proto.Block) error
	Get(string) (*proto.Block, error)
}

type MemoryBlockStore struct {
	blocks map[string]*proto.Block
	mu     sync.RWMutex
}

func NewMemoryBlockStore() *MemoryBlockStore {
	return &MemoryBlockStore{
		blocks: make(map[string]*proto.Block),
	}
}

func (s *MemoryBlockStore) Put(block *proto.Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash := hex.EncodeToString(core.HashBlock(block))
	s.blocks[hash] = block
	return nil
}

func (s *MemoryBlockStore) Get(hash string) (*proto.Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	block, ok := s.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash [%s] does not exist", hash)
	}
	return block, nil
}
