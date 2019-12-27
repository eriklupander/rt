package mat

func NewComputation() Computation {
	containers := make([]Shape, 8)
	containers = containers[:0]
	return Computation{
		T:          0,
		Object:     nil,
		Point:      NewPoint(0, 0, 0),
		EyeVec:     NewVector(0, 0, 0),
		NormalVec:  NewVector(0, 0, 0),
		Inside:     false,
		OverPoint:  NewPoint(0, 0, 0),
		UnderPoint: NewPoint(0, 0, 0),
		ReflectVec: NewVector(0, 0, 0),
		N1:         0,
		N2:         0,

		localPoint:   NewPoint(0, 0, 0),
		containers:   containers,
		cachedOffset: NewVector(0, 0, 0),
	}
}

func PrepareComputationForIntersectionPtr(i Intersection, r Ray, comps *Computation, xs ...Intersection) {
	comps.T = i.T
	comps.Object = i.S
	PositionPtr(r, i.T, &comps.Point)
	NegatePtr(r.Direction, &comps.EyeVec)
	comps.NormalVec = NormalAt(i.S, comps.Point, &i) //  fix
	//comps.NormalVec = NormalAtPtr(i.S, comps.Point, &i, &comps.localPoint) //  fix
	ReflectPtr(r.Direction, comps.NormalVec, &comps.ReflectVec)

	comps.Inside = false
	if Dot(comps.EyeVec, comps.NormalVec) < 0 {
		comps.Inside = true
		NegatePtr(comps.NormalVec, &comps.NormalVec) // fix
	}
	MultiplyByScalarPtr(comps.NormalVec, Epsilon, &comps.cachedOffset)
	AddPtr(comps.Point, comps.cachedOffset, &comps.OverPoint)
	SubPtr(comps.Point, comps.cachedOffset, &comps.UnderPoint)

	comps.N1 = 1.0
	comps.N2 = 1.0

	comps.containers = comps.containers[:0] // make([]Shape, 0)
	for idx := range xs {
		if xs[idx].S.ID() == i.S.ID() && i.T == xs[idx].T {
			if len(comps.containers) == 0 {
				comps.N1 = 1.0
			} else {
				comps.N1 = comps.containers[len(comps.containers)-1].GetMaterial().RefractiveIndex
			}
		}

		index := indexOf(xs[idx].S, comps.containers)
		if index > -1 {
			copy(comps.containers[index:], comps.containers[index+1:])    // Shift a[i+1:] left one indexs[idx].
			comps.containers[len(comps.containers)-1] = nil               // Erase last element (write zero value).
			comps.containers = comps.containers[:len(comps.containers)-1] // Truncate slice.
		} else {
			comps.containers = append(comps.containers, xs[idx].S)
		}

		if xs[idx].S.ID() == i.S.ID() && xs[idx].T == i.T {
			if len(comps.containers) == 0 {
				comps.N2 = 1.0
			} else {
				comps.N2 = comps.containers[len(comps.containers)-1].GetMaterial().RefractiveIndex
			}
			break
		}
	}
}

// TODO refactor so tests uses the other one instead.
func PrepareComputationForIntersection(i Intersection, r Ray, xs ...Intersection) Computation {
	pos := Position(r, i.T)
	eyeVec := Negate(r.Direction)
	normalVec := NormalAt(i.S, pos, &i)
	reflectVec := Reflect(r.Direction, normalVec)
	inside := false
	if Dot(eyeVec, normalVec) < 0 {
		inside = true
		normalVec = Negate(normalVec)
	}
	overPoint := Add(pos, MultiplyByScalar(normalVec, Epsilon))
	underPoint := Sub(pos, MultiplyByScalar(normalVec, Epsilon))

	n1 := 1.0
	n2 := 1.0

	var containers = make([]Shape, 0)
	for idx := range xs {
		if xs[idx].S.ID() == i.S.ID() && i.T == xs[idx].T {
			if len(containers) == 0 {
				n1 = 1.0
			} else {
				n1 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
		}

		index := indexOf(xs[idx].S, containers)
		if index > -1 {
			copy(containers[index:], containers[index+1:]) // Shift a[i+1:] left one indexs[idx].
			containers[len(containers)-1] = nil            // Erase last element (write zero value).
			containers = containers[:len(containers)-1]    // Truncate slice.
		} else {
			containers = append(containers, xs[idx].S)
		}

		if xs[idx].S.ID() == i.S.ID() && xs[idx].T == i.T {
			if len(containers) == 0 {
				n2 = 1.0
			} else {
				n2 = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
			break
		}
	}

	return Computation{
		T:          i.T,
		Object:     i.S,
		Point:      pos,
		EyeVec:     eyeVec,
		NormalVec:  normalVec,
		ReflectVec: reflectVec,
		Inside:     inside,
		OverPoint:  overPoint,
		UnderPoint: underPoint,
		N1:         n1,
		N2:         n2,
	}
}

func indexOf(s Shape, list []Shape) int {
	for idx := range list {
		if list[idx].ID() == s.ID() {
			return idx
		}
	}
	return -1
}

type Computation struct {
	T          float64
	Object     Shape
	Point      Tuple4
	EyeVec     Tuple4
	NormalVec  Tuple4
	Inside     bool
	OverPoint  Tuple4
	UnderPoint Tuple4
	ReflectVec Tuple4
	N1         float64
	N2         float64

	// cached stuff
	localPoint   Tuple4
	containers   []Shape
	cachedOffset Tuple4
}

func NewLightData() LightData {
	return LightData{
		Ambient:        NewColor(0, 0, 0),
		Diffuse:        NewColor(0, 0, 0),
		Specular:       NewColor(0, 0, 0),
		EffectiveColor: NewColor(0, 0, 0),
		LightVec:       NewVector(0, 0, 0),
		ReflectVec:     NewVector(0, 0, 0),
	}
}

// LightData is used for pre-allocated memory for lighting computations.
type LightData struct {
	Ambient        Tuple4
	Diffuse        Tuple4
	Specular       Tuple4
	EffectiveColor Tuple4
	LightVec       Tuple4
	ReflectVec     Tuple4
}
