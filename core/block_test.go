package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xcus33me/chainer/crypto"
	"github.com/xcus33me/chainer/utils"
)

func TestSignVerifyBlock(t *testing.T) {
	var (
		block   = utils.RandomBlock()
		privKey = crypto.GeneratePrivateKey()
		pubKey  = privKey.Public()
	)

	sig := SignBlock(privKey, block)
	assert.Equal(t, 64, len(sig.Bytes()))
	assert.True(t, sig.Verify(pubKey, HashBlock(block)))

	assert.Equal(t, block.PublicKey, pubKey.Bytes())
	assert.Equal(t, block.Signature, sig.Bytes())
	assert.True(t, VerifyBlock(block))

	invalidPrivKey := crypto.GeneratePrivateKey()
	block.PublicKey = invalidPrivKey.Public().Bytes()
	assert.False(t, VerifyBlock(block))
}

func TestHashBlock(t *testing.T) {
	block := utils.RandomBlock()
	hash := HashBlock(block)
	assert.Equal(t, len(hash), 32)
}
