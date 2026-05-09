package bmazerend

import (
	"fmt"
	"math"
	"math/rand"
	"time"
	"unsafe"

	"github.com/lunarisnia/boomer-maze/internal/bmazerend/window"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

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
	for x := ax; x <= bx; x++ {
		if steep {
			fb.SetPixel(int32(y), x, color)
		} else {
			fb.SetPixel(x, int32(y), color)
		}
		y += float64(by-ay) / float64(bx-ax)
	}
}

var lines []func(ctx *window.WindowContext)

func draw(ctx *window.WindowContext) {
	// render
	ctx.Framebuffer.Clear(0xFF000000) // black

	// line(100, 100, 220, 120, ctx.Framebuffer, window.ColorGreen)
	// line(220, 120, 200, 400, ctx.Framebuffer, window.ColorGreen)
	// line(200, 400, 100, 100, ctx.Framebuffer, window.ColorGreen)
	start := time.Now()
	for _, l := range lines {
		l(ctx)
	}

	// blit framebuffer to texture
	ctx.Texture.Update(nil, unsafe.Pointer(&ctx.Framebuffer.Pixels[0]), int(ctx.Framebuffer.Width)*4)
	ctx.Renderer.Clear()
	ctx.Renderer.Copy(ctx.Texture, nil, nil)
	ctx.Renderer.Present()
	fmt.Println("RenderTime: ", time.Since(start))
}

func Run() {
	win := window.New("Test", screenWidth, screenHeight)
	defer win.Destroy()
	win.SetDrawFunction(draw)

	for range 10_000 {
		ax := rand.Int31() % screenWidth
		ay := rand.Int31() % screenHeight
		bx := rand.Int31() % screenWidth
		by := rand.Int31() % screenHeight
		col := rand.Uint32()
		lines = append(lines, func(ctx *window.WindowContext) {
			line(ax, ay, bx, by, ctx.Framebuffer, col)
		})
	}

	win.Run()
}
