package inn

import (
	"fmt"
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/entities/x/move"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/menus"
)

var stayInMenu bool
var nextscene string

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
			//ui text
			render.NewHeap(true),
		)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) / 4

		exit := btn.New(menus.BtnCfgA, btn.Layers(2), btn.Pos(menuX, menuY), btn.Text("Start Run"), btn.Binding(func(int, interface{}) int {
			nextscene = "run"
			stayInMenu = false
			return 0
		}))

		fmt.Println("How high are the buttons", exit.Y())

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "placeholder_inn.png"))
		render.Draw(innBackground, 0)

		// A way to enter the run
		innDoor := characters.NewDoor()
		iW, iH := innDoor.R.GetDims()
		innDoor.SetPos(float64(oak.ScreenWidth-iW), float64(oak.ScreenHeight-iH)/2) //Center the door on the right side
		render.Draw(innDoor.R, 1)

		innSpace := floatgeom.NewRect2(0, 0, float64(oak.ScreenWidth), float64(oak.ScreenHeight)-32) //Adjusted for the current size of the spearman

		// TODO: remove the spearman from here
		s := characters.NewSpearman(float64(oak.ScreenWidth)/2, float64(oak.ScreenHeight/2))
		s.Bind(func(id int, _ interface{}) int {
			ply, ok := event.GetEntity(id).(characters.Player)
			if !ok {
				dlog.Error("Non-player sent to player binding")
			}

			move.WASD(ply)
			move.Limit(ply, innSpace)
			<-s.RSpace.CallOnHits()
			//collision.HitLabel()
			return 0
		}, "EnterFrame")
		s.Speed = physics.NewVector(5, 5) // We actually allow players to move around in the inn!

		s.RSpace.Add(collision.Label(characters.LabelDoor),
			(func(s1, s2 *collision.Space) {
				nextscene = "run"
				stayInMenu = false
			}))

		render.Draw(s.R, 2, 1)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&nextscene),
}
