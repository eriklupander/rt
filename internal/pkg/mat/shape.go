package mat

type Shape interface {
	ID() int64
	GetTransform() Mat4x4
	GetInverse() Mat4x4
	GetInverseTranspose() Mat4x4
	SetTransform(transform Mat4x4)
	GetMaterial() Material
	SetMaterial(material Material)
	IntersectLocal(ray Ray) []Intersection
	NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4
	GetLocalRay() Ray
	GetParent() Shape
	SetParent(shape Shape)
	CastsShadow() bool
	Name() string

	// Init()
}

func WorldToObject(shape Shape, point Tuple4) Tuple4 {
	if shape.GetParent() != nil {
		point = WorldToObject(shape.GetParent(), point)
	}
	return MultiplyByTuple(shape.GetInverse(), point)
}

func WorldToObjectPtr(shape Shape, point Tuple4, out *Tuple4) {
	if shape.GetParent() != nil {
		WorldToObjectPtr(shape.GetParent(), point, &point)
	}
	i := shape.GetInverse()
	MultiplyByTuplePtr(&i, &point, out)
}

func NormalToWorld(shape Shape, normal Tuple4) Tuple4 {
	normal = MultiplyByTuple(shape.GetInverseTranspose(), normal)
	normal[3] = 0.0 // set w to 0
	normal = Normalize(normal)

	if shape.GetParent() != nil {
		normal = NormalToWorld(shape.GetParent(), normal)
	}
	return normal
}

func NormalToWorldPtr(shape Shape, normal *Tuple4) {
	it := shape.GetInverseTranspose()
	MultiplyByTuplePtr(&it, normal, normal)
	normal[3] = 0.0 // set w to 0
	NormalizePtr(normal, normal)

	if shape.GetParent() != nil {
		NormalToWorldPtr(shape.GetParent(), normal)
	}
}
