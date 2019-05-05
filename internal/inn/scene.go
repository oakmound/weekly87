package inn

import (
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
	"github.com/oakmound/weekly87/internal/dtools"
	"github.com/oakmound/weekly87/internal/settings"
)

var stayInMenu bool
var nextscene string
var music klg.Audio

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
			labels.Blocking: color.RGBA{125, 200, 10, 255},
		}
		render.Draw(debugTree, 2, 1000)

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "placeholder_inn.png"))
		render.Draw(innBackground, 0)

		// A way to enter the run
		doodads.NewInnDoor()

		innSpace := floatgeom.NewRect2(0, 0, float64(oak.ScreenWidth), float64(oak.ScreenHeight)-32) //Adjusted for the current size of the spearman

		//Create an example person to navigate the space
		NewInnWalker(innSpace)

		// Additional inn aspects
		doodads.NewFurniture(130, 130, 100, float64(oak.ScreenHeight)-130) // Left Table

		doodads.NewFurniture(480, 225, 195, 70) // top Table
		doodads.NewFurniture(480, 430, 185, 70) // bottom Table

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
		music.Stop()
		return nextscene, nil
	},
}
