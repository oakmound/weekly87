package inn

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"path/filepath"
	"strings"
	"sync"
	"time"

	klg "github.com/200sc/klangsynthese/audio"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render/mod"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/dtools"
	"github.com/oakmound/weekly87/internal/layer"
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

		innSpace := floatgeom.NewRect2(0, 0, float64(oak.ScreenWidth), float64(oak.ScreenHeight)-32) //Adjusted for the current size of the spearman

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
		partyBackground := render.NewColorBox(206, 52, color.RGBA{90, 90, 200, 255})
		partyBackground.SetPos(30, 20)
		render.Draw(partyBackground, layer.Play, 3)

		ptyOffset := floatgeom.Point2{players.WallOffset, 30}
		ptycon.Players[0].Position = ptyOffset
		pty, err := ptycon.NewParty(true)
		if err != nil {
			dlog.Error(err)
			return
		}
		for _, p := range pty.Players {
			render.Draw(p.R, layer.Play, 4)
		}

		interactDelay := time.Second
		pcLastInteract := time.Now()
		interactLock := &sync.Mutex{}
		//Create an example person to navigate the space
		pc := NewInnWalker(innSpace, npcScale, pty.Players[0].Swtch.Copy().(*render.Switch))
		// Interact with NPCs
		pc.RSpace.Add(labels.NPC, func(_, n *collision.Space) {
			// Limit interaction rate of player
			interactLock.Lock()
			if pcLastInteract.Add(interactDelay).After(time.Now()) {
				interactLock.Unlock()
				return
			}
			npc, ok := n.CID.E().(*NPC)
			if !ok {
				interactLock.Unlock()
				dlog.Error("Non-npc sent to npc binding")
				return
			}
			pcLastInteract = time.Now()
			interactLock.Unlock()

			dlog.Info("Adding a class to the party")
			curRecord.PartyComp = append(curRecord.PartyComp, npc.Class)
			for _, p := range pty.Players {
				p.Destroy()
				p.R.Undraw()
			}
			if len(curRecord.PartyComp) > 4 {
				curRecord.PartyComp = curRecord.PartyComp[1:]
			}
			ptycon.Players = players.ClassConstructor(curRecord.PartyComp)
			ptycon.Players[0].Position = ptyOffset

			pty, err = ptycon.NewParty(true)
			if err != nil {
				dlog.Error(err)
				return
			}
			pc.R.Undraw()
			pc.R = pty.Players[0].Swtch.Copy().Modify(mod.Scale(npcScale, npcScale))
			render.Draw(pc.R, layer.Play, 2)

			for _, p := range pty.Players {
				render.Draw(p.R, layer.Play, 4)
			}

		})

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
