package gmath

type Number interface {
	int | float64 | int32
}

type Vector3[T Number] struct {
	X T
	Y T
	Z T
}

// a - b
func Sub[T Number](a, b Vector3[T]) Vector3[T] {
	return Vector3[T]{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func Cross[T Number](a, b Vector3[T]) Vector3[T] {
	return Vector3[T]{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func Cross2D[T Number](a, b Vector3[T]) T {
	return a.X*b.Y - a.Y*b.X
}
