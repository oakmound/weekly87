package run

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
)

var (
	baseSeed       int64
	currentSection int64
)

// Todo generate section

type Section struct {
	// Todo: these should be at the generated level,
	// made into a render.Composite, then converted into
	// a sprite via composite.ToSprite
	ground   [64][24]render.Renderable
	wall     [64][12]render.Renderable
	entities []characters.Character
}

func TestSection() *Section {
	s := &Section{}
	for x := 0; x < len(s.ground); x++ {
		for y := 0; y < len(s.ground[x]); y++ {
			s.ground[x][y] = render.NewColorBox(
				16, 16, color.RGBA{0, uint8(rand.Intn(125)), 255, 255},
			)
		}
	}
	for x := 0; x < len(s.wall); x++ {
		for y := 0; y < len(s.wall[x]); y++ {
			s.wall[x][y] = render.NewColorBox(
				16, 16, color.RGBA{0, uint8(rand.Intn(10)), 120, 255},
			)
		}
	}
	ch := characters.NewChest(0)
	ch.SetPos(800, 400)
	s.entities = append(s.entities, ch)
	return s
}

func (s *Section) Draw() {
	groundOffset := float64(oak.ScreenHeight) * 1 / 3
	for x := 0; x < len(s.ground); x++ {
		for y := 0; y < len(s.ground[x]); y++ {
			s.ground[x][y].SetPos(float64(x)*16, groundOffset+float64(y)*16)
			render.Draw(s.ground[x][y], 0)
		}
	}
	for x := 0; x < len(s.wall); x++ {
		for y := 0; y < len(s.wall[x]); y++ {
			s.wall[x][y].SetPos(float64(x)*16, float64(y)*16)
			render.Draw(s.wall[x][y], 0)
		}
	}
	for _, e := range s.entities {
		render.Draw(e.GetRenderable(), 2, 1)
	}
	// Worry about entities later
}

func (s *Section) Shift(shift float64) {
	for x := 0; x < len(s.ground); x++ {
		for y := 0; y < len(s.ground[x]); y++ {
			s.ground[x][y].ShiftX(shift)
		}
	}
	for x := 0; x < len(s.wall); x++ {
		for y := 0; y < len(s.wall[x]); y++ {
			s.wall[x][y].ShiftX(shift)
		}
	}
	for _, e := range s.entities {
		ShiftMoverX(e, shift)
		fmt.Println(e.Vec().X(), e.Vec().Y())
	}
}
