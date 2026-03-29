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
	txCopy := &proto.Transaction{
		Version: tx.Version,
		Inputs:  make([]*proto.TxInput, len(tx.Inputs)),
		Outputs: tx.Outputs,
	}

	for i, input := range tx.Inputs {
		txCopy.Inputs[i] = &proto.TxInput{
			PrevTxHash: input.PrevTxHash,
			PrevOutIdx: input.PrevOutIdx,
			PublicKey:  input.PublicKey,
			// Do not copy signature
		}
	}

	b, err := pb.Marshal(txCopy)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)
	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	txHash := HashTransaction(tx)
	for _, input := range tx.Inputs {
		sig, err := crypto.SignatureFromBytes(input.Signature)
		if err != nil {
			return false
		}

		pubKey, err := crypto.PublicKeyFromBytes(input.PublicKey)
		if err != nil {
			return false
		}

		if !sig.Verify(pubKey, txHash) {
			return false
		}
	}

	return true
}
