package window

type Framebuffer struct {
	Pixels []uint32
	Width  int32
	Height int32
}

func newFramebuffer(w, h int32) *Framebuffer {
	return &Framebuffer{
		Pixels: make([]uint32, w*h),
		Width:  w,
		Height: h,
	}
}

func (fb *Framebuffer) Clear(color uint32) {
	for i := range fb.Pixels {
		fb.Pixels[i] = color
	}
}

func (fb *Framebuffer) SetPixel(x, y int32, color uint32) {
	if x < 0 || x >= fb.Width || y < 0 || y >= fb.Height {
		return
	}
	fb.Pixels[y*fb.Width+x] = color
}
