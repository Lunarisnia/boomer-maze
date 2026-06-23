package bmazerend

import (
	"math"
	"unsafe"

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
	minX := int32(0)
	minY := int32(0)

	maxX := int32(0)
	maxY := int32(0)
	if t.A.X <= t.B.X && t.A.X <= t.C.X {
		minX = t.A.X
	} else if t.B.X <= t.A.X && t.B.X <= t.C.X {
		minX = t.B.X
	} else if t.C.X <= t.A.X && t.C.X <= t.B.X {
		minX = t.C.X
	}

	if t.A.X >= t.B.X && t.A.X >= t.C.X {
		maxX = t.A.X
	} else if t.B.X >= t.A.X && t.B.X >= t.C.X {
		maxX = t.B.X
	} else if t.C.X >= t.A.X && t.C.X >= t.B.X {
		maxX = t.C.X
	}

	if t.A.Y <= t.B.Y && t.A.Y <= t.C.Y {
		minY = t.A.Y
	} else if t.B.Y <= t.A.Y && t.B.Y <= t.C.Y {
		minY = t.B.Y
	} else if t.C.Y <= t.A.Y && t.C.Y <= t.B.Y {
		minY = t.C.Y
	}

	if t.A.Y >= t.B.Y && t.A.Y >= t.C.Y {
		maxY = t.A.Y
	} else if t.B.Y >= t.A.Y && t.B.Y >= t.C.Y {
		maxY = t.B.Y
	} else if t.C.Y >= t.A.Y && t.C.Y >= t.B.Y {
		maxY = t.C.Y
	}

	return BoundingBox{
		MinX: minX,
		MinY: minY,

		MaxX: maxX,
		MaxY: maxY,
	}
}

func signedTriangleArea(triangle gmath.Triangle[int32]) float64 {
	a := (float64(triangle.B.Y) - float64(triangle.A.Y)) * (float64(triangle.B.X) + float64(triangle.A.X))
	b := (float64(triangle.C.Y) - float64(triangle.B.Y)) * (float64(triangle.C.X) + float64(triangle.B.X))
	c := (float64(triangle.A.Y) - float64(triangle.C.Y)) * (float64(triangle.A.X) + float64(triangle.C.X))
	return 0.5 * (a + b + c)
}

func rasterize(fb *window.Framebuffer, triangles []gmath.Triangle[int32]) {
	// Sort seems unnecessary for this method
	// sort.Slice(triangles, func(i, j int) bool {
	// 	closest := triangles[i]
	// 	b := triangles[j]
	// 	if (closest.A.Y < b.A.Y) || (closest.B.Y < b.B.Y) || (closest.C.Y < b.C.Y) {
	// 		closest = b
	// 	}
	// 	return closest == b
	// })

	// For every sorted triangle points
	// WIP: Still not rendering triangle correctly, need to check if bounding box are correct
	for _, t := range triangles {
		boundingBox := findBoundingBox(t)
		totalArea := signedTriangleArea(t)
		// Iterate over all pixels on the screen
		for x := boundingBox.MinY; x < boundingBox.MaxY; x++ {
			for y := boundingBox.MinX; y < boundingBox.MaxX; y++ {
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
				fb.SetPixel(y, x, t.Color)
			}
		}
	}
}

func draw(ctx *window.WindowContext) {
	// render
	ctx.Framebuffer.Clear(0xFF000000) // black

	a := gmath.Vector3[int32]{
		X: 100,
		Y: 100,
		Z: 0,
	}
	b := gmath.Vector3[int32]{
		X: 200,
		Y: 100,
		Z: 0,
	}
	c := gmath.Vector3[int32]{
		X: 200,
		Y: 300,
		Z: 0,
	}

	// a2 := gmath.Vector3[int32]{
	// 	X: 300,
	// 	Y: 200,
	// 	Z: 0,
	// }
	// b2 := gmath.Vector3[int32]{
	// 	X: 340,
	// 	Y: 200,
	// 	Z: 0,
	// }
	// c2 := gmath.Vector3[int32]{
	// 	X: 500,
	// 	Y: 400,
	// 	Z: 0,
	// }
	t1 := gmath.NewTriangle(a, b, c, window.ColorBlue)
	// t2 := gmath.NewTriangle(a2, b2, c2, window.ColorMagenta)
	rasterize(ctx.Framebuffer, []gmath.Triangle[int32]{t1})

	// blit framebuffer to texture
	ctx.Texture.Update(nil, unsafe.Pointer(&ctx.Framebuffer.Pixels[0]), int(ctx.Framebuffer.Width)*4)
	ctx.Renderer.Clear()
	ctx.Renderer.Copy(ctx.Texture, nil, nil)
	ctx.Renderer.Present()
}

func Run() {
	win := window.New("Test", screenWidth, screenHeight)
	defer win.Destroy()
	win.SetDrawFunction(draw)

	win.Run()
}
