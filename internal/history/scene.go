package history

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/oakmound/oak/mouse"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/records"
)

var stayInMenu bool
var nextscene string

// Scene to display our settings
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		fnt := render.DefFontGenerator.Copy()
		fnt.Color = render.FontColor("Blue")
		fnt.Size = 14
		blueFnt := fnt.Generate()

		fmt.Println("Starting history scene")
		stayInMenu = true
		nextscene = "history"
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)
		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) * 3 / 4

		menuBackground, _ := render.LoadSprite("", filepath.Join("raw", "standard_placeholder.png"))
		render.Draw(menuBackground, 0)

		r := records.Load()
		fmt.Println(r)
		textY := 60.0

		cleared := strconv.FormatInt(r.SectionsCleared, 10)

		sectionText := blueFnt.NewStrText("Sections Cleared: "+cleared, float64(oak.ScreenWidth)/2-80, textY)
		textY += 40
		render.Draw(sectionText, 2, 2)

		btn.New(menus.BtnCfgA,
			btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
			btn.Pos(menuX, menuY),
			btn.Text("Return To Menu"), btn.Binding(mouse.ClickOn, func(int, interface{}) int {
				nextscene = "startup"
				stayInMenu = false
				return 0
			}))

		text := render.DefFont().NewStrText("Save Management is under construction", float64(oak.ScreenWidth)/2-100, float64(oak.ScreenHeight)/4)
		render.Draw(text, 0, 1)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&nextscene),
}
