package bmazerend

import (
	"fmt"
	"log"
	"math/rand/v2"
	"unsafe"

	"github.com/lunarisnia/boomer-maze/internal/bmazerend/window"
	"github.com/lunarisnia/boomer-maze/internal/gmath"
	"github.com/lunarisnia/boomer-maze/internal/wavefront"
)

var loadedModel *wavefront.Object
var faceColors []uint32
var angle float64 = 0.0

const (
	screenWidth  = 800
	screenHeight = 600
)

// NOTE: A stub function
func triangle(fb *window.Framebuffer, t gmath.Triangle[int32]) {
	fmt.Println("Does jack shit")
}

func loadRotatingFace() {
	model, err := wavefront.LoadModel("./models/african_head.obj")
	if err != nil {
		log.Fatalln(err)
	}
	loadedModel = model
	for range len(model.Faces) {
		faceColors = append(faceColors, rand.Uint32())
	}
}

func drawRotatingFace(ctx *window.WindowContext) {
	// render
	ctx.Framebuffer.Clear(0xFF000000) // black

	for i, indice := range loadedModel.Faces {
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

		ax, ay, aVisible := gmath.LocalToScreen(aLocal, pos, screenWidth, screenHeight, 60)
		bx, by, bVisible := gmath.LocalToScreen(bLocal, pos, screenWidth, screenHeight, 60)
		cx, cy, cVisible := gmath.LocalToScreen(cLocal, pos, screenWidth, screenHeight, 60)
		if !aVisible || !bVisible || !cVisible {
			continue
		}

		triangle(ctx.Framebuffer, gmath.NewTriangle(
			gmath.Vector3[int32]{X: int32(ay), Y: int32(ax), Z: 0},
			gmath.Vector3[int32]{X: int32(by), Y: int32(bx), Z: 0},
			gmath.Vector3[int32]{X: int32(cy), Y: int32(cx), Z: 0},
			faceColors[i],
		))
	}
	angle += 0.001

	// blit framebuffer to texture
	ctx.Texture.Update(nil, unsafe.Pointer(&ctx.Framebuffer.Pixels[0]), int(ctx.Framebuffer.Width)*4)
	ctx.Renderer.Clear()
	ctx.Renderer.Copy(ctx.Texture, nil, nil)
	ctx.Renderer.Present()
}
