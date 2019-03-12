package inn

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/entities/x/move"
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
		ply, ok := event.GetEntity(id).(*entities.Interactive)
		if !ok {
			dlog.Error("Non-player sent to player binding")
		}

		move.WASD(ply)
		move.Limit(ply, innSpace)
		<-s.RSpace.CallOnHits()
		swch := ply.R.(*render.Switch)
		if ply.Delta.X() != 0 || ply.Delta.Y() != 0 {
			if ply.Delta.X() > 0 {
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
