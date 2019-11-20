package mat

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestRotateX(t *testing.T) {
	point := NewPoint(0, 1, 0)
	halfQuarterRotation := RotateX(math.Pi / 4)
	fullQuarterRotation := RotateX(math.Pi / 2)

	p2 := MultiplyByTuple(halfQuarterRotation, point)
	assert.Equal(t, 0.0, p2.Get(0))
	assert.InEpsilon(t, math.Sqrt(2.0)/2.0, p2.Get(1), Epsilon, fmt.Sprintf("expected %v, got %v", math.Sqrt(2.0)/2.0, p2.Get(1)))
	assert.InEpsilon(t, math.Sqrt(2.0)/2.0, p2.Get(2), Epsilon, fmt.Sprintf("expected %v, got %v", math.Sqrt(2.0)/2.0, p2.Get(2)))

	p3 := MultiplyByTuple(fullQuarterRotation, point)
	assert.Equal(t, 0.0, p3.Get(0))
	assert.True(t, Eq(0.0, p3.Get(1)))
	assert.Equal(t, 1.0, p3.Get(2))
}

func TestRotateXInverse(t *testing.T) {
	point := NewPoint(0, 1, 0)
	halfQuarterRotation := RotateX(math.Pi / 4)
	inverseHQ := Inverse(halfQuarterRotation)
	p2 := MultiplyByTuple(inverseHQ, point)
	assert.Equal(t, 0.0, p2.Get(0))
	assert.InEpsilon(t, math.Sqrt(2.0)/2.0, p2.Get(1), Epsilon)
	assert.InEpsilon(t, -math.Sqrt(2.0)/2.0, p2.Get(2), Epsilon)
}

func TestRotateY(t *testing.T) {
	point := NewPoint(0, 0, 1)
	halfQuarterRotation := RotateY(math.Pi / 4)
	fullQuarterRotation := RotateY(math.Pi / 2)

	p2 := MultiplyByTuple(halfQuarterRotation, point)
	assert.InEpsilon(t, math.Sqrt(2.0)/2.0, p2.Get(0), Epsilon, fmt.Sprintf("expected %v, got %v", math.Sqrt(2.0)/2.0, p2.Get(0)))
	assert.Equal(t, 0.0, p2.Get(1))
	assert.InEpsilon(t, math.Sqrt(2.0)/2.0, p2.Get(2), Epsilon, fmt.Sprintf("expected %v, got %v", math.Sqrt(2.0)/2.0, p2.Get(2)))

	p3 := MultiplyByTuple(fullQuarterRotation, point)
	assert.Equal(t, 1.0, p3.Get(0))
	assert.True(t, Eq(0.0, p3.Get(1)))
	assert.True(t, Eq(0.0, p3.Get(2)))
}

func TestRotateZ(t *testing.T) {
	point := NewPoint(0, 1, 0)
	halfQuarterRotation := RotateZ(math.Pi / 4)
	fullQuarterRotation := RotateZ(math.Pi / 2)

	p2 := MultiplyByTuple(halfQuarterRotation, point)
	assert.InEpsilon(t, -math.Sqrt(2.0)/2.0, p2.Get(0), Epsilon, fmt.Sprintf("expected %v, got %v", math.Sqrt(2.0)/2.0, p2.Get(0)))
	assert.InEpsilon(t, math.Sqrt(2.0)/2.0, p2.Get(1), Epsilon, fmt.Sprintf("expected %v, got %v", math.Sqrt(2.0)/2.0, p2.Get(1)))
	assert.Equal(t, 0.0, p2.Get(2))

	p3 := MultiplyByTuple(fullQuarterRotation, point)
	assert.Equal(t, -1.0, p3.Get(0))
	assert.True(t, Eq(0.0, p3.Get(1)))
	assert.True(t, Eq(0.0, p3.Get(2)))
}
