package mat

type Shape interface {
	ID() int64
	GetTransform() Mat4x4
	SetTransform(transform Mat4x4)
	GetMaterial() Material
	SetMaterial(material Material)
	IntersectLocal(ray Ray) []Intersection
	NormalAtLocal(point Tuple4) Tuple4
	GetLocalRay() Ray
}
