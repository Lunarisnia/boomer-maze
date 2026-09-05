package window

import (
	"bytes"
	"image"
	"testing"
)

func TestCopyRGBA(t *testing.T) {
	fb := newFramebuffer(2, 2)
	fb.Pixels = []uint32{0xFFFF0000, 0xFF00FF00, 0xFF0000FF, 0x78123456}
	dst := image.NewRGBA(image.Rect(0, 0, 2, 2))
	fb.copyRGBA(dst)
	want := []byte{255, 0, 0, 255, 0, 255, 0, 255, 0, 0, 255, 255, 0x12, 0x34, 0x56, 0x78}
	if !bytes.Equal(dst.Pix, want) {
		t.Fatalf("RGBA pixels = %v, want %v", dst.Pix, want)
	}
}
