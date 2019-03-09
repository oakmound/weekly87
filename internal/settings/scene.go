package settings

import (
	"fmt"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/menus"
)

var stayInMenu bool
var nextscene string

// Scene to display our settings
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInMenu = true
		nextscene = "settings"
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) / 4

		exit := btn.New(menus.BtnCfgA, btn.Pos(menuX, menuY), btn.Text("Return To Menu"), btn.Binding(func(int, interface{}) int {
			nextscene = "startup"
			stayInMenu = false
			return 0
		}))

		fmt.Println("How high are the buttons", exit.Y())

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&nextscene),
}
