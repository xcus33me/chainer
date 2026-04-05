package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privateKey := GeneratePrivateKey()
	assert.Equal(t, PrivateKeyLen, len(privateKey.Bytes()))

	publicKey := privateKey.Public()
	assert.Equal(t, PublicKeyLen, len(publicKey.Bytes()))

}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "87f20c7321d0c670911c172f27ba026b46ecd7eaaa16018e7d0cc71a718d2866"
		privKey    = NewPrivateKeyFromString(seed)
		addressStr = "7afa1c0fd4516f61585d260f5e603284fa131d53"
	)
	assert.Equal(t, PrivateKeyLen, len(privKey.Bytes()))

	address := privKey.Public().Address()
	assert.Equal(t, addressStr, address.String())
}

func TestPrivateKeySign(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.Public()
	msg := []byte("googi")

	sig := privateKey.Sign(msg)
	assert.True(t, sig.Verify(publicKey, msg))

	// Test with invalid msg
	assert.False(t, sig.Verify(publicKey, []byte("foo")))

	// Test with invalid pubKey
	invalidPrivateKey := GeneratePrivateKey()
	invalidPublicKey := invalidPrivateKey.Public()
	assert.False(t, sig.Verify(invalidPublicKey, msg))
}

func TestPublicKeyToAddress(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.Public()
	address := pubKey.Address()

	assert.Equal(t, AddressLen, len(address.Bytes()))
}
