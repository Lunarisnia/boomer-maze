package window

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Window interface {
	Run()
	SetDrawFunction(DrawFunc)
	Destroy()
}

type DrawFunc func(ctx *WindowContext)

type WindowContext struct {
	Framebuffer *Framebuffer
	Renderer    *sdl.Renderer
	Texture     *sdl.Texture
}

type window struct {
	window      *sdl.Window
	renderer    *sdl.Renderer
	texture     *sdl.Texture
	framebuffer *Framebuffer

	drawFunc DrawFunc
}

func New(title string, width int32, height int32) Window {
	var err error
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(fmt.Sprintf("sdl init: %v", err))
	}

	w := window{
		drawFunc: func(ctx *WindowContext) {},
	}

	w.window, err = sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		width, height,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(fmt.Sprintf("create window: %v", err))
	}

	w.renderer, err = sdl.CreateRenderer(w.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(fmt.Sprintf("create renderer: %v", err))
	}

	w.texture, err = w.renderer.CreateTexture(
		sdl.PIXELFORMAT_ARGB8888,
		sdl.TEXTUREACCESS_STREAMING,
		width, height,
	)
	if err != nil {
		panic(fmt.Sprintf("create texture: %v", err))
	}

	w.framebuffer = newFramebuffer(width, height)
	return &w
}

func (w *window) Run() {
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

		ctx := &WindowContext{
			Framebuffer: w.framebuffer,
			Renderer:    w.renderer,
			Texture:     w.texture,
		}
		w.drawFunc(ctx)
	}
}

func (w *window) SetDrawFunction(drawFunc DrawFunc) {
	w.drawFunc = drawFunc
}

func (w *window) Destroy() {
	defer sdl.Quit()
	defer w.window.Destroy()
	defer w.renderer.Destroy()
	defer w.texture.Destroy()
}
