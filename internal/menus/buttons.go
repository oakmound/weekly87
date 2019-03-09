package menus

import "github.com/oakmound/oak/entities/x/btn"
import "github.com/oakmound/oak/entities/x/mods"
import "golang.org/x/image/colornames"

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
)
