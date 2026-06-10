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
	for _, t := range triangles {
		// Iterate over all pixels on the screen
		for j := range int32(screenHeight) {
			for i := range int32(screenWidth) {
				// Check if point is inside the triangle
				if gmath.IsInsideTriangle(t.A, t.B, t.C, gmath.Vector3[int32]{
					X: i,
					Y: j,
					Z: 0,
				}) {
					fb.SetPixel(i, j, t.Color)
				}
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

	a2 := gmath.Vector3[int32]{
		X: 300,
		Y: 50,
		Z: 0,
	}
	b2 := gmath.Vector3[int32]{
		X: 400,
		Y: 200,
		Z: 0,
	}
	c2 := gmath.Vector3[int32]{
		X: 50,
		Y: 200,
		Z: 0,
	}
	t1 := gmath.NewTriangle(a, b, c, window.ColorRed)
	t2 := gmath.NewTriangle(a2, b2, c2, window.ColorMagenta)
	rasterize(ctx.Framebuffer, []gmath.Triangle[int32]{t1, t2})

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
