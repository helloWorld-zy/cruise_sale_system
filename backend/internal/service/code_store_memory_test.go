package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeStoreMemory(t *testing.T) {
	store := NewInMemoryCodeStore()
	err := store.Save("phone", "1234")
	assert.NoError(t, err)

	ok := store.Verify("phone", "1234")
	assert.True(t, ok)

	ok = store.Verify("phone", "wrong")
	assert.False(t, ok)
}
