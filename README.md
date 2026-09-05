# Boomer Maze

A Go software rasterizer displayed inside a giu (Dear ImGui) window.

## Run

Requires Go 1.25 or newer and a C/C++ toolchain with CGO enabled. On macOS,
install the Xcode Command Line Tools with `xcode-select --install` if needed.
The default giu backend uses GLFW/OpenGL; SDL2 is no longer required.

```sh
go run .
```

The framebuffer keeps its 800x600 rendering resolution and scales to fit the
resizable `Renderer` subwindow without changing its aspect ratio. Press Escape
or close the main window to exit.

## Development

```sh
go test ./...
go build ./...
```

The rasterizer writes packed ARGB pixels in `internal/bmazerend/bmazerend.go`.
The window layer in `internal/bmazerend/window/window.go` converts them to RGBA
and displays a giu texture inside the `Renderer` subwindow. Add UI widgets to
its `giu.Window("Renderer").Layout`.
The texture is refreshed when its image changes; the UI targets 60 FPS.
