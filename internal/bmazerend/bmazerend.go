package bmazerend

import (
	"math"

	"github.com/lunarisnia/boomer-maze/internal/bmazerend/window"
	"github.com/lunarisnia/boomer-maze/internal/gmath"
	"github.com/lunarisnia/boomer-maze/internal/valueutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

// Bressenham line drawing algorithm, I don't understand most of this
func line(ax int32, ay int32, bx int32, by int32, fb *window.Framebuffer, color uint32) {
	steep := math.Abs(float64(ax-bx)) < math.Abs(float64(ay-by))
	if steep {
		ax, ay = valueutil.Swap(ax, ay)
		bx, by = valueutil.Swap(bx, by)
	}
	if ax > bx {
		ax, bx = valueutil.Swap(ax, bx)
		ay, by = valueutil.Swap(ay, by)
	}
	y := float64(ay)
	var iError int32 = 0
	for x := ax; x <= bx; x++ {
		if steep {
			fb.SetPixel(int32(y), x, color)
		} else {
			fb.SetPixel(x, int32(y), color)
		}
		iError += 2 * int32(math.Abs(float64(by-ay)))
		if iError > bx-ax {
			if by > ay {
				y += 1
			} else {
				y += -1
			}
			iError -= int32(2 * (float64(bx) - float64(ax)))
		}
	}
}

type BoundingBox struct {
	MinX int32
	MinY int32

	MaxX int32
	MaxY int32
}

func findBoundingBox(t gmath.Triangle[int32]) BoundingBox {
	return BoundingBox{
		MinX: min(min(t.A.X, t.B.X), t.C.X),
		MinY: min(min(t.A.Y, t.B.Y), t.C.Y),

		MaxX: max(max(t.A.X, t.B.X), t.C.X),
		MaxY: max(max(t.A.Y, t.B.Y), t.C.Y),
	}
}

func signedTriangleArea(triangle gmath.Triangle[int32]) float64 {
	a := (float64(triangle.B.Y) - float64(triangle.A.Y)) * (float64(triangle.B.X) + float64(triangle.A.X))
	b := (float64(triangle.C.Y) - float64(triangle.B.Y)) * (float64(triangle.C.X) + float64(triangle.B.X))
	c := (float64(triangle.A.Y) - float64(triangle.C.Y)) * (float64(triangle.A.X) + float64(triangle.C.X))
	return 0.5 * (a + b + c)
}

func rasterize(fb *window.Framebuffer, triangles []gmath.Triangle[int32]) {
	for _, t := range triangles {
		boundingBox := findBoundingBox(t)
		totalArea := signedTriangleArea(t)
		if totalArea < 1 {
			continue
		}
		// Iterate over all pixels on the screen
		for y := boundingBox.MinY; y < boundingBox.MaxY; y++ {
			for x := boundingBox.MinX; x < boundingBox.MaxX; x++ {
				alpha := signedTriangleArea(gmath.NewTriangle(gmath.Vector3[int32]{
					X: x,
					Y: y,
					Z: 0,
				}, t.B, t.C, window.ColorBlue)) / totalArea
				beta := signedTriangleArea(gmath.NewTriangle(gmath.Vector3[int32]{
					X: x,
					Y: y,
					Z: 0,
				}, t.C, t.A, window.ColorBlue)) / totalArea
				gamma := signedTriangleArea(gmath.NewTriangle(gmath.Vector3[int32]{
					X: x,
					Y: y,
					Z: 0,
				}, t.A, t.B, window.ColorBlue)) / totalArea
				// Check if point is inside the triangle
				isOutside := alpha < 0 || beta < 0 || gamma < 0
				if isOutside {
					continue
				}

				red := uint8(alpha * float64(t.A.Z))
				green := uint8(beta * float64(t.B.Z))
				blue := uint8(gamma * float64(t.C.Z))
				fb.SetPixel(x, y, window.ARGB(0xFF, red, green, blue))
			}
		}
	}
}

func triangle(fb *window.Framebuffer, t gmath.Triangle[int32]) {
	rasterize(fb, []gmath.Triangle[int32]{t})
}

func draw(ctx *window.WindowContext) {
	ctx.Framebuffer.Clear(0xFF000000)
	t := gmath.Triangle[int32]{
		A: gmath.Vector3[int32]{
			X: 100,
			Y: 400,
			Z: 255,
		},
		B: gmath.Vector3[int32]{
			X: 300,
			Y: 50,
			Z: 255,
		},
		C: gmath.Vector3[int32]{
			X: 400,
			Y: 70,
			Z: 255,
		},
		Color: window.ColorGreen,
	}
	triangle(ctx.Framebuffer, t)

}

func Run() {
	win := window.New("Boomer Maze", screenWidth, screenHeight)
	win.SetDrawFunction(draw)

	win.Run()
}
