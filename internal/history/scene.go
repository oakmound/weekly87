package history

import (
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

		dlog.Verb("Entering the History Scene")
		stayInMenu = true
		nextscene = "history"
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)

		menuBackground, _ := render.LoadSprite("", filepath.Join("raw", "standard_placeholder.png"))
		render.Draw(menuBackground, 0)

		textBackingX := oak.ScreenWidth / 3

		textBacking := render.NewColorBox(textBackingX, oak.ScreenHeight*3/5, color.RGBA{120, 120, 120, 210})
		textBacking.SetPos(float64(oak.ScreenWidth)/18, 80)
		render.Draw(textBacking, 1)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 6
		menuY := float64(textBacking.Bounds().Max.Y) + 84

		r := records.Load()
		dlog.Verb("Records loaded:", r)
		textY := 120.0
		textX := float64(oak.ScreenWidth) / 5

		historyTitle := titleFnt.NewStrText("Your Past Game Stats!", textX, textY)
		historyTitle.Center()
		render.Draw(historyTitle, 2, 2)
		textY += 40

		cleared := strconv.FormatInt(r.SectionsCleared, 10)
		sectionText := blueFnt.NewStrText("Total Sections Cleared: "+cleared, textX, textY)
		sectionText.Center()
		render.Draw(sectionText, 2, 2)
		textY += 40

		farthest := strconv.FormatInt(r.FarthestGoneInSections, 10)
		farthestText := blueFnt.NewStrText("Farthest Section Reached: "+farthest, textX, textY)
		farthestText.Center()
		render.Draw(farthestText, 2, 2)
		textY += 40

		newSavePressed := 0
		newSaveStr := "Are you sure"

		var nStartBtn btn.Btn
		nStartBtn = btn.New(menus.BtnCfgB,
			btn.Color(menus.Purple),
			btn.Pos(menuX, menuY),
			btn.Text("New Save File"), btn.Binding(mouse.ClickOn, func(int, interface{}) int {
				if newSavePressed < 3 {
					type SetStringer interface {
						SetString(string)
					}
					newSaveStr += "?"
					nStartBtn.(SetStringer).SetString(newSaveStr)
					newSavePressed++
					return 0
				}
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
		menuY += 40
		returnBtn := btn.New(menus.BtnCfgB,
			btn.Color(menus.Red),
			btn.Pos(menuX, menuY),
			btn.Text("Return To Menu"), btn.Binding(mouse.ClickOn, func(int, interface{}) int {
				nextscene = "startup"
				stayInMenu = false
				return 0
			}))

		spcs := []*collision.Space{}
		btnList := []btn.Btn{nStartBtn, returnBtn}
		for _, b := range btnList {
			spcs = append(spcs, b.GetSpace())
		}
		selector.New(
			menus.ButtonSelectorSpacesNoWrap(spcs, btnList),
			selector.MouseBindings(true),
		)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&nextscene),
}
