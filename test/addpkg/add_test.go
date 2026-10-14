package addpkg

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Logf("TestAdd")
	assert.Equal(t, Add(3, 1), 4)
}
