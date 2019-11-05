package end

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/physics"

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
	"github.com/oakmound/weekly87/internal/menus/selector"
	"github.com/oakmound/weekly87/internal/records"
	"github.com/oakmound/weekly87/internal/run"
	"github.com/oakmound/weekly87/internal/sfx"
)

// Init to be called after oak start up to get our asset reference
func Init() {
	_, err := render.LoadSprite("", filepath.Join("raw", "wood_junk.png"))
	dlog.Error("Something went wrong loading in wood junk " + err.Error())
}

var (
	graveX   = 90.0
	graveY   = 300.0
	npcScale = 1.6 //TODO: remove need for this and rename
)

var stayInEndScene bool
var endSceneNextScene string

type intStringer struct {
	i *int
}

func (is intStringer) String() string {
	return strconv.Itoa(*is.i)
}

var (
	presentationX, presentationY, startY, pitX, pitY, hopDistance float64
)

// Scene is a scene in the same package as run to allow for easy variable access.
//If there is time at the e nd we can look at what vbariables this touches and get them exported or passed onwards so this can have its own package
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		r := records.Load()
		justVisiting := false

		outcome, ok := data.(run.Outcome)
		if !ok {
			justVisiting = true
			outcome = run.Outcome{R: r.LastRun}
		}
		dlog.Info("Just visting end game:", justVisiting)

		// STANDARD SETUP STUFF
		runInfo := outcome.R
		stayInEndScene = true
		endSceneNextScene = "inn"

		render.SetDrawStack(layer.Get()...)

		// Given that oak.Screen variables change in game we dont farm this out to an init
		presentationX = float64(oak.ScreenWidth / 2)
		presentationY = float64(oak.ScreenHeight/2) + 20
		startY = float64(oak.ScreenHeight/2) + 20

		pitX = float64(oak.ScreenWidth) - 180.0
		pitY = float64(oak.ScreenHeight) - 75.0

		hopDistance = 100.0

		// Update info in the records with info from the most recent run
		// TODO: debate moving this to the run scene's end function

		sc := int64(runInfo.SectionsCleared)
		r.SectionsCleared += sc
		if sc > r.FarthestGoneInSections {
			r.FarthestGoneInSections = sc
		}

		// Display variables
		currentDeathToll := r.Deaths

		chestTotal := 0
		deadChests := 0
		if !justVisiting {
			for _, pl := range runInfo.Party.Players {
				playerChestValue := 0
				for _, j := range pl.ChestValues {
					playerChestValue += int(j)
				}
				if !pl.Alive {
					r.Deaths++
					deadChests += playerChestValue
				} else {
					chestTotal += playerChestValue
				}
			}
		}

		// For the next run TODO: move to run
		r.BaseSeed = int64(runInfo.SectionsCleared) + 1

		r.Wealth += chestTotal
		r.EnemiesDefeated += runInfo.EnemiesDefeated

		r.Store()

		fnt := render.DefFontGenerator.Copy()

		fnt.Color = render.FontColor("Blue")
		fnt.Size = 14
		blueFnt := fnt.Generate()
		fnt.Color = render.FontColor("Black")
		fnt.Size = 20
		graves := fnt.Generate()

		render.SetDrawStack(layer.Get()...)
		debugTree := dtools.NewThickColoredRTree(collision.DefTree, 4, labels.ColorMap)
		render.Draw(debugTree, layer.Play, 1000)

		// Make the graveyard backing
		endBackground, _ := render.LoadSprite("", filepath.Join("raw", "end_scene.png"))
		render.Draw(endBackground, layer.Ground)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := 16.0

		if !justVisiting {
			b := btn.New(menus.BtnCfgB,
				btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
				btn.Pos(menuX, menuY), btn.Text("Return To Inn"),
				btn.Binding(mouse.ClickOn, func(int, interface{}) int {
					stayInEndScene = false
					return 0
				}))
			selector.New(
				selector.Layers(layer.UI, 3),
				selector.HorzArrowControl(),
				selector.JoystickHorzDpadControl(),
				selector.SelectTrigger(key.Down+key.Spacebar),
				selector.SelectTrigger("A"+joystick.ButtonUp),

				selector.Spaces(b.GetSpace()),

				selector.InteractTrigger(key.Down+key.B, "boot"),
				selector.InteractTrigger("X"+joystick.ButtonUp, "boot"),

				selector.DestroyTrigger(key.Down+key.Escape),
				selector.DestroyTrigger("B"+joystick.ButtonUp),
				selector.MouseBindings(true),
				selector.Callback(func(i int, _ ...interface{}) {
					sfx.Play("selected")
					b.Trigger(mouse.ClickOn, nil)
				}),
				selector.Display(func(pt floatgeom.Point2) render.Renderable {
					poly, err := render.NewPolygon(
						floatgeom.Point2{0, 0},
						floatgeom.Point2{pt.X(), 0},
						floatgeom.Point2{pt.X(), pt.Y()},
						floatgeom.Point2{0, pt.Y()},
					)
					dlog.ErrorCheck(err)
					return poly.GetThickOutline(menus.Gold, 2)
				}),
			)

		}

		goldPit := floatgeom.NewRect2WH(670, float64(oak.ScreenHeight)-100, 330, 100)
		makeGoldParticles(r.Wealth, goldPit)

		textY := 40.0
		textX := float64(oak.ScreenWidth) / 6

		currentDeathTollp := &currentDeathToll
		render.Draw(graves.NewStrText("Current Total Deaths:", textX, textY), 1, 10)
		render.Draw(graves.NewText(intStringer{currentDeathTollp}, textX+200, textY), 1, 10)

		textY += 40

		titling := blueFnt.NewStrText("Last Run Info:", textX, textY)
		textX += 120
		render.Draw(titling, 2, 2)

		sectionText := blueFnt.NewStrText("Sections Cleared: "+strconv.Itoa(runInfo.SectionsCleared), textX, textY)
		textX += 200
		render.Draw(sectionText, 2, 2)

		enemy := blueFnt.NewStrText("Enemies Defeated: "+strconv.FormatInt(runInfo.EnemiesDefeated, 10), textX, textY)
		textX += 200
		render.Draw(enemy, 2, 2)

		chestValues := blueFnt.NewStrText("Chest Value: "+strconv.Itoa(chestTotal), textX, textY)
		textX += 200
		render.Draw(chestValues, 2, 2)

		// Current Run info
		textX = float64(oak.ScreenWidth) / 6
		textY += 30

		titling = blueFnt.NewStrText("Overall Stats:", textX, textY)
		textX += 120
		render.Draw(titling, 2, 2)

		sectionText = blueFnt.NewStrText("Sections Cleared: "+strconv.Itoa(int(r.SectionsCleared)), textX, textY)
		textX += 200
		render.Draw(sectionText, 2, 2)

		enemy = blueFnt.NewStrText("Enemies Defeated: "+strconv.FormatInt(r.EnemiesDefeated, 10), textX, textY)
		textX += 200
		render.Draw(enemy, 2, 2)

		chestValues = blueFnt.NewStrText("Chest Value: "+strconv.Itoa(r.Wealth), textX, textY)
		textX += 200
		render.Draw(chestValues, 2, 2)

		debugElements := []render.Renderable{}
		// debug locations
		debugElements = append(debugElements, render.NewColorBox(5, 5, color.RGBA{200, 200, 10, 255}))
		debugElements[0].SetPos(presentationX, presentationY)

		debugElements = append(debugElements, render.NewColorBox(5, 5, color.RGBA{200, 200, 200, 255}))
		debugElements[1].SetPos(graveX, presentationY)

		debugElements = append(debugElements, render.NewColorBox(5, 5, color.RGBA{10, 200, 200, 255}))
		debugElements[2].SetPos(pitX, pitY)
		for _, r := range debugElements {
			render.Draw(r, layer.Debug, 18)
		}

		// A way to return to the inn scene
		dWidth := 150.0
		doodads.NewCustomInnDoor("inn", (float64(oak.ScreenWidth)-dWidth)/2, float64(oak.ScreenHeight)-20, dWidth, 20)
		// Block off the top of the inn from being walkable
		doodads.NewFurniture(0, 0, float64(oak.ScreenWidth), 187) // top of inn

		addDebugCommands(debugTree, debugElements)
		if justVisiting == true {
			visitEnter(r.PartyComp)
			return
		}

		// Extra running info
		presentSpoils(runInfo.Party, currentDeathTollp, 0)

	},
	Loop: scene.BooleanLoop(&stayInEndScene),

	End: scene.GoToPtr(&endSceneNextScene),
}

// investigate allows us to investigate and poke around the end game with our living characters
func investigate(party *players.Party) {
	fmt.Println("Ending actions can be taken here")

}

// visitEnter frame binding
func visitEnter(pComp []players.PartyMember) {
	ptycon := players.PartyConstructor{
		Players:    players.ClassConstructor(pComp),
		MaxPlayers: len(pComp),
	}

	pty, err := ptycon.NewParty(true)
	if err != nil {
		dlog.Error(err)
		return
	}
	pc := newEndWalker(npcScale, pty.Players)

	stop1 := float64(oak.ScreenHeight) - 140

	stop2 := float64(oak.ScreenWidth/2) - 200

	pc.Front.Delta = physics.NewVector(0, -8)
	pc.Front.Bind(func(id int, _ interface{}) int {
		p, ok := event.GetEntity(id).(*entities.Interactive)
		if !ok {
			dlog.Error("Non-player sent to player binding")
		}
		_, y := p.GetPos()
		if y < stop1 {
			pc.State = overridable

			p.RSpace.Add(collision.Label(labels.Door), (func(s1, s2 *collision.Space) {
				_, ok := s2.CID.E().(*doodads.InnDoor)
				if !ok {
					dlog.Error("Non-door sent to inndoor binding")
					return
				}
				stayInEndScene = false
			}))

			pc.Front.Delta = physics.NewVector(-4, 0)
			pc.Front.Bind(func(id int, _ interface{}) int {
				if pc.State == playing {
					investigate(pty)
					return event.UnbindSingle
				}
				p, ok := event.GetEntity(id).(*entities.Interactive)
				if !ok {
					dlog.Error("Non-player sent to player binding")
				}
				x, _ := p.GetPos()

				if x < stop2 {

					pc.State = playing
					investigate(pty)
					return event.UnbindSingle

				}

				return 0
			}, "EnterFrame")
			return event.UnbindSingle
		}

		return 0
	}, "EnterFrame")

}

func addDebugCommands(debugTree *dtools.Rtree, debugElements []render.Renderable) {
	oak.ResetCommands()
	oak.AddCommand("debug", func(args []string) {
		dlog.Warn("Cheating to toggle debug mode")
		if debugTree.DrawDisabled {
			dlog.Warn("Debug turned off")
			debugTree.DrawDisabled = false
			for _, r := range debugElements {
				render.Draw(r, layer.Debug, 18)
			}
			return
		}
		dlog.Warn("Debug turned on")
		debugTree.DrawDisabled = true
		for _, r := range debugElements {
			r.Undraw()
		}

	})

	dlog.Info("Current Debug Commands are: ", strings.Join(oak.GetDebugKeys(), " , "))
	oak.AddCommand("help", func(args []string) {
		dlog.Info("Current Debug Commands are: ", strings.Join(oak.GetDebugKeys(), " , "))
	})
}
