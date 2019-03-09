package characters

import (
	"fmt"
	"path/filepath"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type Spearman struct {
	*entities.Moving
}

func (s *Spearman) Init() event.CID {
	return event.NextID(s)
}

func NewSpearman(x, y float64) *Spearman {
	s := &Spearman{}
	// r := render.NewColorBox(playerWidth, playerHeight, color.RGBA{255, 0, 0, 255})
	r := render.NewSwitch("walkRT", s.loadAnimationMap())
	s.Moving = entities.NewMoving(x, y, playerWidth, playerHeight, r, nil, s.Init(), 0)
	s.Speed = physics.NewVector(0, 5)

	// s.R = render.NewCompoundR("walkRT", s.loadAnimationMap())
	// h.animation = ch.R.(*render.Compound)
	return s
}

func (s *Spearman) Attack1() {
	fmt.Println("Attacking!")
}

func (s *Spearman) loadAnimationMap() map[string]render.Modifiable {

	animFilePath := (filepath.Join("16x32", "warrior.png"))

	walkRT, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 0, 2, 0}...)
	dlog.ErrorCheck(err)
	walkLT := walkRT.Copy().Modify(mod.FlipX)

	return map[string]render.Modifiable{

		"walkRT": render.NewReverting(walkRT),
		"walkLT": render.NewReverting(walkLT),
	}
}
