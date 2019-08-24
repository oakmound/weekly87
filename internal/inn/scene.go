package inn

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	klg "github.com/200sc/klangsynthese/audio"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/key"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/dtools"
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
		doodads.NewCustomInnDoor("startup", 490, 40, 100, 102)

		// Create doodads for tables
		uglymugger, _ := render.LoadSprites("", filepath.Join("16x16", "ugly_mugger.png"), 16, 16, 0)
		prettyMugs := []*render.Sprite{uglymugger[0][0], uglymugger[0][1], uglymugger[1][0]}

		// Block off the top of the inn from being walkable
		doodads.NewFurniture(0, 0, float64(oak.ScreenWidth), 140) // top of inn

		// Additional inn such as tables
		leftT := doodads.NewFurniture(130, 130, 100, float64(oak.ScreenHeight)-130)
		leftT.SetOrnaments(prettyMugs, 3+rand.Intn(7))

		topT := doodads.NewFurniture(480, 230, 190, 60)
		topT.SetOrnaments(prettyMugs, rand.Intn(6))

		botT := doodads.NewFurniture(480, 430, 180, 55)
		botT.SetOrnaments(prettyMugs, rand.Intn(6))

		// Create the notes on the notice board
		noteSpace := floatgeom.NewRect2WH(240, 60, 85, 65)
		noteHeight := 3
		for i := 0; i < 8+rand.Intn(5)*3; i++ {
			noteHeight = doodads.NewNote(noteSpace, noteHeight)
		}

		// Create all possible NPCs
		npcScale := 1.6
		npcs := []NPC{
			NewInnNPC(players.Swordsman, npcScale, 680, 230).FaceLeft(true),
			NewInnNPC(players.Mage, npcScale, 440, 210),
			NewInnNPC(players.WhiteMage, npcScale, 670, 423).FaceLeft(true),
			NewInnNPC(players.Berserker, npcScale, 445, 460),
			NewInnNPC(players.BlueMage, npcScale, 241, 210).FaceLeft(true),
			NewInnNPC(players.Paladin, npcScale, 240, 280).FaceLeft(true),
			NewInnNPC(players.Spearman, npcScale, 243, 400).FaceLeft(true),
			NewInnNPC(players.TimeMage, npcScale, 675, 477).FaceLeft(true),
		}

		// Inn does quite a few operations on our record (mainly for party purposes)
		curRecord = records.Load()

		// Simple metric for determining number of NPCs in room
		progress := int(math.Min(float64(curRecord.SectionsCleared)/10.0, float64(len(npcs))))
		futureNpcs := npcs[progress:len(npcs)]
		npcs = npcs[0:progress]
		for _, fn := range futureNpcs {
			fn.Destroy() // Being extra safe
		}
		for _, np := range npcs {
			np.Activate()
		}

		// Create the player and the display of the party at the top of the screen
		ptycon := players.PartyConstructor{
			Players: players.ClassConstructor(curRecord.PartyComp),
		}
		// Draw the party in top left
		partyBackground, _ := render.LoadSprite("", filepath.Join("raw", "selector_background.png"))

		partyBackground.SetPos(30, 20)

		ptyOffset := floatgeom.Point2{players.WallOffset, 30}
		ptycon.Players[0].Position = ptyOffset
		pty, err := ptycon.NewParty(true)
		if err != nil {
			dlog.Error(err)
			return
		}

		//interactDelay := time.Second
		//pcLastInteract := time.Now()
		//interactLock := &sync.Mutex{}
		//Create an example person to navigate the space
		pc := NewInnWalker(npcScale, pty.Players)

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
			interactBtn := getConfirmBtn()
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

		event.GlobalBind(func(int, interface{}) int {
			if lastInteractedNPC != nil {
				npc := lastInteractedNPC
				npc.Button.Undraw()
				npcW, _ := npc.R.GetDims()
				bkgW, bkgH := partyBackground.GetDims()
				// Disable all controls
				pc.inMenu = true
				// Spawn a box with the party above the npc being selected
				partyBackground.SetPos(npc.X()+float64(npcW-bkgW)/2, npc.Y()-10-float64(bkgH))
				ptycon.Players[0].Position = floatgeom.Point2{partyBackground.X() + 20, partyBackground.Y() + 10}
				render.Draw(partyBackground, layer.UI, 1)

				pty, err := ptycon.NewParty(true)
				dlog.ErrorCheck(err)
				spcs := make([]*collision.Space, len(pty.Players))
				for i, p := range pty.Players {
					render.Draw(p.R, layer.UI, 2)
					spcs[i] = p.GetSpace()
				}

				// Show a confirm button (and a cancel button)
				cnfrm := getConfirmBtn()
				cnfrm.SetPos(partyBackground.X(), partyBackground.Y()+float64(bkgH)+2)
				render.Draw(cnfrm, layer.UI, 1)
				cancl := getCancelBtn()
				canclW, _ := cancl.GetDims()
				cancl.SetPos(partyBackground.X()+float64(bkgW-canclW), partyBackground.Y()+float64(bkgH)+2)
				render.Draw(cancl, layer.UI, 1)

				// Let arrow keys / joystick or mouse even control which party member is selected
				selector.New(
					selector.Layers(layer.UI, 3),
					selector.HorzArrowControl(),
					selector.JoystickHorzDpadControl(),
					selector.Spaces(spcs...),
					selector.Callback(func(i int) {
						// modify party
						curRecord.PartyComp[i] = npc.Class
						ptycon.Players = players.ClassConstructor(curRecord.PartyComp)

					}),
					selector.Cleanup(func(i int) {
						// undraw menu
						for _, p := range pty.Players {
							p.R.Undraw()
							collision.Remove(p.GetSpace())
						}
						cnfrm.Undraw()
						cancl.Undraw()
						partyBackground.Undraw()

						pty, err = ptycon.NewParty(true)
						if err != nil {
							dlog.Error(err)
							return
						}
						pc.SetParty(pty.Players)

						// end selection
						pc.inMenu = false
						event.Trigger("EndPartySelect", nil)
					}),
					selector.SelectTrigger(key.Down+key.Spacebar),
					selector.SelectTrigger("A"+joystick.ButtonUp),
					selector.DestroyTrigger("EndPartySelect"),
					selector.DestroyTrigger(key.Down+key.Escape),
					selector.DestroyTrigger("B"+joystick.ButtonUp),
				)
			}
			return 0
		}, key.Down+key.ReturnEnter)

		bkgMusic, err = music.Start(true, "inn1.wav")
		dlog.ErrorCheck(err)

		// Clear, set and report on the debug commands available
		oak.ResetCommands()
		oak.AddCommand("resetParty", func(args []string) {
			curRecord.PartyComp = []int{players.Spearman}
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

func getConfirmBtn() render.Renderable {
	// todo: joystick
	return render.NewColorBox(50, 20, color.RGBA{100, 100, 255, 255})
}

func getCancelBtn() render.Renderable {
	// todo: joystick
	return render.NewColorBox(50, 20, color.RGBA{255, 100, 100, 255})
}
