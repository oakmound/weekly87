package inn

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"path/filepath"
	"strings"
	"sync"
	"time"

	klg "github.com/200sc/klangsynthese/audio"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/dtools"
	"github.com/oakmound/weekly87/internal/keyviz"
	"github.com/oakmound/weekly87/internal/layer"
	"github.com/oakmound/weekly87/internal/menus/selector"
	"github.com/oakmound/weekly87/internal/music"
	"github.com/oakmound/weekly87/internal/records"
)

var stayInMenu bool
var nextscene string
var bkgMusic *klg.Audio
var curRecord *records.Records

// Scene  to display the inn
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {

		// event.AddAliases(map[string][]string
		// 	"EndPartySelect": []string{
		// 		key.Down+key.Escape,
		// 		"B"+joystick.ButtonUp,
		// 	},
		// )

		stayInMenu = true
		nextscene = "inn"
		render.SetDrawStack(layer.Get()...)
		debugTree := dtools.NewThickColoredRTree(collision.DefTree, 4, labels.ColorMap)
		render.Draw(debugTree, layer.Play, 1000)

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "placeholder_inn.png"))
		render.Draw(innBackground, layer.Ground)

		// A way to enter the run
		doodads.NewInnDoor("run")
		// A way to go back to menu screen
		doodads.NewCustomInnDoor("endGame", 509, 40, 100, 102)

		// Create doodads for tables
		uglymugger, _ := render.LoadSprites("", filepath.Join("16x16", "ugly_mugger.png"), 16, 16, 0)
		prettyMugs := []*render.Sprite{uglymugger[0][0], uglymugger[0][1], uglymugger[1][0]}

		// Block off the top of the inn from being walkable
		doodads.NewFurniture(0, 0, float64(oak.ScreenWidth), 140) // top of inn

		// Additional inn such as tables
		doodads.NewFurniture(130, 130, 100, float64(oak.ScreenHeight)-130)
		// leftT.PlaceConsumablesa(prettyMugs, 3+rand.Intn(7))

		// TODO: update placement strats for consumables
		doodads.NewConsumable(150, 150, prettyMugs[0])

		topT := doodads.NewFurniture(480, 230, 190, 60)
		topT.SetOrnaments(prettyMugs, rand.Intn(6))

		botT := doodads.NewFurniture(480, 430, 180, 55)
		botT.SetOrnaments(prettyMugs, rand.Intn(6))

		// Create the notes on the notice board
		noteSpace := floatgeom.NewRect2WH(240, 60, 105, 65)
		noteHeight := 3
		for i := 0; i < 8+rand.Intn(5)*3; i++ {
			noteHeight = doodads.NewNote(noteSpace, noteHeight)
		}

		// Create all possible NPCs
		npcScale := 1.6
		npcs := []*NPC{
			NewInnNPC(players.Swordsman, npcScale, 243, 400).FaceLeft(true),
			NewInnNPC(players.Mage, npcScale, 440, 210),
			NewInnNPC(players.WhiteMage, npcScale, 670, 423).FaceLeft(true),
			NewInnkeeper(npcScale, 90, 200),
			NewInnNPC(players.Berserker, npcScale, 445, 460),
			NewInnNPC(players.BlueMage, npcScale, 241, 210).FaceLeft(true),
			NewInnNPC(players.Paladin, npcScale, 240, 280).FaceLeft(true),
			// 	NewInnNPC(players.Spearman, npcScale, 675, 477).FaceLeft(true),
			// 	NewInnNPC(players.TimeMage, npcScale, 680, 230).FaceLeft(true),
		}

		// Start: Swordsman, size 1 party
		// 3 Sections: Mage
		// 10 Sections: Two person party, white mage
		// 25 Sections: Berserker
		// 45 Sections: Three person party
		// 75 Sections: BlueMage
		// 120 Sections: Paladin
		// 200 Sections: Four person party

		// Future: More modes
		// Custom: Choose your own abilities, model, color
		// Chaos: All abilities, models, colors are random (no duplicate abilities for one char)

		// Inn does quite a few operations on our record (mainly for party purposes)
		curRecord = records.Load()

		charUnlocks := []int{
			0,
			3,
			10,
			15,
			25,
			75,
			120,
		}
		progress := len(charUnlocks)
		for i, cu := range charUnlocks {
			if int64(cu) > curRecord.SectionsCleared {
				progress = i
				break
			}
		}
		futureNpcs := npcs[progress:len(npcs)]
		npcs = npcs[0:progress]
		for _, np := range npcs {
			np.Activate()
		}
		for _, fn := range futureNpcs {
			fn.Destroy() // Being extra safe
		}

		partySizeUnlocks := []int{
			0,
			10,
			45,
			200,
		}
		partySize := len(partySizeUnlocks)
		for i, psu := range partySizeUnlocks {
			if int64(psu) > curRecord.SectionsCleared {
				partySize = i
				break
			}
		}

		// Create the player and the display of the party at the top of the screen
		ptycon := players.PartyConstructor{
			Players:    players.ClassConstructor(curRecord.PartyComp),
			MaxPlayers: partySize,
		}

		partyBackground, _ := render.LoadSprite("", filepath.Join("raw", "selector_background.png"))

		partyBackground.SetPos(30, 20)

		ptyOffset := floatgeom.Point2{players.WallOffset, 30}
		ptycon.Players[0].Position = ptyOffset
		pty, err := ptycon.NewParty(true)
		if err != nil {
			dlog.Error(err)
			return
		}
		for _, p := range pty.Players {
			collision.Remove(p.GetSpace())
		}

		pc := newInnWalker(npcScale, pty.Players)

		// Lazy impl for start game walking
		pc.front.Delta = physics.NewVector(4, 0)
		pc.front.Bind(func(id int, _ interface{}) int {
			if pc.gameState == playing {
				return 1
			}
			p, ok := event.GetEntity(id).(*entities.Interactive)
			if !ok {
				dlog.Error("Non-player sent to player binding")
			}
			x, _ := p.GetPos()
			if x > float64(oak.ScreenWidth)/3*2 {
				pc.front.Delta = physics.NewVector(0, 4)
				pc.front.Bind(func(id int, _ interface{}) int {
					if pc.gameState == playing {
						return 1
					}
					p, ok := event.GetEntity(id).(*entities.Interactive)
					if !ok {
						dlog.Error("Non-player sent to player binding")
					}
					x, y := p.GetPos()

					if y > float64(oak.ScreenHeight/2)+38 {
						if x < float64(oak.ScreenWidth)/2 {
							pc.gameState = playing
							return event.UnbindSingle
						}
						pc.front.Delta = physics.NewVector(-4, 0)
					}

					return 0
				}, "EnterFrame")
				return event.UnbindSingle
			}

			return 0
		}, "EnterFrame")

		var lastInteractedNPC *NPC
		// Interact with NPCs
		pc.front.RSpace.Add(labels.NPC, func(_, n *collision.Space) {
			// Todo: what if we're touching multiple NPCS?
			// Keep track of most recently touched npc
			npc, ok := n.CID.E().(*NPC)
			if !ok {
				dlog.Error("Non-npc sent to npc binding")
				return
			}
			if lastInteractedNPC == npc {
				// reset undraw countdown
				npc.UndrawButtonAt = time.Now().Add(1 * time.Second)
				return
			}
			lastInteractedNPC = npc
			fmt.Println(lastInteractedNPC)

			// Make a pop up above the NPC being touched
			// as a command prompt
			interactBtn := getInteractBtn()
			w, h := interactBtn.GetDims()
			npcW, _ := npc.R.GetDims()

			interactBtn.SetPos(npc.X()+float64(npcW-w)/2, npc.Y()-10-float64(h))
			render.Draw(interactBtn, layer.UI, 1)
			npc.Button = interactBtn

			npc.UndrawButtonAt = time.Now().Add(1 * time.Second)

			npc.Bind(func(id int, f interface{}) int {
				frame, ok := f.(int)
				if !ok {
					dlog.Error("Expected int in enterframe")
					return 0
				}
				if frame%15 == 0 {
					npc, ok := event.GetEntity(id).(*NPC)
					if !ok {
						dlog.Error("Non-npc sent to npc binding")
						return 0
					}
					if time.Now().After(npc.UndrawButtonAt) {
						npc.Button.Undraw()
						if lastInteractedNPC == npc {
							lastInteractedNPC = nil
						}
					}
				}
				return 0
			}, "EnterFrame")
		})

		interactLock := &sync.Mutex{}
		partySelectStart := func(int, interface{}) int {

			interactLock.Lock()
			if lastInteractedNPC == nil {
				interactLock.Unlock()
				return 0
			}
			npc := lastInteractedNPC
			lastInteractedNPC = nil
			npc.Button.Undraw()
			interactLock.Unlock()
			npcW, _ := npc.R.GetDims()
			bkgW, bkgH := partyBackground.GetDims()
			// Disable all controls
			pc.gameState = inMenu
			// Spawn a box with the party above the npc being selected
			partyBackground.SetPos(npc.X()+float64(npcW-bkgW)/2, npc.Y()-10-float64(bkgH))
			ptycon.Players[0].Position = floatgeom.Point2{partyBackground.X() + 20, partyBackground.Y() + 10}
			render.Draw(partyBackground, layer.UI, 1)

			fmt.Println("Ptycon players", len(ptycon.Players))
			pty, err := ptycon.NewParty(true)
			dlog.ErrorCheck(err)
			spcs := make([]*collision.Space, 0)
			for _, p := range pty.Players {
				render.Draw(p.R, layer.UI, 2)
				spcs = append(spcs, p.GetSpace())
				if p.Special1 == nil {
					break
				}
			}

			// Show a confirm button, a cancel button and a boot button
			cnfrm := getConfirmBtn()
			cnfrm.SetPos(partyBackground.X(), partyBackground.Y()+float64(bkgH)+2)
			render.Draw(cnfrm, layer.UI, 1)
			boot := getBootButton()
			intrctW, _ := boot.GetDims()
			boot.SetPos(partyBackground.X()+float64(bkgW-intrctW), partyBackground.Y()+float64(bkgH)+2)
			render.Draw(boot, layer.UI, 1)
			cancl := getCancelBtn()
			canclW, canclH := cancl.GetDims()
			cancl.SetPos(partyBackground.X()+float64(bkgW-canclW), partyBackground.Y()-float64(canclH-2))
			render.Draw(cancl, layer.UI, 1)

			var cSelect *selector.Selector

			// Let arrow keys / joystick or mouse even control which party member is selected
			cSelect, _ = selector.New(
				selector.Layers(layer.UI, 3),
				selector.HorzArrowControl(),
				selector.JoystickHorzDpadControl(),
				selector.Spaces(spcs...),
				selector.Callback(func(i int, data ...interface{}) {
					if len(data) == 0 {
						// modify party
						if len(curRecord.PartyComp) <= i {
							curRecord.PartyComp = append(curRecord.PartyComp, players.PartyMember{})
						}
						curRecord.PartyComp[i].PlayerClass = npc.Class
						ptycon.Players = players.ClassConstructor(curRecord.PartyComp)
						return
					}

					dlog.Info("trying to interact  with selector ")
					op, ok := data[0].(string)
					if !ok {
						dlog.Warn("Inn selector recieved a non-string", data[0])
						return
					}
					switch op {
					case "boot": // we kick the person and shift party left
						if len(curRecord.PartyComp) == 1 || (i == 0 && curRecord.PartyComp[1].PlayerClass == 0) {
							return
						}

						for j, p := range curRecord.PartyComp[i+1:] {
							curRecord.PartyComp[j+i] = p
						}

						curRecord.PartyComp[len(curRecord.PartyComp)-1] = players.PartyMember{}

					}
					ptycon.Players = players.ClassConstructor(curRecord.PartyComp)
					cSelect.Cleanup(i)

				}),
				selector.Cleanup(func(i int) {
					// undraw menu
					for _, p := range pty.Players {
						p.R.Undraw()
						collision.Remove(p.GetSpace())
						mouse.Remove(p.GetSpace())
					}
					cnfrm.Undraw()
					cancl.Undraw()
					boot.Undraw()
					partyBackground.Undraw()

					pty, err = ptycon.NewParty(true)
					if err != nil {
						dlog.Error(err)
						return
					}
					pc.setParty(pty.Players)
					for _, p := range pty.Players {
						collision.Remove(p.GetSpace())
					}

					// end selection
					pc.gameState = playing
					event.Trigger("EndPartySelect", nil)
				}),
				selector.SelectTrigger(key.Down+key.Spacebar),
				selector.SelectTrigger("A"+joystick.ButtonUp),

				selector.InteractTrigger(key.Down+key.B, "boot"),
				selector.InteractTrigger("X"+joystick.ButtonUp, "boot"),

				selector.DestroyTrigger("EndPartySelect"),
				selector.DestroyTrigger(key.Down+key.Escape),
				selector.DestroyTrigger("B"+joystick.ButtonUp),
				selector.MouseBindings(true),

				selector.MouseRight(selector.MouseInteract("boot")),
			)

			return 0
		}
		event.GlobalBind(partySelectStart, key.Down+key.ReturnEnter)
		event.GlobalBind(partySelectStart, "A"+joystick.ButtonUp)

		bkgMusic, err = music.Start(true, "inn3.wav")
		dlog.ErrorCheck(err)

		// Clear, set and report on the debug commands available
		oak.ResetCommands()
		oak.AddCommand("resetParty", func(args []string) {
			curRecord.PartyComp = []players.PartyMember{{players.Swordsman, 0, "Dan the Almost Default"}}
			ptycon.Players = players.ClassConstructor(curRecord.PartyComp)
			ptycon.Players[0].Position = ptyOffset
			for _, p := range pty.Players {
				p.R.Undraw()
				debugTree.Remove(p.RSpace.Space)
			}
			pty, err2 := ptycon.NewParty(true)
			if err2 != nil {
				dlog.Error(err2)
				return
			}

			for _, p := range pty.Players {
				render.Draw(p.R, layer.Play, 4)
			}
		})

		oak.AddCommand("shake", func(args []string) {
			//TODO: determine if default shaker is even noticable

			ss := oak.ScreenShaker{
				Random: false,
				Magnitude: floatgeom.Point2{
					3,
					3,
				},
			}
			ss.Shake(time.Duration(1000) * time.Millisecond)
		})
		oak.AddCommand("debug", func(args []string) {
			if debugTree.DrawDisabled {
				debugTree.DrawDisabled = false
				return
			}
			debugTree.DrawDisabled = true
		})
		fullscreen := false
		oak.AddCommand("fullscreen", func(args []string) {
			fullscreen = !fullscreen
			err := oak.SetFullScreen(fullscreen)
			if err != nil {
				fullscreen = !fullscreen
				fmt.Println(err)
			}
			return
		})
		dlog.Info("Current Debug Commands are: ", strings.Join(oak.GetDebugKeys(), " , "))

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End: func() (string, *scene.Result) {
		curRecord.Store()
		(*bkgMusic).Stop()
		return nextscene, nil
	},
}

func getInteractBtn() render.Renderable {
	// Todo: change to space?
	txt := "Enter"
	if oak.MostRecentInput == oak.Joystick {
		txt = "A"
	}
	keyImg := keyviz.Generator{
		Text:     txt,
		TextSize: 12,
		Color:    color.RGBA{100, 100, 255, 255},
	}.Generate()
	// todo: joystick has different style!
	// todo: composite with icon
	return render.NewSprite(0, 0, keyImg.(*image.RGBA))
}

func getConfirmBtn() render.Renderable {
	txt := "Space"
	if oak.MostRecentInput == oak.Joystick {
		txt = "A"
	}
	keyImg := keyviz.Generator{
		Text:     txt,
		TextSize: 12,
		Color:    color.RGBA{100, 100, 255, 255},
	}.Generate()
	// todo: joystick
	// todo: composite with icon
	return render.NewSprite(0, 0, keyImg.(*image.RGBA))
}

func getCancelBtn() render.Renderable {
	txt := "ESC"
	if oak.MostRecentInput == oak.Joystick {
		txt = "B"
	}
	keyImg := keyviz.Generator{
		Text:     txt,
		TextSize: 12,
		Color:    color.RGBA{255, 100, 100, 255},
	}.Generate()
	// todo: joystick
	// todo: composite with icon
	return render.NewSprite(0, 0, keyImg.(*image.RGBA))
}

func getBootButton() render.Renderable {
	txt := "B"
	if oak.MostRecentInput == oak.Joystick {
		txt = "X"
	}
	keyImg := keyviz.Generator{
		Text:     txt,
		TextSize: 12,
		Color:    color.RGBA{50, 150, 50, 255},
	}.Generate()
	// todo: joystick
	// todo: composite with icon
	return render.NewSprite(0, 0, keyImg.(*image.RGBA))
}
