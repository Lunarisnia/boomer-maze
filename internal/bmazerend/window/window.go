package window

import (
	"image"

	"github.com/AllenDang/giu"
)

type Window interface {
	Run()
	SetDrawFunction(DrawFunc)
}

type DrawFunc func(ctx *WindowContext)

type WindowContext struct {
	Framebuffer *Framebuffer
}

type window struct {
	window      *giu.MasterWindow
	framebuffer *Framebuffer
	drawFunc    DrawFunc
}

func New(title string, width int32, height int32) Window {
	w := &window{
		window:      giu.NewMasterWindow(title, int(width), int(height), 0),
		framebuffer: newFramebuffer(width, height),
		drawFunc:    func(ctx *WindowContext) {},
	}
	w.window.SetTargetFPS(60)
	w.window.RegisterKeyboardShortcuts(giu.WindowShortcut{Key: giu.KeyEscape, Callback: w.window.Close})
	return w
}

func (w *window) Run() {
	surface := image.NewRGBA(image.Rect(0, 0, int(w.framebuffer.Width), int(w.framebuffer.Height)))
	texture := &giu.ReflectiveBoundTexture{Surface: surface}

	backend := giu.Context.Backend().(*giu.GLFWBackend)
	beforeDestroy := backend.BeforeDestroyHook()
	backend.SetBeforeDestroyContextHook(func() {
		texture.ForceRelease()
		if beforeDestroy != nil {
			beforeDestroy()
		}
	})

	ctx := &WindowContext{Framebuffer: w.framebuffer}
	w.window.Run(func() {
		w.drawFunc(ctx)
		w.framebuffer.copyRGBA(surface)
		giu.Window("Renderer").Pos(10, 30).Size(float32(w.framebuffer.Width)+32, float32(w.framebuffer.Height)+64).Layout(
			giu.Custom(func() {
				width, height := giu.GetAvailableRegion()

				scale := min(width/float32(w.framebuffer.Width), height/float32(w.framebuffer.Height))
				if scale > 0 {
					texture.ToImageWidget().Size(float32(w.framebuffer.Width)*scale, float32(w.framebuffer.Height)*scale).Build()
				}
			}),
		)
		giu.Update()
	})
}

func (w *window) SetDrawFunction(drawFunc DrawFunc) {
	w.drawFunc = drawFunc
}
