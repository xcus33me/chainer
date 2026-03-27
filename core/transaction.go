package core

import (
	"crypto/sha256"

	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/proto"
	pb "google.golang.org/protobuf/proto"
)

func SignTransaction(pk *crypto.PrivateKey, tx *proto.Transaction) *crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		sig := crypto.SignatureFromBytes(input.Signature)
		pubKey := crypto.PublicKeyFromBytes(input.PublicKey)

		// We should make signature nil because we dont have it for now
		savedSig := input.Signature
		input.Signature = nil
		valid := sig.Verify(pubKey, HashTransaction(tx))
		input.Signature = savedSig

		if !valid {
			return false
		}
	}

	return true
}
