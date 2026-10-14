package addpkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSub(t *testing.T) {
	t.Logf("TestSub")
	assert.Equal(t, Sub(3, 1), 2)
}
