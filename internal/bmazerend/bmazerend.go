package bmazerend

import (
	"log"
	"math"
	"math/rand"
	"unsafe"

	"github.com/lunarisnia/boomer-maze/internal/bmazerend/window"
	"github.com/lunarisnia/boomer-maze/internal/gmath"
	"github.com/lunarisnia/boomer-maze/internal/valueutil"
	"github.com/lunarisnia/boomer-maze/internal/wavefront"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var loadedModel *wavefront.Object
var faceColors []uint32
var angle float64 = 0.0

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
				fb.SetPixel(y, x, t.Color)
			}
		}
	}
}

func triangle(fb *window.Framebuffer, t gmath.Triangle[int32]) {
	rasterize(fb, []gmath.Triangle[int32]{t})
}

func draw(ctx *window.WindowContext) {
	// render
	ctx.Framebuffer.Clear(0xFF000000) // black

	for i, indice := range loadedModel.Faces {
		aLocal := loadedModel.Vertice[indice.X]
		bLocal := loadedModel.Vertice[indice.Y]
		cLocal := loadedModel.Vertice[indice.Z]
		pos := gmath.Vector3[float64]{
			X: 0.0,
			Y: 0.0,
			Z: 1.5,
		}

		aLocal = gmath.RotateY(aLocal, angle)
		bLocal = gmath.RotateY(bLocal, angle)
		cLocal = gmath.RotateY(cLocal, angle)

		ax, ay, aVisible := gmath.LocalToScreen(aLocal, pos, screenWidth, screenHeight, 60)
		bx, by, bVisible := gmath.LocalToScreen(bLocal, pos, screenWidth, screenHeight, 60)
		cx, cy, cVisible := gmath.LocalToScreen(cLocal, pos, screenWidth, screenHeight, 60)
		if !aVisible || !bVisible || !cVisible {
			continue
		}

		triangle(ctx.Framebuffer, gmath.NewTriangle(
			gmath.Vector3[int32]{X: int32(ay), Y: int32(ax), Z: 0},
			gmath.Vector3[int32]{X: int32(by), Y: int32(bx), Z: 0},
			gmath.Vector3[int32]{X: int32(cy), Y: int32(cx), Z: 0},
			faceColors[i],
		))
	}
	angle += 0.001

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

	model, err := wavefront.LoadModel("./models/african_head.obj")
	if err != nil {
		log.Fatalln(err)
	}
	loadedModel = model
	for range len(model.Faces) {
		faceColors = append(faceColors, rand.Uint32())
	}

	win.Run()
}

// TODO: Next goal: Understand why backface culling hack worked
