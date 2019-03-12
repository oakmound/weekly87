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
type Medic struct{}

func (s *Swordsman) Special1() {
	fmt.Println("Attacking!")
}

func (s *Swordsman) Special2() {
	fmt.Println("Attacking!2")
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

func (s *Archer) Special1() {
	fmt.Println("Attacking! As an archer")
}
func (s *Archer) Special2() {
	fmt.Println("Attacking! As an archer2")
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

func (s *Medic) Special1() {
	fmt.Println("help1")
}
func (s *Medic) Special2() {
	fmt.Println("Heal2")
}

func (s *Medic) loadAnimationMap() map[string]render.Modifiable {

	animFilePath := (filepath.Join("16x32", "warrior.png"))

	walkRT, err := render.LoadSheetSequence(animFilePath, 16, 32, 0, 8, []int{1, 0, 2, 0, 0, 0}...)
	dlog.ErrorCheck(err)
	walkRT.Modify(mod.GiftTransform(gift.ColorBalance(0, 200, 0)))
	walkLT := walkRT.Copy().Modify(mod.FlipX)

	return map[string]render.Modifiable{

		"walkRT": render.NewReverting(walkRT),
		"walkLT": render.NewReverting(walkLT),
	}
}
