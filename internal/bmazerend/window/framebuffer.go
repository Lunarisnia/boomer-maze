package window

import "image"

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

func (fb *Framebuffer) copyRGBA(dst *image.RGBA) {
	for y := 0; y < int(fb.Height); y++ {
		for x := 0; x < int(fb.Width); x++ {
			pixel := fb.Pixels[y*int(fb.Width)+x]
			i := y*dst.Stride + x*4
			dst.Pix[i] = uint8(pixel >> 16)
			dst.Pix[i+1] = uint8(pixel >> 8)
			dst.Pix[i+2] = uint8(pixel)
			dst.Pix[i+3] = uint8(pixel >> 24)
		}
	}
}

func (fb *Framebuffer) SetPixel(x, y int32, color uint32) {
	if x < 0 || x >= fb.Width || y < 0 || y >= fb.Height {
		return
	}
	fb.Pixels[y*fb.Width+x] = color
}
