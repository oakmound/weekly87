package inn

import (
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/menus"
)

var stayInMenu bool
var nextscene string

// Scene  to display the inn
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInMenu = true
		nextscene = "inn"
		render.SetDrawStack(
			// ground
			render.NewCompositeR(),
			// entities
			render.NewHeap(false),
			// ui
			render.NewHeap(true),
		)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) / 4

		btn.New(menus.BtnCfgA, btn.Layers(2, 0), btn.Pos(menuX, menuY), btn.Text("Start Run"), btn.Binding(func(int, interface{}) int {
			nextscene = "run"
			stayInMenu = false
			return 0
		}))

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "placeholder_inn.png"))
		render.Draw(innBackground, 0)

		// A way to enter the run
		innDoor := characters.NewDoor()
		iW, iH := innDoor.R.GetDims()
		innDoor.SetPos(float64(oak.ScreenWidth-iW), float64(oak.ScreenHeight-iH)/2) //Center the door on the right side
		render.Draw(innDoor.R, 1)

		text := render.DefFont().NewStrText("Hit the button or walk out of the inn to start the game!", float64(oak.ScreenWidth)/2-120, float64(oak.ScreenHeight)/4-40)
		render.Draw(text, 2, 1)

		innSpace := floatgeom.NewRect2(0, 0, float64(oak.ScreenWidth), float64(oak.ScreenHeight)-32) //Adjusted for the current size of the spearman

		NewInnWalker(innSpace)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End:  scene.GoToPtr(&nextscene),
}
