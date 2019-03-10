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
	ground   *render.Sprite
	wall     *render.Sprite
	entities []characters.Character
}

type sectionGenerator struct {
	// Todo: these should be at the generated level,
	// made into a render.Composite, then converted into
	// a sprite via composite.ToSprite
	ground   [64][24]render.Modifiable
	wall     [64][12]render.Modifiable
	entities []characters.Character
}

func (sg *sectionGenerator) generate() *Section {
	groundOffset := float64(oak.ScreenHeight) * 1 / 3
	// Place ground and wall appropariately in composites and
	// create sprites
	s := &Section{}
	gcmp := render.NewCompositeM()
	for x, col := range sg.ground {
		for y, r := range col {
			r.SetPos(float64(x)*16, groundOffset+float64(y)*16)
			gcmp.Append(r)
		}
	}
	s.ground = gcmp.ToSprite()
	wcmp := render.NewCompositeM()
	for x, col := range sg.wall {
		for y, r := range col {
			r.SetPos(float64(x)*16, float64(y)*16)
			wcmp.Append(r)
		}
	}
	s.wall = wcmp.ToSprite()
	s.entities = sg.entities
	return s
}

func TestSection() *Section {
	s := &sectionGenerator{}
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
	return s.generate()
}

func (s *Section) Draw() {
	render.Draw(s.ground)
	render.Draw(s.wall)
	for _, e := range s.entities {
		render.Draw(e.GetRenderable(), 2, 1)
	}
}

func (s *Section) Shift(shift float64) {
	s.ground.ShiftX(shift)
	s.wall.ShiftX(shift)
	for _, e := range s.entities {
		ShiftMoverX(e, shift)
		fmt.Println(e.Vec().X(), e.Vec().Y())
	}
}
