package gmath

type Triangle[T Number] struct {
	A     Vector3[T]
	B     Vector3[T]
	C     Vector3[T]
	Color uint32
}

func NewTriangle[T Number](a, b, c Vector3[T], color uint32) Triangle[T] {
	return Triangle[T]{
		A:     a,
		B:     b,
		C:     c,
		Color: color,
	}
}

// IsInsideTriangle checks whether point p is inside triangle abc.
//
// Each triangle edge acts like a directed line that splits the plane into two
// halves. For edge AB, Cross2D(B-A, P-A) tells which side of that directed edge
// point P is on:
//
//   - positive: P is on one side of AB
//   - negative: P is on the other side of AB
//   - zero: P is exactly on the line AB
//
// A triangle is the overlap of the three inside half-planes made by edges AB,
// BC, and CA. If the triangle vertices are ordered consistently, an inside
// point stays on the same side of every edge. If one edge test has the opposite
// sign, the point is outside the triangle.
func IsInsideTriangle[T Number](a, b, c, p Vector3[T]) bool {
	ab := Sub(b, a)
	bc := Sub(c, b)
	ca := Sub(a, c)
	list := []Vector3[T]{
		ab, bc, ca,
	}
	rawPoints := []Vector3[T]{
		a, b, c,
	}

	for i := range len(list) {
		toPoint := Sub(p, rawPoints[i])
		if Cross2D(list[i], toPoint) < 0 {
			return false
		}
	}
	return true
}
