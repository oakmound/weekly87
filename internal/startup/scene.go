package startup

import (
	"fmt"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/menus"
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

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) / 4

		start := btn.New(menus.BtnCfgA, btn.Color(colornames.Green), btn.Pos(menuX, menuY), btn.Text("Start Game"),
			btn.Binding(func(int, interface{}) int {
				nextscene = "inn"
				stayInMenu = false
				return 0
			}))
		menuY += menus.BtnHeightA * 1.5
		load := btn.New(menus.BtnCfgA, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Load Game"))
		menuY += menus.BtnHeightA * 1.5
		settings := btn.New(menus.BtnCfgA, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Settings"),
			btn.Binding(func(int, interface{}) int {
				nextscene = "settings"
				stayInMenu = false
				return 0
			}))
		menuY += menus.BtnHeightA * 1.5
		credits := btn.New(menus.BtnCfgA, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Credits"))
		menuY += menus.BtnHeightA * 1.5
		exit := btn.New(menus.BtnCfgA, btn.Pos(menuX, menuY), btn.Text("Exit Game"))
		// render.Draw()

		fmt.Println("How high are the buttons", start.Y(), load.Y(), settings.Y(), credits.Y(), exit.Y())

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End:  scene.GoToPtr(&nextscene),
}
