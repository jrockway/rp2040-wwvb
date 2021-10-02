package screen

import (
	"fmt"
	"image"
	"image/color"

	"github.com/mdp/smallfont"
)

type SetTexter interface {
	SetText(string)
}

type Displayer interface {
	Display() error
	SetPixel(int16, int16, color.RGBA)
}

type Screen struct {
	Display Displayer
	n       int
	buf     [64]byte
}

func (s *Screen) Printf(x string, args ...interface{}) {
	str := fmt.Sprintf(x, args...)
	for _, c := range str {
		if c == '\n' {
			// If we're at the end of the line, doing nothing makes a new line.
			if s.n%16 != 0 {
				eol := 16*(1+(s.n/16)) - s.n
				for j := 0; j < eol; j++ {
					s.buf[s.n] = ' '
					s.n = (s.n + 1) % len(s.buf)
				}
			}
			continue
		}
		s.buf[s.n] = byte(c)
		s.n = (s.n + 1) % len(s.buf)
	}
	if d, ok := s.Display.(SetTexter); ok {
		d.SetText(string(s.buf[:]))
	}
	img := image.NewRGBA(image.Rect(0, 0, 128, 32))
	sf := smallfont.Context{
		Dst:    img,
		StartX: 0,
		StartY: 0,
		Font:   smallfont.Font8x8,
		Color:  image.White,
	}
	sf.Draw(s.buf[0:16], 0, 0)
	sf.Draw(s.buf[16:32], 0, 8)
	sf.Draw(s.buf[32:48], 0, 16)
	sf.Draw(s.buf[48:64], 0, 24)
	for x := 0; x < 128; x++ {
		for y := 0; y < 32; y++ {
			s.Display.SetPixel(int16(x), int16(y), img.RGBAAt(x, y))
		}
	}
	s.Display.Display()
}

func (s *Screen) Clear() {
	for i := 0; i < len(s.buf); i++ {
		s.buf[i] = 0
	}
	s.n = 0
}
