package gmath

import "math"

func LocalToScreen(local Vector3[float64], objPos Vector3[float64], screenW, screenH int, fovDeg float64) (int, int, bool) {
	worldX := local.X + objPos.X
	worldY := local.Y + objPos.Y
	worldZ := local.Z + objPos.Z

	if worldZ <= 0.01 {
		return 0, 0, false
	}

	aspect := float64(screenW) / float64(screenH)
	fovRad := fovDeg * math.Pi / 180.0
	f := 1.0 / math.Tan(fovRad/2.0)

	ndcX := (worldX / worldZ) * f / aspect
	ndcY := (worldY / worldZ) * f

	screenX := int((ndcX*0.5 + 0.5) * float64(screenW))
	screenY := int((1.0 - (ndcY*0.5 + 0.5)) * float64(screenH))

	return screenX, screenY, true
}
