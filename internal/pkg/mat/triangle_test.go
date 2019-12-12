package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTriangle(t *testing.T) {
	p1 := NewPoint(0, 1, 0)
	p2 := NewPoint(-1, 0, 0)
	p3 := NewPoint(1, 0, 0)
	tri := NewTriangle(p1, p2, p3)

	assert.Equal(t, p1, tri.P1)
	assert.Equal(t, p2, tri.P2)
	assert.Equal(t, p3, tri.P3)

	assert.Equal(t, NewVector(-1, -1, 0), tri.E1)
	assert.Equal(t, NewVector(1, -1, 0), tri.E2)
	assert.Equal(t, NewVector(0, 0, -1), tri.N)
}

func TestFindNormalOnTriangle(t *testing.T) {
	tri := NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0))
	n1 := tri.NormalAtLocal(NewPoint(0, 0.5, 0), nil)
	n2 := tri.NormalAtLocal(NewPoint(-0.5, 0.75, 0), nil)
	n3 := tri.NormalAtLocal(NewPoint(0.5, 0.25, 0), nil)
	assert.Equal(t, tri.N, n1)
	assert.Equal(t, tri.N, n2)
	assert.Equal(t, tri.N, n3)
}

func TestIntersectTriangleParallellMisses(t *testing.T) {
	tri := NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0))
	ray := NewRay(NewPoint(0, -1, -2), NewVector(0, 1, 0))
	xs := tri.IntersectLocal(ray)
	assert.Len(t, xs, 0)
}

func TestIntersectTriangleMissesP1P3Edge(t *testing.T) {
	tri := NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0))
	ray := NewRay(NewPoint(1, 1, -2), NewVector(0, 0, 1))
	xs := tri.IntersectLocal(ray)
	assert.Len(t, xs, 0)
}

func TestIntersectTriangleMissesP1P2Edge(t *testing.T) {
	tri := NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0))
	ray := NewRay(NewPoint(-1, 1, -2), NewVector(0, 0, 1))
	xs := tri.IntersectLocal(ray)
	assert.Len(t, xs, 0)
}

func TestIntersectTriangleMissesP2P3Edge(t *testing.T) {
	tri := NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0))
	ray := NewRay(NewPoint(0, -1, -2), NewVector(0, 0, 1))
	xs := tri.IntersectLocal(ray)
	assert.Len(t, xs, 0)
}

func TestIntersectTriangleHits(t *testing.T) {
	tri := NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0))
	ray := NewRay(NewPoint(0, 0.5, -2), NewVector(0, 0, 1))
	xs := tri.IntersectLocal(ray)
	assert.Len(t, xs, 1)
	assert.Equal(t, 2.0, xs[0].T)
}
