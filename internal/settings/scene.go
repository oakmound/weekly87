package settings

import (
	"fmt"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/entities/x/mods"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"golang.org/x/image/colornames"
)

var stayInMenu bool
var nextscene string

// SettingsScene  to display our settings
var SettingsScene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInMenu = true
		nextscene = "settings"
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)

		btnHeight := 30.0
		btnWidth := 120.0

		menuX := (float64(oak.ScreenWidth) - btnWidth) / 2
		menuY := float64(oak.ScreenHeight) / 4

		btnCfg := btn.And(
			btn.Width(btnWidth),
			btn.Height(btnHeight),
			btn.Mod(mods.HighlightOff(colornames.Blue, 3, 0, 0)),
			btn.Mod(mods.InnerHighlightOff(colornames.Black, 1, 0, 0)),
			btn.TxtOff(btnWidth/4, btnHeight/3), //magic numbers from main menu
		)

		exit := btn.New(btnCfg, btn.Pos(menuX, menuY), btn.Text("Return To Menu"), btn.Binding(func(int, interface{}) int {
			nextscene = "startup"
			stayInMenu = false
			return 0
		}))

		fmt.Println("How high are the buttons", exit.Y())

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	// scene.GoTo("inn"),
	End: scene.GoTo(nextscene),
}
