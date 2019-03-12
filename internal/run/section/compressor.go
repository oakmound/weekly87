package section

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
)

type compressor struct {
	ground   [70][24]render.Modifiable
	wall     [70][12]render.Modifiable
	entities []characters.Character
}

func (sg *compressor) generate() *Section {

	// Place ground and wall appropariately in composites and
	// create sprites
	s := &Section{}
	groundOffset := float64(oak.ScreenHeight) * 1 / 3
	gcmp1 := render.NewCompositeM()
	for x, col := range sg.ground {
		for y, r := range col {
			r.SetPos(float64(x)*16, groundOffset+float64(y)*16)
			gcmp1.Append(r)
		}
	}
	s.ground = gcmp1.ToSprite()
	wcmp := render.NewCompositeM()
	for x, col := range sg.wall {
		for y, r := range col {
			r.SetPos(float64(x)*16, float64(y)*16)
			wcmp.Append(r)
		}
	}
	s.wall = wcmp.ToSprite()
	// Todo: attach all entities at offsets?
	s.wall.Vector = s.wall.AttachX(s.ground, 0)
	s.entities = make([]characters.Character, len(sg.entities))
	copy(s.entities, sg.entities)
	sg.entities = nil
	return s
}
