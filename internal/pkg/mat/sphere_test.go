package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSphereHasIdentity(t *testing.T) {
	sphere := NewSphere()
	assert.True(t, Equals(*sphere.Transform, *IdentityMatrix))
}
