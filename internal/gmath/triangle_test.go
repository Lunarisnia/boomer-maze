package gmath

import "testing"

func TestIsInsideTriangleDocumentsCurrentWindingLimitation(t *testing.T) {
	a := Vector3[int32]{X: 0, Y: 0, Z: 0}
	b := Vector3[int32]{X: 10, Y: 0, Z: 0}
	c := Vector3[int32]{X: 0, Y: 10, Z: 0}
	p := Vector3[int32]{X: 2, Y: 2, Z: 0}

	if !IsInsideTriangle(a, b, c, p) {
		t.Fatal("expected point inside triangle with clockwise screen-space winding")
	}

	if IsInsideTriangle(a, c, b, p) {
		t.Fatal("current implementation is expected to reject reversed winding")
	}
}

func TestIsInsideTriangleDocumentsCurrentDegenerateTriangleBehavior(t *testing.T) {
	a := Vector3[int32]{X: 0, Y: 0, Z: 0}
	b := Vector3[int32]{X: 5, Y: 0, Z: 0}
	c := Vector3[int32]{X: 10, Y: 0, Z: 0}
	p := Vector3[int32]{X: 5, Y: 0, Z: 0}

	if !IsInsideTriangle(a, b, c, p) {
		t.Fatal("current implementation is expected to accept this degenerate case")
	}
}
