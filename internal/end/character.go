package end

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/joys"
	"github.com/oakmound/weekly87/internal/layer"
)

// Consider whether there will be more states, if not consider combining with inn's states
const (
	prePlay = iota
	inMenu
	overridable
	playing
)

type endWalker struct {
	*players.FreeWalker
}

var (
	iEndPartyPosX = float64(oak.ScreenWidth / 2)
	iEndPartyPosY = float64(oak.ScreenHeight)
)

// setParty for the endWalker. Create any nessecary players and update if not
// Is responsible for inital location!
func (iw *endWalker) setParty(plys []*players.Player) {
	if len(plys) == 0 {
		dlog.Error("Need at least one party member")
		return
	}

	if iw.Front == nil {
		iw.Front = entities.NewInteractive(
			float64(oak.ScreenWidth/2),
			float64(oak.ScreenHeight)-80,
			16*iw.Scale,
			32*iw.Scale,
			plys[0].Swtch.Copy().Modify(mod.Scale(iw.Scale, iw.Scale)),
			nil,
			0,
			0,
		)
		iw.bindFront()
	} else {
		old := iw.Front.R
		old.Undraw()
		iw.Front.R = plys[0].Swtch.Copy().Modify(mod.Scale(iw.Scale, iw.Scale))
		iw.Front.R.SetPos(old.X(), old.Y())
		iw.Front.R.(*render.Switch).Set(old.(*render.Switch).Get())
	}
	render.Draw(iw.Front.R, layer.Play, players.MaxPartySize)

	for i := 1; i < len(plys); i++ {
		if i >= len(iw.Followers) {
			// make a new one for this position
			iw.Followers = append(iw.Followers, entities.NewInteractive(
				iw.Front.X(),
				iw.Front.Y(),
				16*iw.Scale,
				32*iw.Scale,
				plys[i].Swtch.Copy().Modify(mod.Scale(iw.Scale, iw.Scale)),
				nil,
				0,
				0,
			))

		} else {
			old := iw.Followers[i-1].R
			old.Undraw()
			iw.Followers[i-1].R = plys[i].Swtch.Copy().Modify(mod.Scale(iw.Scale, iw.Scale))
			iw.Followers[i-1].SetPos(old.X(), old.Y())
			iw.Followers[i-1].R.(*render.Switch).Set(old.(*render.Switch).Get())

		}
		render.Draw(iw.Followers[i-1].R, layer.Play, players.MaxPartySize-1)
	}
}

// newendWalker creates a special character for the inn
func newEndWalker(scale float64, plys []*players.Player) *endWalker {

	iw := &endWalker{
		&players.FreeWalker{Scale: scale},
	}
	iw.setParty(plys)

	return iw
}

func (iw *endWalker) bindFront() {

	iw.Front.Bind(func(id int, _ interface{}) int {
		p, ok := event.GetEntity(id).(*entities.Interactive)
		if !ok {
			dlog.Error("Non-player sent to player binding")
		}

		switch iw.State {
		case inMenu:
			p.Delta.Zero()
			return 0
		case playing:
			players.FreeWalkControls(p)
		case overridable:
			lowestID := joys.LowestID()
			js := joys.StickState(lowestID)
			if oak.IsDown(key.UpArrow) || js.StickLY > 8000 ||
				oak.IsDown(key.DownArrow) || js.StickLY < -8000 ||
				oak.IsDown(key.LeftArrow) || js.StickLX < -8000 ||
				oak.IsDown(key.RightArrow) || js.StickLX > 8000 {
				iw.State = playing
			}
		default:

		}

		p.Vector.Add(p.Delta)

		_, h := p.R.GetDims()
		hf := float64(h)

		if p.Vector.Y() < 32 {
			p.Delta.Sub(physics.NewVector(0, p.Vector.Y()-32))
			p.Vector.SetY(32)
		} else if p.Vector.Y() > (float64(oak.ScreenHeight) - hf) {
			p.Delta.Sub(physics.NewVector(0, p.Vector.Y()-(float64(oak.ScreenHeight)-hf)))
			p.Vector.SetY((float64(oak.ScreenHeight) - hf))
		}
		if p.Vector.X() < 220 {
			p.Delta.Sub(physics.NewVector(p.Vector.X()-220, 0))
			p.Vector.SetX(220)
		} else if p.Vector.X()+p.W > float64(oak.ScreenWidth) {
			p.Delta.Sub(physics.NewVector((p.Vector.X()+p.W)-float64(oak.ScreenWidth), 0))
			p.Vector.SetX(float64(oak.ScreenWidth) - p.W)
		}
		p.R.SetPos(p.Vector.X(), p.Vector.Y())
		p.RSpace.Update(p.Vector.X(), p.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
		<-iw.Front.RSpace.CallOnHits()
		if collision.HitLabel(iw.Front.RSpace.Space, labels.Blocking, labels.NPC) != nil {
			p.Vector.Sub(p.Delta)
			p.Delta.Zero()
			p.R.SetPos(p.Vector.X(), p.Vector.Y())
			p.RSpace.Update(p.Vector.X(), p.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
		}
		players.FreeFollow(iw.FreeWalker, p)

		return 0
	}, "EnterFrame")
	iw.Front.Speed = physics.NewVector(5, 5)

	// iw.Front.RSpace.Add(collision.Label(labels.Door), (func(s1, s2 *collision.Space) {
	// 	d, ok := s2.CID.E().(*doodads.InnDoor)
	// 	if !ok {
	// 		dlog.Error("Non-door sent to inndoor binding")
	// 		return
	// 	}
	// 	nextscene = d.NextScene
	// 	stayInMenu = false
	// }))
}
