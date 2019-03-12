package inn

import (
	"fmt"
	"path/filepath"
	"strconv"

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

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "placeholder_inn.png"))
		render.Draw(innBackground, 0)

		// A way to enter the run
		innDoor := characters.NewDoor()
		iW, iH := innDoor.R.GetDims()
		innDoor.SetPos(float64(oak.ScreenWidth-iW), float64(oak.ScreenHeight-iH)/2) //Center the door on the right side
		render.Draw(innDoor.R, 1)

		text := render.DefFont().NewStrText("Hit the button or walk out of the inn to start the game!", float64(oak.ScreenWidth)/2-120, float64(oak.ScreenHeight)/4-40)
		render.Draw(text, 3, 1)

		innSpace := floatgeom.NewRect2(0, 0, float64(oak.ScreenWidth), float64(oak.ScreenHeight)-32) //Adjusted for the current size of the spearman

		pc := characters.NewPc(characters.JobSwordsman, float64(oak.ScreenWidth)/2, float64(oak.ScreenHeight/2))
		pc.Bind(func(id int, _ interface{}) int {
			ply := pcInnBindings(id)
			move.Limit(ply, innSpace)
			<-pc.RSpace.CallOnHits()
			return 0
		}, "EnterFrame")
		pc.Speed = physics.NewVector(5, 5) // We actually allow players to move around in the inn!

		pc.RSpace.Add(collision.Label(characters.LabelDoor),
			(func(s1, s2 *collision.Space) {
				nextscene = "run"
				stayInMenu = false
			}))

		render.Draw(pc.R, 2, 1)

		menuY += menus.BtnHeightA * 1.5
		menuX -= (menus.BtnWidthA * (.125*float64((characters.CurrentParty.MaxSize%2)) + float64((characters.CurrentParty.MaxSize / 2))))

		selectPC1 := btn.New(menus.BtnCfgA, btn.Layers(2), btn.Pos(menuX, menuY), btn.Text("Change PC 1"), btn.Binding(func(int, interface{}) int {
			pc.R.Undraw()
			pc.SetJob((pc.Job + 1) % characters.JobMax)
			render.Draw(pc.R, 2, 1)
			return 0
		}))

		for x := 2; x <= characters.CurrentParty.MaxSize; x++ {
			menuX += menus.BtnWidthA * 5 / 4
			pctmp := characters.NewPc(characters.JobArcher, float64(pc.X()-18), float64(oak.ScreenHeight/2))
			pctmp.Bind(func(id int, _ interface{}) int {
				pcInnBindings(id)
				return 0
			}, "EnterFrame")
			pctmp.Speed = pc.Speed
			render.Draw(pctmp.R, 2, 1)
			btnTxt := "Change PC " + strconv.Itoa(x)
			btn.New(menus.BtnCfgA, btn.Layers(2), btn.Pos(menuX, menuY), btn.Text(btnTxt), btn.Binding(func(int, interface{}) int {
				fmt.Println("Whoo")
				pctmp.R.Undraw()
				pctmp.SetJob((pctmp.Job + 1) % characters.JobMax)
				render.Draw(pctmp.R, 2, 1)
				return 0
			}))

		}

		fmt.Println("How high are the buttons", exit.Y(), selectPC1.X())

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&nextscene),
}

func pcInnBindings(id int) characters.Player {
	ply, ok := event.GetEntity(id).(characters.Player)
	if !ok {
		dlog.Error("Non-player sent to player binding")
	}

	oldDeltaX := ply.GetDelta().X()
	move.WASD(ply)
	if oldDeltaX != ply.GetDelta().X() { //Lazy impl since we dont have facing baked into WASD movement
		img := ply.GetRenderable().(*render.Switch)
		if ply.GetDelta().X() > 0 {
			img.Set("walkRT")
		} else {
			img.Set("walkLT")
		}
	}
	return ply
}
