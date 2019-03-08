package startup

import (
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/entities/x/mods"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"golang.org/x/image/colornames"
)

var stayInMenu bool

var Scene = scene.Scene{
	func(prevScene string, data interface{}) {

		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
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
		btnCfg := btn.And(
			btn.Color(colornames.Blueviolet),
			btn.Width(120),
			btn.Height(30),
			btn.Mod(mods.HighlightOff(colornames.Blue, 3, 0, 0)),
			btn.Mod(mods.InnerHighlightOff(colornames.Black, 1, 0, 0)),
		)
		start := btn.New(btnCfg, btn.Text("Start Game"))
		render.Draw()

	},
	scene.BooleanLoop(&stayInMenu),
	scene.GoTo("inn"),
}
