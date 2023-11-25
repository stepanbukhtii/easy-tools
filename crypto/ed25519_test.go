package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	publicED25519Key  = "MCowBQYDK2VwAyEAwggfjC0JqoTCejZGFvSQoP50+3wdftcIBNkZU1w4JPg="
	privateED25519Key = "MC4CAQAwBQYDK2VwBCIEIPHEJfftv7kbcWtE+PJvUMdpofbVEO2fsSCBQ2UygOHj"
)

func TestED25519Key(t *testing.T) {
	pubKey, err := ParseED25519PublicKey(publicED25519Key)
	assert.NoError(t, err)
	assert.NotNil(t, pubKey)

	privateKey, err := ParseED25519PrivateKey(privateED25519Key)
	assert.NoError(t, err)
	assert.NotNil(t, privateKey)

	assert.Equal(t, pubKey, privateKey.Public())
}
