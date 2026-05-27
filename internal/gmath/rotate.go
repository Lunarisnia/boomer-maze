package gmath

import "math"

func RotateY(v Vector3[float64], angleRad float64) Vector3[float64] {
	c := math.Cos(angleRad)
	s := math.Sin(angleRad)

	return Vector3[float64]{
		X: v.X*c + v.Z*s,
		Y: v.Y,
		Z: -v.X*s + v.Z*c,
	}
}
