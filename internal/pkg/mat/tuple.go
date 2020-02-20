package mat

import "math"

type Tuple4 [4]float64

func (t Tuple4) ResetVector() Tuple4 {
	t[0] = 0.0
	t[1] = 0.0
	t[2] = 0.0
	t[3] = 1.0
	return t
}
func (t Tuple4) ResetPoint() Tuple4 {
	t[0] = 0.0
	t[1] = 0.0
	t[2] = 0.0
	t[3] = 0.0
	return t
}

func NewVector(x, y, z float64) Tuple4 {
	return Tuple4{x, y, z, 0.0}
}
func NewPoint(x, y, z float64) Tuple4 {
	return Tuple4{x, y, z, 1.0}
}
func NewColor(r, g, b float64) Tuple4 {
	return Tuple4{r, g, b, 1.0}
}

func NewTuple() Tuple4 {
	return [4]float64{0, 0, 0, 0}
}
func NewTupleOf(x, y, z, w float64) Tuple4 {
	return [4]float64{x, y, z, w}
}

//func NewTuple4(elems []float64) Tuple4 {
//	return Tuple4{elems[0], elems[1], elems[2], elems[3]}
//}

func (t Tuple4) Get(row int) float64 {
	return t[row]
}

func (t Tuple4) IsVector() bool {
	return t[3] == 0.0
}
func (t Tuple4) IsPoint() bool {
	return t[3] == 1.0
}

func Add(t1, t2 Tuple4) Tuple4 {
	t3 := [4]float64{}
	for i := 0; i < 4; i++ {
		t3[i] = t1.Get(i) + t2.Get(i)
	}
	return t3
}

func AddPtr(t1, t2 Tuple4, t3 *Tuple4) {
	for i := 0; i < 4; i++ {
		t3[i] = t1.Get(i) + t2.Get(i)
	}
}

func Add3(t1, t2, t3 Tuple4, out *Tuple4) {
	for i := 0; i < 4; i++ {
		out[i] = t1.Get(i) + t2.Get(i) + t3.Get(i)
	}
}

func (t3 Tuple4) Add(t2 Tuple4) Tuple4 {
	for i := 0; i < 4; i++ {
		t3[i] = t3.Get(i) + t2.Get(i)
	}
	return t3
}

func Sub(t1, t2 Tuple4) Tuple4 {
	t3 := [4]float64{}
	for i := 0; i < 4; i++ {
		t3[i] = t1.Get(i) - t2.Get(i)
	}
	return t3
}

func SubPtr(t1, t2 Tuple4, out *Tuple4) {
	for i := 0; i < 4; i++ {
		out[i] = t1.Get(i) - t2.Get(i)
	}
}

func Negate(t1 Tuple4) Tuple4 {
	t3 := [4]float64{}
	for i := 0; i < 4; i++ {
		t3[i] = 0 - t1.Get(i)
	}
	return t3
}

func NegatePtr(t1 Tuple4, out *Tuple4) {
	for i := 0; i < 4; i++ {
		out[i] = 0 - t1.Get(i)
	}
}

func MultiplyByScalar(t1 Tuple4, scalar float64) Tuple4 {
	t3 := [4]float64{}
	for i := 0; i < 4; i++ {
		t3[i] = t1.Get(i) * scalar
	}
	return t3
}

func MultiplyByScalarPtr(t1 Tuple4, scalar float64, out *Tuple4) {
	for i := 0; i < 4; i++ {
		out[i] = t1.Get(i) * scalar
	}
}

func (t3 Tuple4) Multiply(scalar float64) Tuple4 {
	for i := 0; i < 4; i++ {
		t3[i] = t3[i] * scalar
	}
	return t3
}

func DivideByScalar(t1 Tuple4, scalar float64) Tuple4 {
	t3 := [4]float64{}
	for i := 0; i < 4; i++ {
		t3[i] = t1.Get(i) / scalar
	}
	return t3
}

// Magnitude measures the length of the passed vector. It's basically pythagoras sqrt(x2 + y2 + z2 + w2)
func Magnitude(t1 Tuple4) float64 {
	return math.Sqrt(t1[0]*t1[0] +
		t1[1]*t1[1] +
		t1[2]*t1[2])

}

func MagnitudePtr(t1 *Tuple4) float64 {
	return math.Sqrt(t1[0]*t1[0] +
		t1[1]*t1[1] +
		t1[2]*t1[2])

}

// Normalize measures the length (magnitude) of the passed Vector. Each component in t1 is then divided my the magnitude
// in order to Normalize it to unit (1) size.
func Normalize(t1 Tuple4) Tuple4 {
	t3 := [4]float64{}
	magnitude := Magnitude(t1)
	for i := 0; i < 4; i++ {
		t3[i] = t1.Get(i) / magnitude
	}
	return t3
}

func NormalizePtr(t1 Tuple4, out *Tuple4) {
	magnitude := Magnitude(t1)
	var x, y, z, w float64

	x = t1[0] / magnitude
	y = t1[1] / magnitude
	z = t1[2] / magnitude
	w = t1[3] / magnitude

	out[0] = x
	out[1] = y
	out[2] = z
	out[3] = w
}

// Dot product is the sum of the products of the corresponding entries of the two sequences of numbers
// a product is simply put the result of a multiplication. The dot product of two tuples is simply
// t1.x * t2.x + t1.y * t2.y + t1.z * t2.z + t1.w * t2.w
func Dot(t1 Tuple4, t2 Tuple4) float64 {
	sum := 0.0
	for i := 0; i < 4; i++ {
		sum += t1[i] * t2[i]
	}
	return sum
}

func Cross(t1 Tuple4, t2 Tuple4) Tuple4 {
	t3 := [4]float64{}

	t3[0] = t1[1]*t2[2] - t1[2]*t2[1]
	t3[1] = t1[2]*t2[0] - t1[0]*t2[2]
	t3[2] = t1[0]*t2[1] - t1[1]*t2[0]
	t3[3] = 0
	return t3
}

func Hadamard(t1 Tuple4, t2 Tuple4) Tuple4 {
	t3 := [4]float64{}
	t3[0] = t1[0] * t2[0]
	t3[1] = t1[1] * t2[1]
	t3[2] = t1[2] * t2[2]
	t3[3] = 1.0
	return t3
}

func HadamardPtr(t1 Tuple4, t2 Tuple4, out *Tuple4) {
	out[0] = t1[0] * t2[0]
	out[1] = t1[1] * t2[1]
	out[2] = t1[2] * t2[2]
	out[3] = 1.0
}

func TupleEquals(t1, t2 Tuple4) bool {
	return Eq(t1[0], t2[0]) &&
		Eq(t1[1], t2[1]) &&
		Eq(t1[2], t2[2]) &&
		Eq(t1[3], t2[3])
}

func TupleXYZEq(t1, t2 Tuple4) bool {
	return Eq(t1[0], t2[0]) &&
		Eq(t1[1], t2[1]) &&
		Eq(t1[2], t2[2])
}
