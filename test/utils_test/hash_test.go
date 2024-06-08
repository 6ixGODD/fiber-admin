package utils_test

import (
	"testing"

	"fiber-admin/pkg/utils/crypt"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	hash, err := crypt.Hash("foo")
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	assert.True(t, crypt.Compare("foo", hash))
	assert.False(t, crypt.Compare("bar", hash))

	pwd, err := crypt.Hash("Admin@123")
	assert.NoError(t, err)
	assert.NotEmpty(t, pwd)

	t.Logf("Hash: %s", pwd)
}
