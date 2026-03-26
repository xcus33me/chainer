package core

import (
	"crypto/sha256"

	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/proto"
	pb "google.golang.org/protobuf/proto"
)

func SignBlock(privKey *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	return privKey.Sign(HashBlock(b))
}

// HashBlock returns a SHA256 of the header
func HashBlock(block *proto.Block) []byte {
	b, err := pb.Marshal(block.Header)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)

	return hash[:]
}
