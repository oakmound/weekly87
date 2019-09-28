package credits

import (
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/mouse"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/menus/selector"
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

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 6
		menuY := float64(oak.ScreenHeight) / 2.7

		fnt := render.DefFontGenerator.Copy()
		fnt.Color = render.FontColor("Black")
		fnt.Size = 60
		titleFnt := fnt.Generate()

		title := titleFnt.NewStrText("Credits", menuX-20, menuY-40)
		render.Draw(title, 2, 12)

		construction := titleFnt.NewStrText("!Under Construction!", float64(oak.ScreenWidth)/3, menuY+70)
		render.Draw(construction, 2, 12)

		creditFenrir := btn.New(
			menus.BtnCfgC,
			btn.Text("Art - LightningFenrir"),
			btn.Pos(menuX, menuY),
			btn.Color(color.RGBA{0, 125, 255, 255}),
		)

		menuY += menus.BtnHeightB * 1.5

		creditPlaus := btn.New(
			menus.BtnCfgC,
			btn.Text("Code - PlausiblyFun"),
			btn.Pos(menuX, menuY),
			btn.Color(color.RGBA{0, 125, 255, 255}),
		)

		menuY += menus.BtnHeightB * 1.5

		credit200sc := btn.New(
			menus.BtnCfgC,
			btn.Text("Code/Music - 200sc"),
			btn.Pos(menuX, menuY),
			btn.Color(color.RGBA{0, 125, 255, 255}),
		)

		menuY += menus.BtnHeightB * 1.5

		returnBtn := btn.New(menus.BtnCfgB,
			btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
			btn.Pos(menuX, menuY),
			btn.Text("Return To Menu"),
			btn.Binding(mouse.ClickOn,
				func(int, interface{}) int {
					nextscene = "startup"
					stayInMenu = false
					return 0
				}))

		spcs := []*collision.Space{}
		btnList := []btn.Btn{returnBtn, credit200sc, creditPlaus, creditFenrir}
		for _, b := range btnList {
			spcs = append(spcs, b.GetSpace())
		}
		selector.New(
			menus.ButtonSelectorSpacesA(spcs, btnList),
			selector.MouseBindings(true),
		)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End:  scene.GoToPtr(&nextscene),
}
