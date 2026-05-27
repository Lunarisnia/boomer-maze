package bmazerend

import (
	"log"
	"math"
	"unsafe"

	"github.com/lunarisnia/boomer-maze/internal/bmazerend/window"
	"github.com/lunarisnia/boomer-maze/internal/gmath"
	"github.com/lunarisnia/boomer-maze/internal/wavefront"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var loadedModel *wavefront.Object
var angle float64 = 0.0

func swap[T any](a T, b T) (T, T) {
	return b, a
}

func line(ax int32, ay int32, bx int32, by int32, fb *window.Framebuffer, color uint32) {
	steep := math.Abs(float64(ax-bx)) < math.Abs(float64(ay-by))
	if steep {
		ax, ay = swap(ax, ay)
		bx, by = swap(bx, by)
	}
	if ax > bx {
		ax, bx = swap(ax, bx)
		ay, by = swap(ay, by)
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

func draw(ctx *window.WindowContext) {
	// render
	ctx.Framebuffer.Clear(0xFF000000) // black

	for _, indice := range loadedModel.Faces {
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

		ax, ay, _ := gmath.LocalToScreen(aLocal, pos, screenWidth, screenHeight, 60)
		bx, by, _ := gmath.LocalToScreen(bLocal, pos, screenWidth, screenHeight, 60)
		cx, cy, _ := gmath.LocalToScreen(cLocal, pos, screenWidth, screenHeight, 60)
		line(int32(ax), int32(ay), int32(bx), int32(by), ctx.Framebuffer, window.ColorRed)
		line(int32(bx), int32(by), int32(cx), int32(cy), ctx.Framebuffer, window.ColorRed)
		line(int32(cx), int32(cy), int32(ax), int32(ay), ctx.Framebuffer, window.ColorRed)
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

	model, err := wavefront.LoadModel("./models/diablo3_pose.obj")
	if err != nil {
		log.Fatalln(err)
	}
	loadedModel = model

	win.Run()
}
