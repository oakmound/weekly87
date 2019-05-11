package inn

import (
	"fmt"
	"image/color"
	"path/filepath"

	klg "github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/audio/filter"

	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/dtools"
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
		render.SetDrawStack(
			// ground
			render.NewCompositeR(),
			// entities
			render.NewHeap(false),
			// ui
			render.NewHeap(true),
		)
		debugTree := dtools.NewThickRTree(collision.DefTree, 4)
		debugTree.ColorMap = map[collision.Label]color.RGBA{
			labels.Door:     color.RGBA{200, 0, 100, 255},
			labels.PC:       color.RGBA{125, 0, 255, 255},
			labels.Blocking: color.RGBA{200, 200, 10, 255},
			labels.NPC:      color.RGBA{125, 200, 10, 255},
		}
		render.Draw(debugTree, 2, 1000)

		r = records.Load()

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "placeholder_inn.png"))
		render.Draw(innBackground, 0)

		// A way to enter the run
		doodads.NewInnDoor()

		innSpace := floatgeom.NewRect2(0, 0, float64(oak.ScreenWidth), float64(oak.ScreenHeight)-32) //Adjusted for the current size of the spearman

		// Additional inn aspects
		doodads.NewFurniture(130, 130, 100, float64(oak.ScreenHeight)-130) // Left Table

		doodads.NewFurniture(480, 225, 195, 70) // top Table
		doodads.NewFurniture(480, 430, 185, 70) // bottom Table
		NewInnNPC(players.Mage, 460, 420)

		ptycon := players.PartyConstructor{
			Players: players.ClassConstructor(
				r.PartyComp),
			// []int{players.Spearman, players.Mage, players.Mage, players.Swordsman}),
			// []int{players.Spearman, players.Spearman, players.Spearman, players.Spearman}),
		}
		ptycon.Players[0].Position = floatgeom.Point2{players.WallOffset, 50}
		pty, err2 := ptycon.NewParty(true)
		if err2 != nil {
			dlog.Error(err2)
			return
		}
		for _, p := range pty.Players {
			render.Draw(p.R, 2, 2)
		}

		//Create an example person to navigate the space
		pc := NewInnWalker(innSpace)
		pc.RSpace.Add(labels.NPC, func(_, n *collision.Space) {
			npc, ok := n.CID.E().(*NPC)
			if !ok {
				dlog.Error("Non-npc sent to npc binding")
				return
			}
			if len(r.PartyComp) < 4 {
				fmt.Println("Updating Party")
				r.PartyComp = append(r.PartyComp, npc.Class)
				ptycon.Players = players.ClassConstructor(r.PartyComp)
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
					render.Draw(p.R, 2, 2)
				}
			}

		})

		// Set up the audio
		var err error
		music, err = audio.Load(filepath.Join("assets", "audio"), "inn1.wav")
		dlog.ErrorCheck(err)
		music, err = music.Copy()
		dlog.ErrorCheck(err)
		music = music.MustFilter(
			filter.Volume(0.5*settings.Active.MusicVolume*settings.Active.MasterVolume),
			filter.LoopOn(),
		)

		music.Play()

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End: func() (string, *scene.Result) {
		r.Store()
		music.Stop()
		return nextscene, nil
	},
}
