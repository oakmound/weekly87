package menus

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/entities/x/btn/grid"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/mouse"
	s "github.com/oakmound/weekly87/internal/menus/selector"
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
	return s.And(
		s.Layers(2, 3),
		s.HorzArrowControl(),
		s.Spaces(spcs...),
		s.Callback(func(i int) {
			btnList[i].Trigger(mouse.ClickOn, nil)
		}),
		s.SelectTrigger(key.Down+key.Spacebar),
		s.DestroyTrigger(key.Down+key.Escape),
	)
}
