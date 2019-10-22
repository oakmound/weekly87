package end

import (
	"time"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/layer"
	"github.com/oakmound/weekly87/internal/sfx"
)

// presentSpoils gained in your run
// leave in style
func presentSpoils(party *players.Party, graveCount *int, index int) {
	if index > len(party.Players)-1 {
		// we are out of players  to present
		investigate(party)
		return
	}
	p := party.Players[index]
	p.CID = p.Init()

	// Safely copy assets
	s := p.Swtch.Copy()
	p.Swtch = s.Modify(mod.Scale(npcScale, npcScale)).(*render.Switch)
	p.Interactive.R = p.Swtch

	//Player enters stage right
	p.SetPos(float64(oak.ScreenWidth-64), startY)
	render.Draw(p.R, layer.Play, 20)
	p.ChestsHeight = 0
	// TODO: Consider centering chests to account for varying sizes
	for _, r := range p.Chests {
		_, h := r.GetDims()
		p.ChestsHeight += float64(h)
		chestHeight := p.ChestsHeight
		r.(*render.Sprite).Vector = r.Attach(p.Vector, -3, -chestHeight)
		render.Draw(r, layer.Play, 21)
	}

	dlog.Verb("Character %d walking through and is Alive:%t ", index, p.Alive)

	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		ply.ShiftPos(-2, 0)

		p.Swtch.Set("walkLT")
		if len(ply.ChestValues) > 0 {
			p.Swtch.Set("walkHold")
		}
		if !ply.Alive {
			p.Swtch.Set("deadLT")

		}

		// TODO:
		// If alive: throw chests with hop and cheer
		//		Card explaining the amount in their chests?
		//		Chest explodes into money and goes into pit
		// If dead: whomp whomp
		// 		eulogoy? name, class, run?

		// Next person in party starts process
		// You walk to end point (graves or bottom)
		// When reach your end point destroy self

		//Move to the center
		if ply.R.X() < presentationX {

			if p.Alive {
				p.Swtch.Set("standRT")
				if len(p.ChestValues) > 0 {
					hop(p)

				} else {
					sfx.Play("ohWell")
				}

				t := time.Now().Add(time.Second)
				p.CheckedBind(func(ply *players.Player, _ interface{}) int {

					if time.Now().After(t) {
						livingExit(p)
						return event.UnbindSingle
					}
					return 0
				}, "EnterFrame")
			} else {
				*graveCount++
				deadMovement(p)
			}

			// kick off the next persons presentation
			presentSpoils(party, graveCount, index+1)
			return event.UnbindSingle

		}
		return 0
	}, "EnterFrame")

}

// deadMovement into the grave
func deadMovement(p *players.Player) {
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		// ply.R.Undraw()

		ply.ShiftPos(-2, 0)
		if ply.R.X() < graveX {
			deathSprites(ply.R.X(), ply.R.Y())
			sfx.Play("dissappear1")
			ply.R.Undraw()
			return event.UnbindSingle
		}

		return 0
	}, "EnterFrame")

}

// livingExit is for those great people who actually survived
func livingExit(p *players.Player) {
	p.Swtch.Set("walkLT")
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		ply.ShiftPos(0, 2)
		if ply.R.Y() > float64(oak.ScreenHeight) {

			return event.UnbindSingle
		}
		return 0
	}, "EnterFrame")
}

// hop up before throwing your chest
// remember to come down!
func hop(p *players.Player) {
	sfx.Play("chestHop1")
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {
		if ply.Y() < presentationY-hopDistance {
			tossChests(ply)
			ply.CheckedBind(func(plyz *players.Player, _ interface{}) int {
				if plyz.Y() > presentationY {
					return event.UnbindSingle
				}
				plyz.ShiftPos(0, 4)
				return 0
			}, "EnterFrame")
			return event.UnbindSingle
		}

		ply.ShiftPos(0, -4)
		return 0
	}, "EnterFrame")
}

// tossChests the player is carrying to the money pit
// make them explode when they get there
func tossChests(p *players.Player) {

	for _, c := range p.Chests {
		c.(*render.Sprite).Vector = c.(*render.Sprite).Vector.Detach()
	}
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		done := false
		for _, cr := range p.Chests {

			cr.ShiftX((pitX - presentationX) / 100.0)
			cr.ShiftY((pitY - presentationY + hopDistance) / 100.0)

			if cr.X() > pitX {
				explodeChest(cr.X(), cr.Y())
				cr.Undraw()
				done = true
			}
		}
		if done {
			return event.UnbindSingle
		}
		return 0
	}, "EnterFrame")
}
