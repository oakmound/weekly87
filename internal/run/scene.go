package run

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/move"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters"
)

var stayInGame bool
var nextscene string
var playerMoveRect floatgeom.Rect2

// facing is whether is game is moving forward or backward,
// 1 means forward, -1 means backward
var facing = 1
var runSpeed = 2.0

// Scene  to display the run
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInGame = true
		nextscene = "inn"

		// There should be some way to draw to a stack based
		// on layer name
		render.SetDrawStack(
			// ground
			render.NewCompositeR(),
			// maybe background / parallax?
			// wall
			render.NewCompositeR(),
			// entities
			render.NewHeap(false),
			// maybe effects?
			// ui
			render.NewHeap(true),
			render.NewDrawFPS(),
			render.NewLogicFPS(),
		)

		playerMoveRect = floatgeom.NewRect2(0, float64(oak.ScreenHeight)*1/3, float64(oak.ScreenWidth),
			float64(oak.ScreenHeight))
		// Todo: add collision with chests, when this happpens the chest
		// 1. needs to be collected
		// 2. If we're going forward, start going back
		// 3. Shift the player move rect gradually if we just started moving back
		// 4. Flip enemies / characters as needed

		s := characters.NewSpearman(50, float64(oak.ScreenHeight/2))
		s.Bind(func(id int, _ interface{}) int {
			ply, ok := event.GetEntity(id).(characters.Player)
			if !ok {
				dlog.Error("Non-player sent to player binding")
			}
			// The idea behind splitting up the move functions
			// flawed when they're all working together--we only want
			// to shift everything -once-, otherwise there are jitters
			// or other awkward bits to moving around.
			vec := ply.Vec()
			delta := ply.GetDelta()
			spd := ply.GetSpeed()
			r := ply.GetRenderable()
			sp := ply.GetSpace()

			delta.Zero()

			delta.SetX(runSpeed)
			if oak.IsDown("W") {
				delta.Add(physics.NewVector(0, -spd.Y()))
			}
			if oak.IsDown("S") {
				delta.Add(physics.NewVector(0, spd.Y()))
			}

			vec.Add(delta)
			oak.ViewPos.X += int(runSpeed)
			// This is 6, when it should be 32
			//_, h := r.GetDims()
			hf := 32.0
			if vec.Y() < playerMoveRect.Min.Y() {
				vec.SetY(playerMoveRect.Min.Y())
			} else if vec.Y() > playerMoveRect.Max.Y()-hf {
				vec.SetY(playerMoveRect.Max.Y() - hf)
			}
			r.SetPos(vec.X(), vec.Y())
			sp.Update(vec.X(), vec.Y(), sp.GetW(), sp.GetH())
			return 0
		}, "EnterFrame")
		render.Draw(s.R, 2, 2)

		// todo populate baseseed
		tracker := NewSectionTracker(baseSeed)
		sct := tracker.Next()
		sct.Draw()
		nextSct := tracker.Next()
		nextSct.SetBackgroundX(sct.X() + sct.W())
		nextSct.Draw()

		event.GlobalBind(func(int, interface{}) int {
			//sct.Shift(runSpeed * float64(-facing))
			//nextSct.Shift(runSpeed * float64(-facing))
			// This calculation needs to be modified based
			// on how much of the screen a section takes up.
			// If a section takes up more than one screen,
			// this is fine, otherwise it needs to change a little
			w := sct.W()
			offLeft := oak.ViewPos.X - int(w)
			if offLeft >= 0 {
				// We need a way to make these actions draw-level atomic
				// Or a way to fake it so there isn't a blip
				oak.ViewPos.X = offLeft
				nextSct.Shift(-w)
				sct.Destroy()
				sct = nextSct
				// Todo: shift player, not locally stored s
				s.ShiftX(-w)
				go func() {
					nextSct = tracker.Next()
					nextSct.SetBackgroundX(sct.X() + w)
					nextSct.Draw()
				}()
			}
			return 0
		}, "EnterFrame")

		// Maybe there's a countdown timer

		// There should be a player running to the right from the left
		// side of the screen

		// We need a scrolling background

		// The state of the game is generated based on combining a base seed
		// and the current section the player is in. When the game is first started
		// base seed is populated randomly and stored in a settings file, then
		// incremented as sections are cleared

		// We also need to keep track of changes to each section like enemies destroyed
		// This means map[int64][]int, where the slice is list of enemies destroyed
		// with enemies identified by order they are made in

		// The inn should have some image on the left of the first section

		// Background should probably be very basic hallway with tile types
		// and different themes populate the tile types

		// Character types: 10 Sec cooldown ability, 30 sec
		// Spearman - Shove up - Attack in front
		// Warrior - Shove back - Shove all enemies back
		// Cleric - Slow down run speed for short period - Revive
		// Ranger - Y Speed boost for short period - Shoot arrow in front
		// Rogue - Invisible for short period - Blink / jump forward
		// Paladin - Invincible for short period -
		// Mage - Spawns Fire - Freeze all enemies in place

		// Enemy types:
		// 1. Stands in the way and hurts if you touch it
		// 2. Goes up and down and hurts if you touch it
		// 3. Charges right to left and hurts if you touch it
		// 3a. As 3, but can slightly curve during charge

		// Treasure types:
		// bigger and fancier boxes with fancier colors
		// maybe 3 sizes
		// Color order: brown, then dim gray, then black,
		// white, yellow, green, blue, purple, orange, red, silver

		// for right now the one character can chain boxes behind them
		// infinitely, eventually there should be an upgrade thing

		// Sections jitter a bit and it is annoying, must fix

	},
	Loop: scene.BooleanLoop(&stayInGame),
	// scene.GoTo("inn"),
	End: scene.GoToPtr(&nextscene),
}

func ShiftMoverX(mvr move.Mover, x float64) {
	vec := mvr.Vec()
	vec.ShiftX(x)
	mvr.GetRenderable().SetPos(vec.X(), vec.Y())
	sp := mvr.GetSpace()
	sp.Update(vec.X(), vec.Y(), sp.GetW(), sp.GetH())
}
