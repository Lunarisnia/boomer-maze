package bmazerend

import (
	"math"
	"unsafe"

	"github.com/lunarisnia/boomer-maze/internal/bmazerend/window"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func line(ax int32, ay int32, bx int32, by int32, fb *window.Framebuffer, color uint32) {
	for t := 0.0; t < 1.0; t += 0.002 {
		x := int32(math.Round(float64(ax) + float64(bx-ax)*t))
		y := int32(math.Round(float64(ay) + float64(by-ay)*t))
		fb.SetPixel(x, y, color)
	}
}

func draw(ctx *window.WindowContext) {
	// render
	ctx.Framebuffer.Clear(0xFF000000) // black

	line(100, 100, 220, 120, ctx.Framebuffer, window.ColorGreen)
	line(220, 120, 200, 400, ctx.Framebuffer, window.ColorGreen)
	line(200, 400, 100, 100, ctx.Framebuffer, window.ColorGreen)

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
