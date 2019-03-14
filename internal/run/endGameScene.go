package run

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"

	"github.com/oakmound/weekly87/internal/records"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/menus"
)

// EndScene is a scene in the same package as run to allow for easy variable access.
//If there is time at the end we can look at what vbariables this touches and get them exported or passed onwards so this can have its own package
var EndScene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInGame = true
		nextscene = "load"
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)

		// TODO: This should be in a more central place.
		// Allows us to have text that shows up on a white background
		fnt := render.DefFontGenerator.Copy()
		fnt.Color = render.FontColor("Red")
		fnt.Size = 16
		redFnt := fnt.Generate()
		fnt.Color = render.FontColor("Blue")
		fnt.Size = 14
		blueFnt := fnt.Generate()

		menuBackground, _ := render.LoadSprite("", filepath.Join("raw", "standard_placeholder.png"))
		menuBackground.Modify(mod.FlipX)
		render.Draw(menuBackground, 0)

		textBackingX := oak.ScreenWidth / 3

		textBacking := render.NewColorBox(textBackingX, oak.ScreenHeight*2/3, color.RGBA{120, 120, 120, 190})
		textBacking.SetPos(float64(oak.ScreenWidth)*0.33, 40)
		render.Draw(textBacking, 1)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) * 3 / 4

		btn.New(menus.BtnCfgB, btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3), btn.Pos(menuX, menuY), btn.Text("Return To Menu"), btn.Binding(func(int, interface{}) int {
			nextscene = "startup"
			stayInGame = false
			return 0
		}))

		textY := 60.0

		text := redFnt.NewStrText("Your Ending Statistics", float64(oak.ScreenWidth)/2-80, textY)
		textY += 40
		render.Draw(text, 2, 2)

		sectionText := blueFnt.NewStrText("Sections Cleared: "+strconv.Itoa(runInfo.SectionsCleared), float64(oak.ScreenWidth)/2-80, textY)
		textY += 40
		render.Draw(sectionText, 2, 2)

		enemy := blueFnt.NewStrText("Enemies Defeated: "+strconv.Itoa(runInfo.EnemiesDefeated), float64(oak.ScreenWidth)/2-80, textY)
		textY += 40
		render.Draw(enemy, 2, 2)

		fmt.Println("Cleared out", runInfo.SectionsCleared)
		chestTotal := 0
		for i := 0; i < len(runInfo.Party); i++ {
			// playerJson, _ := json.Marshal(runInfo.Party[i].ChestValues)

			playerChestValue := 0
			for _, j := range runInfo.Party[i].ChestValues {
				playerChestValue += int(j)
			}
			fmt.Println("Chesty Value :", playerChestValue)
			charText := blueFnt.NewStrText("Chests Acquired by Player:"+strconv.Itoa(playerChestValue), float64(textBackingX), textY)
			render.Draw(charText, 2, 3)
			textY += 20

			chestTotal += playerChestValue

			chestMin := oak.ScreenWidth/2 - textBackingX/2
			textX := chestMin

			for _, j := range runInfo.Party[i].ChestValues {
				ch := doodads.NewChest(j)

				xInc, _ := ch.R.GetDims()

				textX += xInc + 10
				if textX > chestMin+textBackingX {
					textX = chestMin
					textY += 40
				}

				ch.SetPos(float64(textX), float64(textY))
				render.Draw(ch.GetRenderable(), 2, 1)
			}
		}

		// partyJson, _ := json.Marshal(runInfo.Party)

		r := records.Load()
		sc := int64(runInfo.SectionsCleared)
		r.SectionsCleared += sc
		if sc > r.FarthestGoneInSections {
			r.FarthestGoneInSections = sc
		}
		// For the next run
		BaseSeed += int64(runInfo.SectionsCleared) + 1
		r.Store()
	},
	Loop: scene.BooleanLoop(&stayInGame),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&nextscene),
}
