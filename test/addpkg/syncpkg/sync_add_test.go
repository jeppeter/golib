package syncpkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func SyncTestAdd(t *testing.T) {
	t.Logf("SyncTestAdd")
	assert.Equal(t, SyncAdd(3, 1), 4)
}
