package menus

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/entities/x/btn/grid"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/mouse"
	s "github.com/oakmound/weekly87/internal/menus/selector"
	"github.com/oakmound/weekly87/internal/sfx"
)

// ButtonSelectorA sets up a standard button selector for menus
func ButtonSelectorA(selectors grid.Grid) s.Option {
	btnList := []btn.Btn{}
	spcs := []*collision.Space{}
	for _, selectList := range selectors {
		btnList = append(btnList, selectList...)
		for _, button := range selectList {
			spcs = append(spcs, button.GetSpace())
		}
	}
	return ButtonSelectorSpacesA(spcs, btnList)
}

// ButtonSelectorSpacesA sets up a standard button selector for menus given a set of spaces and btns
// For times when the menu is not in a grid.Grid
func ButtonSelectorSpacesA(spcs []*collision.Space, btnList []btn.Btn) s.Option {
	return s.And(
		s.Layers(2, 3),
		s.VertArrowControl(),
		s.JoystickVertDpadControl(),
		s.Spaces(spcs...),
		s.Callback(func(i int, _ ...interface{}) {
			sfx.Play("selected")
			btnList[i].Trigger(mouse.ClickOn, nil)
		}),
		s.SelectTrigger(key.Down+key.Spacebar),
		s.SelectTrigger("A"+joystick.ButtonUp),
		s.SelectTrigger("Start"+joystick.ButtonUp),
		s.DestroyTrigger(key.Down+key.Escape),
		s.Wraps(true),
	)
}
