package characters

import (
	"fmt"
	"path/filepath"

	"github.com/disintegration/gift"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak/render"
)

type CharacterClass interface {
	Special()
	loadAnimationMap() map[string]render.Modifiable
}

type Swordsman struct{}
type Archer struct{}

func (s *Swordsman) Special() {
	fmt.Println("Attacking!")
}

func (s *Swordsman) loadAnimationMap() map[string]render.Modifiable {

	animFilePath := (filepath.Join("16x32", "warrior.png"))

	walkRT, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 0, 2, 0, 0, 0}...)
	dlog.ErrorCheck(err)
	walkLT := walkRT.Copy().Modify(mod.FlipX)

	return map[string]render.Modifiable{

		"walkRT": render.NewReverting(walkRT),
		"walkLT": render.NewReverting(walkLT),
	}
}

func (s *Archer) Special() {
	fmt.Println("Attacking! As an archer")
}

func (s *Archer) loadAnimationMap() map[string]render.Modifiable {

	animFilePath := (filepath.Join("16x32", "warrior.png"))

	walkRT, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 0, 2, 0, 0, 0}...)
	dlog.ErrorCheck(err)
	walkRT.Modify(mod.GiftTransform(gift.ColorBalance(100, 0, 0)))
	walkLT := walkRT.Copy().Modify(mod.FlipX)

	return map[string]render.Modifiable{

		"walkRT": render.NewReverting(walkRT),
		"walkLT": render.NewReverting(walkLT),
	}
}
