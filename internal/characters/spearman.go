package characters

import (
	"fmt"
	"path/filepath"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

var _ Player = &Spearman{}

type Spearman struct {
	basePlayer
}

func (s *Spearman) Init() event.CID {
	return event.NextID(s)
}

func NewSpearman(x, y float64) *Spearman {
	s := &Spearman{}
	r := render.NewSwitch("walkRT", s.loadAnimationMap())
	s.Interactive = entities.NewInteractive(x, y, playerWidth, playerHeight, r, nil, s.Init(), 0)
	collision.Add(s.RSpace.Space)
	s.Speed = physics.NewVector(0, 5)
	s.alive = true

	// h.animation = ch.R.(*render.Compound)
	return s
}

func (s *Spearman) Special1() {
	fmt.Println("Attacking!")
}

func (s *Spearman) loadAnimationMap() map[string]render.Modifiable {

	animFilePath := (filepath.Join("16x32", "warrior.png"))

	walkRT, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 0, 2, 0, 0, 0}...)
	dlog.ErrorCheck(err)
	walkLT := walkRT.Copy().Modify(mod.FlipX)

	return map[string]render.Modifiable{
		"walkRT": render.NewReverting(walkRT),
		"walkLT": render.NewReverting(walkLT),
	}
}
