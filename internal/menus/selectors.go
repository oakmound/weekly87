package menus

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/entities/x/btn/grid"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/alg/floatgeom"
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
			sp := button.GetSpace()
			sp.Location.Max = sp.Location.Max.Add(floatgeom.Point3{1,1,0})
			spcs = append(spcs, sp)
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
		s.Display(func(pt floatgeom.Point2) render.Renderable {
			poly, err := render.NewPolygon(
				floatgeom.Point2{0, 0},
				floatgeom.Point2{pt.X(), 0},
				floatgeom.Point2{pt.X(), pt.Y()},
				floatgeom.Point2{0, pt.Y()},
			)
			dlog.ErrorCheck(err)
			return poly.GetThickOutline(Gold, 2)
		}),
	)
}
