package menus

import (
	"github.com/oakmound/oak/alg/intgeom"
	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"

	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/entities/x/mods"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render/mod"
	"golang.org/x/image/colornames"

	"image/color"
)

// Shared menu configuration parameters
var (
	BtnHeightA = 30.0
	BtnWidthA  = 120.0
	BtnCfgA    = btn.And(
		btn.Width(BtnWidthA),
		btn.Height(BtnHeightA),
		btn.Mod(mods.HighlightOff(colornames.Blue, 3, 0, 0)),
		btn.Mod(mods.InnerHighlightOff(colornames.Black, 1, 0, 0)),
		btn.TxtOff(BtnWidthA/4, BtnHeightA/3), //magic numbers from main menu
	)
	BtnHeightB = 32.0
	BtnWidthB  = 128.0
	BtnCfgB    = btn.And(
		BtnCfgC,
		btn.Binding(mouse.Start, func(id int, _ interface{}) int {
			b := event.GetEntity(id).(btn.Btn)
			if _, ok := b.Metadata("inactive"); ok {
				return 0
			}
			r := b.GetRenderable()
			m := r.(render.Modifiable)
			m.Filter(mod.Brighten(25))
			return 0
		}),
		btn.Binding(mouse.Stop, func(id int, _ interface{}) int {
			b := event.GetEntity(id).(btn.Btn)
			if _, ok := b.Metadata("inactive"); ok {
				return 0
			}
			err := btn.Revert(b, 1)
			dlog.ErrorCheck(err)
			return 0
		}),
	)
	BtnCfgC = btn.And(
		btn.Width(BtnWidthB),
		btn.Height(BtnHeightB),
		btn.Mod(mod.And(
			mod.CutRound(.05, .25),
			mods.Inset(func(c color.Color) color.Color {
				return mods.Darker(c, .25)
			}, intgeom.UpLeft),
			mods.Highlight(color.RGBA{170, 170, 170, 200}, 1),
			mods.HighlightOff(color.RGBA{0, 0, 0, 100}, 1, 2, 1),
		)),
		btn.TxtOff(10, 10),
	)
)
