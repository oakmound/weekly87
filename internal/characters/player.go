package characters

import (
	"errors"
	"fmt"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/physics"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

var requiredPlayerAnimations = []string{
	"standRT",
	"standLT",
	"walkRT",
	"walkLT",
	"deadRT",
	"deadLT",
}

type PlayerConstructor struct {
	Position   floatgeom.Point2
	Speed      floatgeom.Point2
	Dimensions floatgeom.Point2
	// The following strings are required in the animation map:
	// "standRT"
	// "standLT"
	// "walkRT"
	// "walkLT"
	// more may be added
	AnimationMap map[string]render.Modifiable
	Bindings     map[string]func(*Player, interface{}) int
	Special1     func() // Todo: flesh out specials
	Special2     func()
	RunSpeed     float64
}

// Copy returns a shallow copy of the constructor.
// Don't modify the animation map or bindings on a copy of a constructor, that part isn't
// deep copied.
func (pc *PlayerConstructor) Copy() *PlayerConstructor {
	return &PlayerConstructor{
		Position:     pc.Position,
		Speed:        pc.Speed,
		Dimensions:   pc.Dimensions,
		AnimationMap: pc.AnimationMap,
		Bindings:     pc.Bindings,
		Special1:     pc.Special1,
		Special2:     pc.Special2,
		RunSpeed:     pc.RunSpeed,
	}
}

type Player struct {
	*entities.Interactive
	facing      string
	Swtch       *render.Switch
	Alive       bool
	RunSpeed    float64
	ChestValues []int64
}

func (p *Player) Init() event.CID {
	return event.NextID(p)
}

func (p *Player) CheckedBind(bnd func(*Player, interface{}) int, ev string) {
	p.Bind(func(id int, data interface{}) int {
		be, ok := event.GetEntity(id).(*Player)
		if !ok {
			dlog.Error("Player binding was called on non-player")
			return event.UnbindSingle
		}
		return bnd(be, data)
	}, ev)
}

const PlayerWallOffset = 50

func (pc *PlayerConstructor) NewPlayer() (*Player, error) {
	if pc.Dimensions == (floatgeom.Point2{}) {
		return nil, errors.New("Dimensions must be provided")
	}
	for _, s := range requiredPlayerAnimations {
		if _, ok := pc.AnimationMap[s]; !ok {
			return nil, errors.New("Animation name " + s + " must be provided")
		}
	}
	p := Player{}

	p.Swtch = render.NewSwitch("walkRT", pc.AnimationMap)
	p.Interactive = entities.NewInteractive(
		pc.Position.X(),
		pc.Position.Y(),
		pc.Dimensions.X(),
		pc.Dimensions.Y(),
		p.Swtch,
		nil,
		p.Init(),
		0,
	)
	p.facing = "RT"
	p.Alive = true
	p.Speed = physics.NewVector(pc.Speed.X(), pc.Speed.Y())
	p.RunSpeed = pc.RunSpeed

	p.RSpace.UpdateLabel(LabelPC)

	p.CheckedBind(func(p *Player, _ interface{}) int {
		p.facing = "LT"
		return 0
	}, "RunBack")

	p.CheckedBind(func(p *Player, _ interface{}) int {
		p.RunSpeed *= -1
		p.CheckedBind(func(p *Player, _ interface{}) int {
			// Shift the player back until against the right wall
			if int(p.X())-oak.ViewPos.X >= oak.ScreenWidth-PlayerWallOffset {
				return event.UnbindSingle
			}
			p.ShiftX(-p.RunSpeed)
			return 0
		}, "EnterFrame")
		return event.UnbindSingle
	}, "RunBack")

	p.CheckedBind(func(p *Player, _ interface{}) int {
		fmt.Println("Kill triggered")
		dlog.ErrorCheck(p.Swtch.Set("dead" + p.facing))
		return 0
	}, "Kill")

	p.CheckedBind(func(p *Player, _ interface{}) int {
		if !p.Alive {
			p.Swtch.Set("dead" + p.facing)
			// This logic has to change once there are multiple characters
			return 0
		}
		//fmt.Println("Player Loc", p.X(), p.Y())
		// The idea behind splitting up the move functions is
		// flawed when they're all working together--we only want
		// to shift everything -once-, otherwise there are jitters
		// or other awkward bits to moving around.
		p.Delta.Zero()

		p.Delta.SetX(p.RunSpeed)
		if oak.IsDown("W") {
			p.Delta.Add(physics.NewVector(0, -p.Speed.Y()))
		}
		if oak.IsDown("S") {
			p.Delta.Add(physics.NewVector(0, p.Speed.Y()))
		}

		p.Vector.Add(p.Delta)
		oak.SetScreen(oak.ViewPos.X+int(p.RunSpeed), oak.ViewPos.Y)
		//oak.ViewPos.X += int(p.RunSpeed)
		// This is 6, when it should be 32
		//_, h := r.GetDims()
		hf := 32.0
		if p.Vector.Y() < float64(oak.ScreenHeight)*1/3 {
			p.Vector.SetY(float64(oak.ScreenHeight) * 1 / 3)
		} else if p.Vector.Y() > (float64(oak.ScreenHeight) - hf) {
			p.Vector.SetY((float64(oak.ScreenHeight) - hf))
		}
		p.R.SetPos(p.Vector.X(), p.Vector.Y())
		p.RSpace.Update(p.Vector.X(), p.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
		<-p.RSpace.CallOnHits()
		if p.Delta.X() != 0 || p.Delta.Y() != 0 {
			p.Swtch.Set("walk" + p.facing)
		} else {
			p.Swtch.Set("stand" + p.facing)
		}
		return 0
	}, "EnterFrame")

	for ev, b := range pc.Bindings {
		p.CheckedBind(b, ev)
	}
	return &p, nil
}
