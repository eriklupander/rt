package mat

import "math"

type Tuple4 struct {
	Elems []float64
}

func NewVector(x, y, z float64) Tuple4 {
	return Tuple4{[]float64{x, y, z, 0.0}}
}
func NewPoint(x, y, z float64) Tuple4 {
	return Tuple4{[]float64{x, y, z, 1.0}}
}
func NewColor(r, g, b float64) Tuple4 {
	return Tuple4{[]float64{r, g, b, 1.0}}
}

func NewTuple4(elems []float64) Tuple4 {
	return Tuple4{Elems: elems}
}

func (t Tuple4) Get(row int) float64 {
	return t.Elems[row]
}

func (t Tuple4) IsVector() bool {
	return t.Elems[3] == 0.0
}
func (t Tuple4) IsPoint() bool {
	return t.Elems[3] == 1.0
}

func Add(t1, t2 Tuple4) Tuple4 {
	t3 := NewTuple4(make([]float64, 4))
	for i := 0; i < 4; i++ {
		t3.Elems[i] = t1.Get(i) + t2.Get(i)
	}
	return t3
}

func AddPtr(t1, t2 Tuple4, out *Tuple4) {
	for i := 0; i < 4; i++ {
		out.Elems[i] = t1.Get(i) + t2.Get(i)
	}
}

func Sub(t1, t2 Tuple4) Tuple4 {
	t3 := NewTuple4(make([]float64, 4))
	for i := 0; i < 4; i++ {
		t3.Elems[i] = t1.Get(i) - t2.Get(i)
	}
	return t3
}

func SubPtr(t1, t2 Tuple4, out *Tuple4) {
	for i := 0; i < 4; i++ {
		out.Elems[i] = t1.Get(i) - t2.Get(i)
	}
}

func Negate(t1 Tuple4) Tuple4 {
	t3 := NewTuple4(make([]float64, 4))
	for i := 0; i < 4; i++ {
		t3.Elems[i] = 0 - t1.Get(i)
	}
	return t3
}

func NegatePtr(t1 Tuple4, out *Tuple4) {
	for i := 0; i < 4; i++ {
		out.Elems[i] = 0 - t1.Get(i)
	}
}

func MultiplyByScalar(t1 Tuple4, scalar float64) Tuple4 {
	t3 := Tuple4{Elems: make([]float64, 4)}
	for i := 0; i < 4; i++ {
		t3.Elems[i] = t1.Get(i) * scalar
	}
	return t3
}

func DivideByScalar(t1 Tuple4, scalar float64) Tuple4 {
	t3 := Tuple4{Elems: make([]float64, 4)}
	for i := 0; i < 4; i++ {
		t3.Elems[i] = t1.Get(i) / scalar
	}
	return t3
}

// Magnitude measures the length of the passed vector. It's basically pythagoras sqrt(x2 + y2 + z2 + w2)
func Magnitude(t1 Tuple4) float64 {
	return math.Sqrt(t1.Elems[0]*t1.Elems[0] +
		t1.Elems[1]*t1.Elems[1] +
		t1.Elems[2]*t1.Elems[2])

}

// Normalize measures the length (magnitude) of the passed Vector. Each component in t1 is then divided my the magnitude
// in order to Normalize it to unit (1) size.
func Normalize(t1 Tuple4) Tuple4 {
	t3 := Tuple4{Elems: make([]float64, 4)}
	magnitude := Magnitude(t1)
	for i := 0; i < 4; i++ {
		t3.Elems[i] = t1.Get(i) / magnitude
	}
	return t3
}

func NormalizePtr(t1 Tuple4, out *Tuple4) {
	magnitude := Magnitude(t1)
	for i := 0; i < 4; i++ {
		out.Elems[i] = t1.Get(i) / magnitude
	}
}

// Dot product is the sum of the products of the corresponding entries of the two sequences of numbers
// a product is simply put the result of a multiplication. The dot product of two tuples is simply
// t1.x * t2.x + t1.y * t2.y + t1.z * t2.z + t1.w * t2.w
func Dot(t1 Tuple4, t2 Tuple4) float64 {
	sum := 0.0
	for i := 0; i < 4; i++ {
		sum += t1.Get(i) * t2.Get(i)
	}
	return sum
}

func Cross(t1 Tuple4, t2 Tuple4) Tuple4 {
	t3 := Tuple4{Elems: make([]float64, 4)}

	t3.Elems[0] = t1.Get(1)*t2.Get(2) - t1.Get(2)*t2.Get(1)
	t3.Elems[1] = t1.Get(2)*t2.Get(0) - t1.Get(0)*t2.Get(2)
	t3.Elems[2] = t1.Get(0)*t2.Get(1) - t1.Get(1)*t2.Get(0)
	t3.Elems[3] = 0
	return t3
}

func Hadamard(t1 Tuple4, t2 Tuple4) Tuple4 {
	t3 := Tuple4{Elems: make([]float64, 4)}
	t3.Elems[0] = t1.Get(0) * t2.Get(0)
	t3.Elems[1] = t1.Get(1) * t2.Get(1)
	t3.Elems[2] = t1.Get(2) * t2.Get(2)
	t3.Elems[3] = 1.0
	return t3
}

func TupleEquals(t1, t2 Tuple4) bool {
	return Eq(t1.Get(0), t2.Get(0)) &&
		Eq(t1.Get(1), t2.Get(1)) &&
		Eq(t1.Get(2), t2.Get(2)) &&
		Eq(t1.Get(3), t2.Get(3))
}
