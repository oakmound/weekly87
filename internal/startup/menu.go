package startup

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

// Scene  to display
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInMenu = true
		nextscene = "startup"
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)

		if prevScene == "" {
			// The game has just started, maybe do some
			// intro stuff
		}
		// Render menu buttons
		// 1. Start game
		// 2. Select save file? <- don't worry about saving progress for first build
		// 3. Settings
		// 4. Credits
		// 5. Exit game

		btnHeight := 30.0
		btnWidth := 120.0

		menuX := (float64(oak.ScreenWidth) - btnWidth) / 2
		menuY := float64(oak.ScreenHeight) / 4

		btnCfg := btn.And(
			btn.Width(btnWidth),
			btn.Height(btnHeight),
			btn.Mod(mods.HighlightOff(colornames.Blue, 3, 0, 0)),
			btn.Mod(mods.InnerHighlightOff(colornames.Black, 1, 0, 0)),
			btn.TxtOff(btnWidth/4, btnHeight/3), //magic numbers
		)

		start := btn.New(btnCfg, btn.Color(colornames.Green), btn.Pos(menuX, menuY), btn.Text("Start Game"))
		menuY += btnHeight * 1.5
		load := btn.New(btnCfg, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Load Game"))
		menuY += btnHeight * 1.5
		settings := btn.New(btnCfg, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Settings"),
			btn.Binding(func(int, interface{}) int {
				nextscene = "settings"
				stayInMenu = false
				return 0
			}))
		menuY += btnHeight * 1.5
		credits := btn.New(btnCfg, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Credits"))
		menuY += btnHeight * 1.5
		exit := btn.New(btnCfg, btn.Pos(menuX, menuY), btn.Text("Exit Game"))
		// render.Draw()

		fmt.Println("How high are the buttons", start.Y(), load.Y(), settings.Y(), credits.Y(), exit.Y())

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End:  scene.GoTo(nextscene),
}
