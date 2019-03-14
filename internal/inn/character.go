package inn

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
)

func NewInnWalker(innSpace floatgeom.Rect2) {
	anims := players.SpearmanConstructor.AnimationMap

	s := entities.NewInteractive(
		float64(oak.ScreenWidth/2),
		float64(oak.ScreenHeight/2),
		16,
		32,
		render.NewSwitch("standRT", anims),
		nil,
		0,
		0,
	)

	s.Bind(func(id int, _ interface{}) int {
		p, ok := event.GetEntity(id).(*entities.Interactive)
		if !ok {
			dlog.Error("Non-player sent to player binding")
		}

		p.Delta.Zero()

		if oak.IsDown("W") {
			p.Delta.Add(physics.NewVector(0, -p.Speed.Y()))
		}
		if oak.IsDown("S") {
			p.Delta.Add(physics.NewVector(0, p.Speed.Y()))
		}
		if oak.IsDown("A") {
			p.Delta.Add(physics.NewVector(-p.Speed.X(), 0))
		}
		if oak.IsDown("D") {
			p.Delta.Add(physics.NewVector(p.Speed.X(), 0))
		}

		p.Vector.Add(p.Delta)
		// This is 6, when it should be 32
		//_, h := r.GetDims()
		hf := 32.0
		if p.Vector.Y() < float64(oak.ScreenHeight)*1/3 {
			p.Vector.SetY(float64(oak.ScreenHeight) * 1 / 3)
		} else if p.Vector.Y() > (float64(oak.ScreenHeight) - hf) {
			p.Vector.SetY((float64(oak.ScreenHeight) - hf))
		}
		if p.Vector.X() < 220 {
			p.Vector.SetX(220)
		} else if p.Vector.X()+p.W > float64(oak.ScreenWidth) {
			p.Vector.SetX(float64(oak.ScreenWidth) - p.W)
		}
		p.R.SetPos(p.Vector.X(), p.Vector.Y())
		p.RSpace.Update(p.Vector.X(), p.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
		<-s.RSpace.CallOnHits()
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
		return 0
	}, "EnterFrame")
	s.Speed = physics.NewVector(5, 5) // We actually allow players to move around in the inn!

	s.RSpace.Add(collision.Label(labels.Door),
		(func(s1, s2 *collision.Space) {
			nextscene = "run"
			stayInMenu = false
		}))

	render.Draw(s.R, 2, 1)
}
