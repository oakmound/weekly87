package history

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/mouse"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/menus/selector"
	"github.com/oakmound/weekly87/internal/records"
)

var stayInMenu bool
var nextscene string

// Scene to display our settings
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		fnt := render.DefFontGenerator.Copy()
		fnt.Color = render.FontColor("Blue")
		fnt.Size = 18
		titleFnt := fnt.Generate()

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

		textBackingX := oak.ScreenWidth / 3

		textBacking := render.NewColorBox(textBackingX, oak.ScreenHeight*2/3, color.RGBA{120, 120, 120, 210})
		textBacking.SetPos(float64(oak.ScreenWidth)*0.33, 40)
		render.Draw(textBacking, 1)

		r := records.Load()
		dlog.Info("Records loaded:", r)
		textY := 60.0

		historyTitle := titleFnt.NewStrText("Your Past Game Stats!", float64(oak.ScreenWidth)/2, textY)
		historyTitle.Center()
		render.Draw(historyTitle, 2, 2)
		textY += 40

		cleared := strconv.FormatInt(r.SectionsCleared, 10)
		sectionText := blueFnt.NewStrText("Total Sections Cleared: "+cleared, float64(oak.ScreenWidth)/2, textY)
		sectionText.Center()
		render.Draw(sectionText, 2, 2)
		textY += 40

		farthest := strconv.FormatInt(r.FarthestGoneInSections, 10)
		farthestText := blueFnt.NewStrText("Farthest Section Reached: "+farthest, float64(oak.ScreenWidth)/2, textY)
		farthestText.Center()
		render.Draw(farthestText, 2, 2)
		textY += 40

		nStartBtn := btn.New(menus.BtnCfgA,
			btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
			btn.Pos(menuX, textY),
			btn.Text("New Save File"), btn.Binding(mouse.ClickOn, func(int, interface{}) int {
				fmt.Println("HAI THERE")
				newPath, err := records.Archive()
				nextscene = "history"
				stayInMenu = false
				if err != nil {
					dlog.Error("Failed to move save file to", newPath, "due to", err)
					return 0
				}
				dlog.Info("Moved the current save file to ", newPath)

				return 0
			}))

		returnBtn := btn.New(menus.BtnCfgA,
			btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
			btn.Pos(menuX, menuY),
			btn.Text("Return To Menu"), btn.Binding(mouse.ClickOn, func(int, interface{}) int {
				nextscene = "startup"
				stayInMenu = false
				return 0
			}))

		spcs := []*collision.Space{}
		btnList := []btn.Btn{returnBtn, nStartBtn}
		for _, b := range btnList {
			spcs = append(spcs, b.GetSpace())
		}
		selector.New(
			menus.ButtonSelectorSpacesA(spcs, btnList),
			selector.MouseBindings(true),
		)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&nextscene),
}
