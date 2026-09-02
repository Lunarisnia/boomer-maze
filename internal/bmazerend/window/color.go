package window

const (
	ColorBlack   = 0xFF_00_00_00
	ColorWhite   = 0xFF_FF_FF_FF
	ColorRed     = 0xFF_FF_00_00
	ColorGreen   = 0xFF_00_FF_00
	ColorBlue    = 0xFF_00_00_FF
	ColorYellow  = 0xFF_FF_FF_00
	ColorCyan    = 0xFF_00_FF_FF
	ColorMagenta = 0xFF_FF_00_FF
	ColorGray    = 0xFF_80_80_80
)

func ARGB(a, r, g, b uint8) uint32 {
	return uint32(a)<<24 | uint32(r)<<16 | uint32(g)<<8 | uint32(b)
}
