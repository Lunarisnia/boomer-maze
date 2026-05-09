package bmazerend

import (
	"unsafe"

	"github.com/lunarisnia/boomer-maze/internal/bmazerend/window"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func draw(ctx *window.WindowContext) {
	// render
	ctx.Framebuffer.Clear(0xFF000000)                     // black
	ctx.Framebuffer.SetPixel(100, 100, window.ColorWhite) // red dot
	ctx.Framebuffer.SetPixel(200, 100, window.ColorWhite) // red dot
	ctx.Framebuffer.SetPixel(200, 400, window.ColorWhite) // red dot

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
