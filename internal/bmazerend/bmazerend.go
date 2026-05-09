package bmazerend

import (
	"fmt"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

type framebuffer struct {
	pixels []uint32
	width  int32
	height int32
}

func newFramebuffer(w, h int32) *framebuffer {
	return &framebuffer{
		pixels: make([]uint32, w*h),
		width:  w,
		height: h,
	}
}

func (fb *framebuffer) clear(color uint32) {
	for i := range fb.pixels {
		fb.pixels[i] = color
	}
}

func (fb *framebuffer) setPixel(x, y int32, color uint32) {
	if x < 0 || x >= fb.width || y < 0 || y >= fb.height {
		return
	}
	fb.pixels[y*fb.width+x] = color
}

func Run() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(fmt.Sprintf("sdl init: %v", err))
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"boomer-maze",
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		screenWidth, screenHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(fmt.Sprintf("create window: %v", err))
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(fmt.Sprintf("create renderer: %v", err))
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_STREAMING,
		screenWidth, screenHeight,
	)
	if err != nil {
		panic(fmt.Sprintf("create texture: %v", err))
	}
	defer texture.Destroy()

	fb := newFramebuffer(screenWidth, screenHeight)

	running := true
	for running {
		// input
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if e.Keysym.Sym == sdl.K_ESCAPE {
					running = false
				}
			}
		}

		// render
		fb.clear(0xFF000000)                                   // black
		fb.setPixel(screenWidth/2, screenHeight/2, 0xFFFF0000) // red dot

		// blit framebuffer to texture
		texture.Update(nil, unsafe.Pointer(&fb.pixels[0]), int(fb.width)*4)
		renderer.Clear()
		renderer.Copy(texture, nil, nil)
		renderer.Present()
	}
}
