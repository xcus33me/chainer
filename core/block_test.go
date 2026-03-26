package core

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/utils"
)

func TestSignBlock(t *testing.T) {
	var (
		block   = utils.RandomBlock()
		privKey = crypto.GeneratePrivateKey()
		pubKey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))
}

func TestHashBlock(t *testing.T) {
	block := utils.RandomBlock()
	hash := HashBlock(block)
	fmt.Println(hex.EncodeToString(hash))
	assert.Equal(t, len(hash), 32)
}
