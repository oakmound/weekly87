package players

import (
	"errors"
	"math"
	"strconv"

	"github.com/oakmound/oak/key"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters/labels"
)

type Party struct {
	event.CID
	Players      []*Player
	Acceleration float64
	speedUps     float64
}

func (p *Party) Init() event.CID {
	return event.NextID(p)
}

func (p *Party) SpeedUp() {
	// 100 sections to get to 20 accel
	// 50 sections to get to 15 accel
	// 25 sections to get to 10 accel
	// 12 sections to get to 5 accel
	p.speedUps++
	p.Acceleration = math.Log10(math.Pow(
		math.Log10(p.speedUps+10), 2)) * 15
	if p.Players[0].RunSpeed == 0 {
		p.Acceleration = 0
	}

}

func (p *Party) CheckedBind(bnd func(*Party, interface{}) int, ev string) {
	p.Bind(func(id int, data interface{}) int {
		be, ok := event.GetEntity(id).(*Party)
		if !ok {
			dlog.Error("Party binding was called on non-party")
			return event.UnbindSingle
		}
		return bnd(be, data)
	}, ev)
}

func (p *Party) RunSpeed() int {
	if p.Players[0].facing == "LT" {
		return int(p.Players[0].RunSpeed - p.Acceleration)
	}
	return int(p.Players[len(p.Players)-1].RunSpeed + p.Acceleration)
}

func (p *Party) Speed() physics.Vector {
	if p.Players[0].facing == "LT" {
		return p.Players[0].Speed
	}
	return p.Players[len(p.Players)-1].Speed
}

func (p *Party) Defeated() bool {
	for _, pl := range p.Players {
		if pl.Alive {
			return false
		}
	}
	return true
}

func (p *Party) ShiftX(f float64) {
	for _, pl := range p.Players {
		pl.ShiftX(f)
	}
}

type PartyConstructor struct {
	Players  []Constructor
	Bindings map[string]func(*Party, interface{}) int
}

func (pc *PartyConstructor) NewParty() (*Party, error) {
	if len(pc.Players) == 0 {
		return nil, errors.New("At least one player must be in a party")
	}

	pty := &Party{}

	const PlayerGap = 50

	for i, pcon := range pc.Players {
		if pcon.Dimensions == (floatgeom.Point2{}) {
			return nil, errors.New("Dimensions must be provided for player " + strconv.Itoa(i))
		}
		for _, s := range requiredAnimations {
			if _, ok := pcon.AnimationMap[s]; !ok {
				return nil, errors.New("Animation name " + s + " must be provided for player " + strconv.Itoa(i))
			}
		}
		p := Player{}

		if pcon.Special1 != nil {
			p.Special1 = pcon.Special1.SetUser(&p)

		}
		if pcon.Special2 != nil {
			p.Special2 = pcon.Special2.SetUser(&p)
		}

		p.Swtch = render.NewSwitch("walkRT", pcon.AnimationMap)
		p.Interactive = entities.NewInteractive(
			pc.Players[0].Position.X()+float64(i)*PlayerGap,
			pc.Players[0].Position.Y(),
			pcon.Dimensions.X(),
			pcon.Dimensions.Y(),
			p.Swtch,
			nil,
			p.Init(),
			0,
		)
		p.facing = "RT"
		p.Alive = true
		p.Speed = physics.NewVector(pcon.Speed.X(), pcon.Speed.Y())
		p.RunSpeed = pcon.RunSpeed
		p.RSpace.UpdateLabel(labels.PC)

		p.CheckedBind(func(p *Player, _ interface{}) int {
			p.facing = "LT"
			return 0
		}, "RunBack")

		p.CheckedBind(func(p *Player, _ interface{}) int {
			dlog.ErrorCheck(p.Swtch.Set("dead" + p.facing))
			return 0
		}, "Kill")

		for ev, b := range pcon.Bindings {
			p.CheckedBind(b, ev)
		}
		pty.Players = append(pty.Players, &p)
	}

	pty.CID = pty.Init()

	pty.CheckedBind(func(pty *Party, _ interface{}) int {
		for i, p := range pty.Players {
			i := i
			p.RunSpeed *= -1
			p.ForcedInvulnerable = true
			p.CheckedBind(func(p *Player, _ interface{}) int {
				// Shift the player back until against the right wall
				if int(p.X())-oak.ViewPos.X >= oak.ScreenWidth-(WallOffset+(len(pty.Players)-1-i)*PlayerGap) {
					p.ForcedInvulnerable = false
					return event.UnbindSingle
				}
				p.ShiftX(float64(-pty.RunSpeed()) * 2)
				return 0
			}, "EnterFrame")
		}
		return event.UnbindSingle
	}, "RunBack")

	pty.CheckedBind(func(pty *Party, _ interface{}) int {
		p0 := pty.Players[0]
		p0.Delta.Zero()

		p0.Delta.SetX(float64(pty.RunSpeed()))
		if oak.IsDown(key.UpArrow) {
			p0.Delta.ShiftY(-pty.Speed().Y())
		}
		if oak.IsDown(key.DownArrow) {
			p0.Delta.ShiftY(pty.Speed().Y())
		}

		p0.Vector.Add(p0.Delta)

		_, h := p0.Swtch.GetDims()
		hf := float64(h)
		if p0.Vector.Y() < float64(oak.ScreenHeight)*1/3 {
			p0.Vector.SetY(float64(oak.ScreenHeight) * 1 / 3)
		} else if p0.Vector.Y() > (float64(oak.ScreenHeight) - hf) {
			p0.Vector.SetY((float64(oak.ScreenHeight) - hf))
		}

		for i, p := range pty.Players {
			// The idea behind splitting up the move functions is
			// flawed when they're all working together--we only want
			// to shift everything -once-, otherwise there are jitters
			// or other awkward bits to moving around.
			if i != 0 {
				p.Vector.Add(p0.Delta)
				p.Vector.SetY(p0.Vector.Y())
			}
			p.R.SetPos(p.Vector.X(), p0.Vector.Y())

			if !p.Alive {
				p.Swtch.Set("dead" + p.facing)
			} else {
				p.RSpace.Update(p.Vector.X(), p.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
				<-p.RSpace.CallOnHits()
				if p0.Delta.X() != 0 || p0.Delta.Y() != 0 {
					if len(p.ChestValues) > 0 {
						p.Swtch.Set("walkHold")
					} else {
						p.Swtch.Set("walk" + p.facing)
					}
				} else {
					if len(p.ChestValues) > 0 {
						p.Swtch.Set("standHold")
					} else {
						p.Swtch.Set("stand" + p.facing)
					}
				}
			}
		}

		oak.SetScreen(oak.ViewPos.X+int(pty.RunSpeed()), oak.ViewPos.Y)

		return 0
	}, "EnterFrame")

	for ev, b := range pc.Bindings {
		pty.CheckedBind(b, ev)
	}

	return pty, nil
}
