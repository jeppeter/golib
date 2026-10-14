package syncpkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSyncAdd(t *testing.T) {
	t.Logf("TestSyncAdd")
	assert.Equal(t, SyncAdd(3, 1), 4)
}
