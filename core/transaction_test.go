package core

import (
	"fmt"
	"testing"

	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/proto"
	"github.com/xcus33me/chainer/utils"
)

// Balance 100 coins
// want to spend 5 coint to 'user_1'
// 2 ouputs
// 5 to the 'user_1'
// 95 back to our address
func TestNewTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	fromAddress := fromPrivKey.Public().Address().Bytes()

	toPrivKey := crypto.GeneratePrivateKey()
	toAddress := toPrivKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash: utils.RandomHash(),
		PrevOutIdx: 0,
		PublicKey:  toPrivKey.Public().Bytes(),
	}

	output1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}

	output2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output1, output2},
	}

	sig := SignTransaction(fromPrivKey, tx)
	input.Signature = sig.Bytes()

	fmt.Printf("%+v\n", tx)
}
