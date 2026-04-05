package core

import (
	"crypto/sha256"

	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/proto"
	pb "google.golang.org/protobuf/proto"
)

func VerifyBlock(b *proto.Block) bool {
	if len(b.PublicKey) != crypto.PublicKeyLen {
		return false
	}

	if len(b.Signature) != crypto.SignatureLen {
		return false
	}

	sig, err := crypto.SignatureFromBytes(b.Signature)
	if err != nil {
		return false
	}

	pubKey, err := crypto.PublicKeyFromBytes(b.PublicKey)
	if err != nil {
		return false
	}

	hash := HashBlock(b)

	return sig.Verify(pubKey, hash)
}

func SignBlock(privKey *crypto.PrivateKey, b *proto.Block) *crypto.Signature {
	sig := privKey.Sign(HashBlock(b))
	b.PublicKey = privKey.Public().Bytes()
	b.Signature = sig.Bytes()

	return sig
}

// HashBlock returns a SHA256 of the header
func HashBlock(block *proto.Block) []byte {
	return HashHeader(block.Header)
}

func HashHeader(header *proto.Header) []byte {
	b, err := pb.Marshal(header)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)

	return hash[:]
}
