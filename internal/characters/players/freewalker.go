package players

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/joys"
)

// free walking party parameters TODO: Consider pulling to a true config
const (
	FrameLag     = 10
	MaxPartySize = 4
)

// FreeWalker allows a party to walk around a scene with some following lag
type FreeWalker struct {
	Front     *entities.Interactive
	Followers []*entities.Interactive
	Scale     float64
	LagDeltas [FrameLag * MaxPartySize]floatgeom.Point2
	LagIdx    int
	State     int
}

// FreeWalkControls shares the standard key controls between freealkers
func FreeWalkControls(p *entities.Interactive) {
	p.Delta.Zero()
	lowestID := joys.LowestID()
	js := joys.StickState(lowestID)
	// Todo: support full analog control

	if oak.IsDown(key.UpArrow) || js.StickLY > 8000 {
		p.Delta.Add(physics.NewVector(0, -p.Speed.Y()))
	}
	if oak.IsDown(key.DownArrow) || js.StickLY < -8000 {
		p.Delta.Add(physics.NewVector(0, p.Speed.Y()))
	}
	if oak.IsDown(key.LeftArrow) || js.StickLX < -8000 {
		p.Delta.Add(physics.NewVector(-p.Speed.X(), 0))
	}
	if oak.IsDown(key.RightArrow) || js.StickLX > 8000 {
		p.Delta.Add(physics.NewVector(p.Speed.X(), 0))
	}

}

// FreeFollow can be used to make the rest of the party follow the leader!
func FreeFollow(fw *FreeWalker, p *entities.Interactive) {
	if p.Delta.Magnitude() != 0 {
		//fmt.Println("delta x, y", p.Delta.X(), p.Delta.Y())
		//fmt.Println("lag index", fw.LagIdx)
		// Store this event in our frame delay
		fw.LagDeltas[fw.LagIdx] = floatgeom.Point2{p.Delta.X(), p.Delta.Y()}

		// Access stored frame deltas, move followers
		for i, fo := range fw.Followers {
			delta := fw.LagDeltas[(fw.LagIdx+FrameLag*(i+1))%len(fw.LagDeltas)]
			fo.Vector.Add(physics.NewVector(delta.X(), delta.Y()))
			fo.R.SetPos(fo.Vector.X(), fo.Vector.Y())

			swch := fo.R.(*render.Switch)
			if delta.X() > 0 {
				swch.Set("walkRT")
			} else {
				swch.Set("walkLT")
			}
		}

		// Shift the contents of the frame deltas
		fw.LagIdx--
		if fw.LagIdx < 0 {
			fw.LagIdx = len(fw.LagDeltas) - 1
		}
	} else {
		for _, fo := range fw.Followers {
			swch := fo.R.(*render.Switch)
			cur := swch.Get()
			err := swch.Set("stand" + string(cur[len(cur)-2:]))
			dlog.ErrorCheck(err)
		}
	}

	swch := p.R.(*render.Switch)
	if p.Delta.X() != 0 || p.Delta.Y() != 0 {
		if p.Delta.X() > 0 {
			swch.Set("walkRT")
		} else {
			swch.Set("walkLT")
		}
	} else {
		cur := swch.Get()
		err := swch.Set("stand" + string(cur[len(cur)-2:]))
		dlog.ErrorCheck(err)
	}
}
