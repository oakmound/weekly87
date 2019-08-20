package end

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/dtools"
	"github.com/oakmound/weekly87/internal/layer"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/records"
	"github.com/oakmound/weekly87/internal/run"
)

var stayInEndScene bool
var endSceneNextScene string

// Scene is a scene in the same package as run to allow for easy variable access.
//If there is time at the e nd we can look at what vbariables this touches and get them exported or passed onwards so this can have its own package
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {

		outcome := data.(run.Outcome)
		fmt.Println(outcome.R)
		runInfo := outcome.R
		stayInEndScene = true
		endSceneNextScene = "inn"

		render.SetDrawStack(layer.Get()...)

		fnt := render.DefFontGenerator.Copy()

		fnt.Color = render.FontColor("Blue")
		fnt.Size = 14
		blueFnt := fnt.Generate()

		// nextscene = "inn"
		render.SetDrawStack(layer.Get()...)
		debugTree := dtools.NewThickColoredRTree(collision.DefTree, 4, labels.ColorMap)
		render.Draw(debugTree, layer.Play, 1000)

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "end_scene.png"))
		render.Draw(innBackground, layer.Ground)

		// Text overlay info
		textBackingX := oak.ScreenWidth / 3

		// textBacking := render.NewColorBox(textBackingX, oak.ScreenHeight*2/3, color.RGBA{120, 120, 120, 190})
		// textBacking.SetPos(float64(oak.ScreenWidth)*0.33, 40)
		// render.Draw(textBacking, 1)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := 32.0

		btn.New(menus.BtnCfgB,
			btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
			btn.Pos(menuX, menuY), btn.Text("Return To Inn"),
			btn.Binding(mouse.ClickOn, func(int, interface{}) int {
				stayInEndScene = false
				return 0
			}))

		textY := 60.0

		//TODO: 2 layers: first current run
		// Second layer totals
		// Each layer has: sections_completed, enemies_defeated, chestvalue

		sectionText := blueFnt.NewStrText("Sections Cleared: "+strconv.Itoa(runInfo.SectionsCleared), float64(oak.ScreenWidth)/2-80, textY)
		textY += 40
		render.Draw(sectionText, 2, 2)

		enemy := blueFnt.NewStrText("Enemies Defeated: "+strconv.Itoa(runInfo.EnemiesDefeated), float64(oak.ScreenWidth)/2-80, textY)
		textY += 40
		render.Draw(enemy, 2, 2)

		fmt.Println("Cleared out", runInfo.SectionsCleared)
		chestTotal := 0
		for _, pl := range runInfo.Party.Players {
			// playerJson, _ := json.Marshal(runInfo.Party[i].ChestValues)

			playerChestValue := 0
			for _, j := range pl.ChestValues {
				playerChestValue += int(j)
			}
			chestTotal += playerChestValue
			fmt.Println("Chest Value :", playerChestValue)
			_ = blueFnt.NewStrText("Chests Acquired by Player:"+strconv.Itoa(playerChestValue), float64(textBackingX), textY)
			// render.Draw(charText, 2, 3)
			textY += 16

			chestMin := oak.ScreenWidth/2 - textBackingX/2
			textX := chestMin
			for _, j := range pl.ChestValues {
				ch := doodads.NewChest(j)

				xInc, _ := ch.R.GetDims()

				textX += xInc + 10
				if textX > chestMin+textBackingX {
					textX = chestMin
					textY += 15
				}

				ch.SetPos(float64(textX), float64(textY))
				render.Draw(ch.GetRenderable(), 2, 1)
			}
			if len(pl.ChestValues) > 0 {
				textY += 30
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
		// BaseSeed := int64(runInfo.SectionsCleared) + 1
		r.Store()

		// Extra running info
		presentSpoils(runInfo.Party, 0)

	},
	Loop: scene.BooleanLoop(&stayInEndScene),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&endSceneNextScene),
}

func presentSpoils(party *players.Party, index int) {
	if index > len(party.Players)-1 {
		return
	}
	p := party.Players[index]
	p.CID = p.Init()
	//Player enters stage right
	p.SetPos(float64(oak.ScreenWidth-64), float64(oak.ScreenHeight/2))
	render.Draw(p.R, layer.Play, 20)
	fmt.Printf("\nCharacter %d walking through as %s ", index, p.R)

	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		ply.R.ShiftX(-2)
		// TODO: the following
		// Player comes to middle ish

		// Do something with spoils or look sad if dead
		// If alive: throw chests with hop and cheer
		//		Card explaining the amount in their chests?
		//		Chest explodes into money and goes into pit
		// If dead: whomp whomp
		// 		eulogoy? name, class, run?

		// Next person in party starts process
		// You walk to end point (graves or bottom)
		// When reach your end point destroy self

		if ply.R.X() < float64(oak.ScreenWidth/2) {

			ply.R.Undraw()
			presentSpoils(party, index+1)
			return event.UnbindSingle
		}
		return 0
	}, "EnterFrame")

}
