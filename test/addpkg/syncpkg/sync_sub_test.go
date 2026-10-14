package syncpkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSyncSub(t *testing.T) {
	t.Logf("TestSyncSub")
	assert.Equal(t, SyncSub(3, 1), 2)
}
