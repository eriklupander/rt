package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIntersection(t *testing.T) {
	sp := NewSphere()
	intersection := NewIntersection(3.5, sp)
	assert.Equal(t, 3.5, intersection.T)
	assert.Equal(t, sp.Id, intersection.S.Id)
}
