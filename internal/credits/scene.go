package credits

import (
	"image/color"
	"path/filepath"

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
		nextscene = "credits"
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)

		menuBackground, _ := render.LoadSprite("", filepath.Join("raw", "standard_placeholder.png"))
		render.Draw(menuBackground, 0)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) * 1 / 4

		btn.New(
			menus.BtnCfgB,
			btn.Text("Art - LightningFenrir"),
			btn.Pos(menuX, menuY),
			btn.Color(color.RGBA{0, 125, 255, 255}),
		)

		menuY += menus.BtnHeightB * 1.5

		btn.New(
			menus.BtnCfgB,
			btn.Text("Code - PlausiblyFun"),
			btn.Pos(menuX, menuY),
			btn.Color(color.RGBA{0, 125, 255, 255}),
		)

		menuY += menus.BtnHeightB * 1.5

		btn.New(
			menus.BtnCfgB,
			btn.Text("Code/Music - 200sc"),
			btn.Pos(menuX, menuY),
			btn.Color(color.RGBA{0, 125, 255, 255}),
		)

		menuY += menus.BtnHeightB * 1.5

		btn.New(menus.BtnCfgB, btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3), btn.Pos(menuX, menuY), btn.Text("Return To Menu"), btn.Binding(func(int, interface{}) int {
			nextscene = "startup"
			stayInMenu = false
			return 0
		}))

		text := render.DefFont().NewStrText("The Credits are under construction", float64(oak.ScreenWidth)/2-100, float64(oak.ScreenHeight)/4)
		render.Draw(text, 0, 1)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End:  scene.GoToPtr(&nextscene),
}
