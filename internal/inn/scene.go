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
	"github.com/200sc/klangsynthese/audio/filter"

	"github.com/oakmound/oak/audio"
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
	"github.com/oakmound/weekly87/internal/records"
	"github.com/oakmound/weekly87/internal/settings"
)

var stayInMenu bool
var nextscene string
var music klg.Audio
var r *records.Records

// Scene  to display the inn
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInMenu = true
		nextscene = "inn"
		render.SetDrawStack(layer.Get()...)
		debugTree := dtools.NewThickRTree(collision.DefTree, 4)
		debugTree.ColorMap = labels.ColorMap
		render.Draw(debugTree, 2, 1000)

		r = records.Load()

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "placeholder_inn.png"))
		render.Draw(innBackground, 0)

		// A way to enter the run
		doodads.NewInnDoor("run")
		doodads.NewCustomInnDoor("startup", 490, 40, 100, 102)

		innSpace := floatgeom.NewRect2(0, 0, float64(oak.ScreenWidth), float64(oak.ScreenHeight)-32) //Adjusted for the current size of the spearman

		uglymugger, _ := render.LoadSprites("", filepath.Join("16x16", "ugly_mugger.png"), 16, 16, 0)
		prettyMugs := []*render.Sprite{uglymugger[0][0], uglymugger[0][1], uglymugger[1][0]}

		doodads.NewFurniture(0, 0, float64(oak.ScreenWidth), 140) // top of inn

		// Additional inn aspects
		doodads.NewFurniture(130, 130, 100, float64(oak.ScreenHeight)-130) // Left Table
		for i := 0; i < 2+rand.Intn(7); i++ {
			doodads.NewOrnament(130, 130, 100, float64(oak.ScreenHeight)-130, prettyMugs[rand.Intn(3)])
		}

		topLocations := []float64{480, 230, 190, 60}
		botLocations := []float64{480, 430, 180, 55}

		doodads.NewFurniture(topLocations[0], topLocations[1], topLocations[2], topLocations[3]) // top Table
		for i := 0; i < rand.Intn(6); i++ {
			doodads.NewOrnament(topLocations[0], topLocations[1], topLocations[2], topLocations[3], prettyMugs[rand.Intn(3)])
		}
		doodads.NewFurniture(botLocations[0], botLocations[1], botLocations[2], botLocations[3]) // bottom Table
		for i := 0; i < rand.Intn(6); i++ {
			doodads.NewOrnament(botLocations[0], botLocations[1], botLocations[2], botLocations[3], prettyMugs[rand.Intn(3)])
		}

		noteSpace := floatgeom.NewRect2WH(240, 60, 85, 65)

		noteHeight := 3
		for i := 0; i < 8+rand.Intn(5)*3; i++ {
			noteHeight = doodads.NewNote(noteSpace, noteHeight)
		}

		npcs := []NPC{}

		npcScale := 1.6

		npcs = append(npcs, NewInnNPC(players.Swordsman, npcScale, 680, 230).FaceLeft(true))
		npcs = append(npcs, NewInnNPC(players.Mage, npcScale, 440, 210))

		npcs = append(npcs, NewInnNPC(players.WhiteMage, npcScale, 670, 423).FaceLeft(true))
		npcs = append(npcs, NewInnNPC(players.Berserker, npcScale, 445, 460))

		npcs = append(npcs, NewInnNPC(players.BlueMage, npcScale, 241, 210).FaceLeft(true))
		npcs = append(npcs, NewInnNPC(players.Paladin, npcScale, 240, 280).FaceLeft(true))

		npcs = append(npcs, NewInnNPC(players.Spearman, npcScale, 243, 400).FaceLeft(true))
		npcs = append(npcs, NewInnNPC(players.TimeMage, npcScale, 675, 477).FaceLeft(true))

		// For now lets use a combined metric for progress of npc unlocks
		progress := int(math.Min(float64(r.SectionsCleared)/10.0, float64(len(npcs))))

		futureNpcs := npcs[progress:len(npcs)]
		npcs = npcs[0:progress]

		dlog.Verb("Future NPCS", len(futureNpcs))
		for _, fn := range futureNpcs {
			// fn.RSpace = nil
			fn.Destroy()
		}
		for _, np := range npcs {
			np.Activate()
		}

		ptycon := players.PartyConstructor{
			Players: players.ClassConstructor(
				r.PartyComp),
			// []int{players.Spearman, players.Mage, players.Mage, players.Swordsman}),
			// []int{players.Spearman, players.Spearman, players.Spearman, players.Spearman}),
		}
		partyBackground := render.NewColorBox(206, 52, color.RGBA{90, 90, 200, 255})
		partyBackground.SetPos(30, 20)
		render.Draw(partyBackground, 2, 3)
		ptyOffset := floatgeom.Point2{players.WallOffset, 30}
		ptycon.Players[0].Position = ptyOffset
		pty, err := ptycon.NewParty(true)
		if err != nil {
			dlog.Error(err)
			return
		}
		for _, p := range pty.Players {
			render.Draw(p.R, 2, 4)
		}

		interactDelay := time.Second
		pcLastInteract := time.Now()
		interactLock := &sync.Mutex{}
		//Create an example person to navigate the space
		pc := NewInnWalker(innSpace, npcScale, pty.Players[0].Swtch.Copy().(*render.Switch))
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
			r.PartyComp = append(r.PartyComp, npc.Class)
			for _, p := range pty.Players {
				p.Destroy()
				p.R.Undraw()
			}
			if len(r.PartyComp) > 4 {
				r.PartyComp = r.PartyComp[1:]
			}
			ptycon.Players = players.ClassConstructor(r.PartyComp)
			ptycon.Players[0].Position = ptyOffset

			pty, err = ptycon.NewParty(true)
			if err != nil {
				dlog.Error(err)
				return
			}
			pc.R.Undraw()
			pc.R = pty.Players[0].Swtch.Copy().Modify(mod.Scale(npcScale, npcScale))
			render.Draw(pc.R, 2, 2)

			for _, p := range pty.Players {
				render.Draw(p.R, 2, 4)
			}

		})

		// err = mouse.PhaseCollision()
		dlog.ErrorCheck(err)

		// Set up the audio
		music, err = audio.Load(filepath.Join("assets", "audio"), "inn1.wav")
		dlog.ErrorCheck(err)
		music, err = music.Copy()
		dlog.ErrorCheck(err)
		music = music.MustFilter(
			filter.Volume(0.5*settings.Active.MusicVolume*settings.Active.MasterVolume),
			filter.LoopOn(),
		)

		music.Play()

		oak.ResetCommands()
		oak.AddCommand("resetParty", func(args []string) {
			r.PartyComp = []int{players.Spearman}
			ptycon.Players = players.ClassConstructor(r.PartyComp)
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
				render.Draw(p.R, 2, 4)
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
		r.Store()
		music.Stop()
		return nextscene, nil
	},
}
