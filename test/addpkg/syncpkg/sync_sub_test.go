package syncpkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func SyncTestSub(t *testing.T) {
	t.Logf("SyncTestSub")
	assert.Equal(t, SyncSub(3, 1), 2)
}
