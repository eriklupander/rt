package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewBoundingBox(t *testing.T) {
	box := NewEmptyBoundingBox()
	assert.Equal(t, NewTupleOf(math.Inf(1), math.Inf(1),math.Inf(1), 0), box.Min)
	assert.Equal(t, NewTupleOf(math.Inf(-1), math.Inf(-1),math.Inf(-1), 0), box.Max)
}

