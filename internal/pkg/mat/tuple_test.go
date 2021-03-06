package mat

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTuple4_IsVector(t *testing.T) {
	v := NewVector(4.3, -4.2, 3.1)
	assert.True(t, v.IsVector())
	assert.False(t, v.IsPoint())

	assert.Equal(t, 4.3, v.Get(0))
	assert.Equal(t, -4.2, v.Get(1))
	assert.Equal(t, 3.1, v.Get(2))
}

func TestTuple4_IsPoint(t *testing.T) {
	p := NewPoint(4.3, -4.2, 3.1)
	assert.True(t, p.IsPoint())
	assert.False(t, p.IsVector())

	assert.Equal(t, 4.3, p.Get(0))
	assert.Equal(t, -4.2, p.Get(1))
	assert.Equal(t, 3.1, p.Get(2))
}

func TestTuple4Add(t *testing.T) {
	t1 := NewPoint(3, -2, 5)
	t2 := NewVector(-2, 3, 1)
	t3 := Add(t1, t2)
	assert.Equal(t, 1.0, t3.Get(0))
	assert.Equal(t, 1.0, t3.Get(1))
	assert.Equal(t, 6.0, t3.Get(2))
	assert.Equal(t, 1.0, t3.Get(3))
}

func TestTuple4Sub(t *testing.T) {
	t1 := NewPoint(3, 2, 1)
	t2 := NewPoint(5, 6, 7)
	t3 := Sub(t1, t2)
	assert.Equal(t, -2.0, t3.Get(0))
	assert.Equal(t, -4.0, t3.Get(1))
	assert.Equal(t, -6.0, t3.Get(2))
	assert.Equal(t, 0.0, t3.Get(3))
}

func TestSubVectorFromPoint(t *testing.T) {
	t1 := NewPoint(3, 2, 1)
	t2 := NewVector(5, 6, 7)

	t3 := Sub(t1, t2)
	assert.Equal(t, -2.0, t3.Get(0))
	assert.Equal(t, -4.0, t3.Get(1))
	assert.Equal(t, -6.0, t3.Get(2))
	assert.Equal(t, 1.0, t3.Get(3))
}

func TestSubVectorFromVector(t *testing.T) {
	t1 := NewVector(3, 2, 1)
	t2 := NewVector(5, 6, 7)

	t3 := Sub(t1, t2)
	assert.Equal(t, -2.0, t3.Get(0))
	assert.Equal(t, -4.0, t3.Get(1))
	assert.Equal(t, -6.0, t3.Get(2))
	assert.Equal(t, 0.0, t3.Get(3))
}

func TestSubtractVectorFromZeroVector(t *testing.T) {
	t1 := NewVector(0, 0, 0)
	t2 := NewVector(1, -2, 3)

	t3 := Sub(t1, t2)
	assert.Equal(t, -1.0, t3.Get(0))
	assert.Equal(t, 2.0, t3.Get(1))
	assert.Equal(t, -3.0, t3.Get(2))
	assert.Equal(t, 0.0, t3.Get(3))
}

func TestNegateTuple(t *testing.T) {
	t1 := NewTupleOf(1, -2, 3, -4)
	t3 := Negate(t1)
	assert.Equal(t, -1.0, t3.Get(0))
	assert.Equal(t, 2.0, t3.Get(1))
	assert.Equal(t, -3.0, t3.Get(2))
	assert.Equal(t, 4.0, t3.Get(3))
}

func TestMultiplyByScalar(t *testing.T) {
	t1 := NewTupleOf(1, -2, 3, -4)
	t3 := MultiplyByScalar(t1, 3.5)
	assert.Equal(t, 3.5, t3.Get(0))
	assert.Equal(t, -7.0, t3.Get(1))
	assert.Equal(t, 10.5, t3.Get(2))
	assert.Equal(t, -14.0, t3.Get(3))
}

func TestMultiplyByScalarFraction(t *testing.T) {
	t1 := NewTupleOf(1, -2, 3, -4)
	t3 := MultiplyByScalar(t1, 0.5)
	assert.Equal(t, 0.5, t3.Get(0))
	assert.Equal(t, -1.0, t3.Get(1))
	assert.Equal(t, 1.5, t3.Get(2))
	assert.Equal(t, -2.0, t3.Get(3))
}

func TestDivideByScalar(t *testing.T) {
	t1 := NewTupleOf(1, -2, 3, -4)
	t3 := DivideByScalar(t1, 2)
	assert.Equal(t, 0.5, t3.Get(0))
	assert.Equal(t, -1.0, t3.Get(1))
	assert.Equal(t, 1.5, t3.Get(2))
	assert.Equal(t, -2.0, t3.Get(3))
}

func TestMagnitude(t *testing.T) {
	tc := []struct {
		tpl Tuple4
		out float64
	}{
		{NewVector(1, 0, 0), 1.0},
		{NewVector(0, 1, 0), 1.0},
		{NewVector(0, 0, 1), 1.0},
		{NewVector(1, 2, 3), math.Sqrt(14)},
		{NewVector(-1, -2, -3), math.Sqrt(14)},
	}

	for _, test := range tc {
		assert.Equal(t, test.out, Magnitude(test.tpl))
	}
}

func TestNormalizeXOnly(t *testing.T) {
	t1 := NewVector(4, 0, 0)
	t3 := Normalize(t1)
	assert.Equal(t, 1.0, t3.Get(0))
	assert.Equal(t, 0.0, t3.Get(1))
	assert.Equal(t, 0.0, t3.Get(2))
}

func TestNormalizeXYZ(t *testing.T) {
	t1 := NewVector(1, 2, 3)
	t3 := Normalize(t1)
	assert.True(t, Eq(0.26726, t3.Get(0)))
	assert.True(t, Eq(0.53452, t3.Get(1)))
	assert.True(t, Eq(0.80178, t3.Get(2)))
}

func TestNormalizedMagnitudeIsOne(t *testing.T) {
	t1 := NewVector(1, 2, 3)
	t3 := Normalize(t1)
	assert.Equal(t, 1.0, Magnitude(t3))
}

func TestDot(t *testing.T) {
	t1 := NewVector(1, 2, 3)
	t2 := NewVector(2, 3, 4)
	dotProduct := Dot(t1, t2)
	assert.Equal(t, 20.0, dotProduct)
}

func TestCross(t *testing.T) {
	t1 := NewVector(1, 2, 3)
	t2 := NewVector(2, 3, 4)
	crossT1 := Cross(t1, t2)
	crossT2 := Cross(t2, t1)
	assert.True(t, TupleEquals(crossT1, NewVector(-1, 2, -1)))
	assert.True(t, TupleEquals(crossT2, NewVector(1, -2, 1)))
}
func TestCrossProduct(t *testing.T) {

	t1 := NewVector(1, 2, 3)
	t2 := NewVector(2, 3, 4)
	out := Tuple4{}
	CrossProduct(&t1, &t2, &out)
	assert.True(t, TupleEquals(out, NewVector(-1, 2, -1)))
}

func TestColorAdd(t *testing.T) {
	c1 := NewColor(0.9, 0.6, 0.75)
	c2 := NewColor(0.7, 0.1, 0.25)
	c3 := Add(c1, c2)
	assert.Equal(t, 1.6, c3.Get(0))
	assert.Equal(t, 0.7, c3.Get(1))
	assert.Equal(t, 1.0, c3.Get(2))
}
func TestColorSub(t *testing.T) {
	c1 := NewColor(0.9, 0.6, 0.75)
	c2 := NewColor(0.7, 0.1, 0.25)
	c3 := Sub(c1, c2)
	assert.InEpsilon(t, 0.2, c3.Get(0), Epsilon)
	assert.InEpsilon(t, 0.5, c3.Get(1), Epsilon)
	assert.InEpsilon(t, 0.5, c3.Get(2), Epsilon)
}
func TestColorMultiplyByScalar(t *testing.T) {
	c1 := NewColor(0.2, 0.3, 0.4)
	c3 := MultiplyByScalar(c1, 2)
	assert.Equal(t, 0.4, c3.Get(0))
	assert.Equal(t, 0.6, c3.Get(1))
	assert.Equal(t, 0.8, c3.Get(2))
}
func TestHadamard(t *testing.T) {
	c1 := NewColor(1, 0.2, 0.4)
	c2 := NewColor(0.9, 1, 0.1)
	c3 := Hadamard(c1, c2)
	assert.InEpsilon(t, 0.9, c3.Get(0), Epsilon)
	assert.InEpsilon(t, 0.2, c3.Get(1), Epsilon)
	assert.InEpsilon(t, 0.04, c3.Get(2), Epsilon)
}
func BenchmarkDot(b *testing.B) {
	var res = 0.0
	t1 := NewVector(1, 2, 3)
	t2 := NewVector(2, 3, 4)
	for i := 0; i < b.N; i++ {
		res = Dot(t1, t2)
	}
	fmt.Printf("%v\n", res)
}
func BenchmarkNormalizePtr(b *testing.B) {
	t1 := NewVector(1, 2, 3)
	out := Tuple4{}
	for i := 0; i < b.N; i++ {
		NormalizePtr(&t1, &out)
	}
	fmt.Printf("%v\n", out)
}
func BenchmarkNormalizePtr2(b *testing.B) {
	t1 := NewVector(1, 2, 3)
	out := Tuple4{}
	for i := 0; i < b.N; i++ {
		NormalizePtr2(t1, &out)
	}
	fmt.Printf("%v\n", out)
}
func BenchmarkNormalize(b *testing.B) {
	t1 := NewVector(1, 2, 3)
	out := Tuple4{}
	for i := 0; i < b.N; i++ {
		out = Normalize(t1)
	}
	fmt.Printf("%v\n", out)
}
func BenchmarkCross(b *testing.B) {
	t1 := NewVector(1.3243, 2.35456, 3.65464)
	t2 := NewVector(2.6563, 3.75672, 4.54654)
	out := Tuple4{}
	for i := 0; i < b.N; i++ {
		out = Cross(t1, t2)
	}
	fmt.Printf("%v\n", out)
}
func BenchmarkCross2(b *testing.B) {
	t1 := NewVector(1.3243, 2.35456, 3.65464)
	t2 := NewVector(2.6563, 3.75672, 4.54654)
	out := Tuple4{}
	for i := 0; i < b.N; i++ {
		Cross2(&t1, &t2, &out)
	}
	fmt.Printf("%v\n", out)
}
func BenchmarkCross2Parallell(b *testing.B) {
	t1 := NewVector(1.3243, 2.35456, 3.65464)
	t2 := NewVector(2.6563, 3.75672, 4.54654)

	b.RunParallel(func(pb *testing.PB) {
		out := Tuple4{}
		for pb.Next() {
			Cross2(&t1, &t2, &out)
		}
		fmt.Printf("%v\n", out)
	})
}
func BenchmarkCrossProduct(b *testing.B) {
	t1 := NewVector(1.3243, 2.35456, 3.65464)
	t2 := NewVector(2.6563, 3.75672, 4.54654)
	out := Tuple4{}
	for i := 0; i < b.N; i++ {
		CrossProduct(&t1, &t2, &out)
	}
	fmt.Printf("%v\n", out)
}

func BenchmarkCrossProductParallell(b *testing.B) {
	t1 := NewVector(1.3243, 2.35456, 3.65464)
	t2 := NewVector(2.6563, 3.75672, 4.54654)

	b.RunParallel(func(pb *testing.PB) {
		out := Tuple4{}
		for pb.Next() {
			CrossProduct(&t1, &t2, &out)
		}
		fmt.Printf("%v\n", out)
	})
}

func BenchmarkAdd(b *testing.B) {
	t1 := NewPoint(3, -2, 5)
	t2 := NewVector(-2, 3, 1)
	var t3 Tuple4
	for i := 0; i < b.N; i++ {
		t3 = Add(t1, t2)
	}
	fmt.Printf("%v\n", t3)
}
func BenchmarkAddPtr(b *testing.B) {
	t1 := NewPoint(3, -2, 5)
	t2 := NewVector(-2, 3, 1)
	var t3 Tuple4
	for i := 0; i < b.N; i++ {
		AddPtr(t1, t2, &t3)
	}
	fmt.Printf("%v\n", t3)
}
func BenchmarkAddPtr2(b *testing.B) {
	t1 := NewPoint(3, -2, 5)
	t2 := NewVector(-2, 3, 1)
	var t3 Tuple4
	for i := 0; i < b.N; i++ {
		AddPtr2(&t1, &t2, &t3)
	}
	fmt.Printf("%v\n", t3)
}
