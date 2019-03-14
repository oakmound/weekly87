package inn

import (
	"path/filepath"

	klg "github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/audio/filter"

	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/doodads"
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

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "placeholder_inn.png"))
		render.Draw(innBackground, 0)

		// A way to enter the run
		doodads.NewInnDoor()

		innSpace := floatgeom.NewRect2(0, 0, float64(oak.ScreenWidth), float64(oak.ScreenHeight)-32) //Adjusted for the current size of the spearman

		NewInnWalker(innSpace)

		var err error
		music, err = audio.Load(filepath.Join("assets", "audio"), "inn1.wav")
		dlog.ErrorCheck(err)
		music = music.MustFilter(
			filter.Volume(0.5*settings.MusicVolume*settings.MasterVolume),
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
