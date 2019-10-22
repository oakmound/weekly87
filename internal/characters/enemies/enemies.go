package enemies

import (
	"errors"
	"fmt"
	"time"

	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/restrictor"
	"github.com/oakmound/weekly87/internal/vfx"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/particle"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

var requiredAnimations = []string{
	"standRT",
	"standLT",
	"walkRT",
	"walkLT",
}

type Constructor struct {
	Position   floatgeom.Point2
	Dimensions floatgeom.Point2

	Speed       floatgeom.Point2
	SpaceOffset physics.Vector

	// The following strings are required in the animation map:
	// "standRT"
	// "standLT"
	// "walkRT"
	// "walkLT"
	// more may be added
	AnimationMap map[string]render.Modifiable
	Bindings     map[string]func(*BasicEnemy, interface{}) int
	Health       int
}

// Copy the data values to new instances of an enemy constructor
func (ec *Constructor) Copy() *Constructor {
	c2 := &Constructor{
		Position:    ec.Position,
		Dimensions:  ec.Dimensions,
		SpaceOffset: ec.SpaceOffset,
		Speed:       ec.Speed,
		// Todo: Assuming right now that the bindings map never gets modified (by a variant)
		Bindings:     ec.Bindings,
		Health:       ec.Health,
		AnimationMap: make(map[string]render.Modifiable, len(ec.AnimationMap)),
	}
	for k, v := range ec.AnimationMap {
		c2.AnimationMap[k] = v.Copy()
	}
	return c2
}

var Constructors [TypeLimit * VariantCount]*Constructor

func setConstructor(eType, size, color int, cons *Constructor) {
	Constructors[(eType*VariantCount)+(size*lastColor)+color] = cons
}

func GetConstructor(eType, size, color int) *Constructor {
	return Constructors[(eType*VariantCount)+(size*lastColor)+color]
}

type BasicEnemy struct {
	*entities.Interactive
	facing        string
	swtch         *render.Switch
	Active        bool
	beenDisplayed bool
	PushBack      physics.Vector
	baseSpeed     physics.Vector
	Health        int
}

func (be *BasicEnemy) Init() event.CID {
	return event.NextID(be)
}

func (be *BasicEnemy) Activate() {
	be.Active = true
	restrictor.Add(be)
}

func (be *BasicEnemy) Destroy() {
	be.Active = false
	be.Interactive.Destroy()
}

func (be *BasicEnemy) DeathEffect(secid, idx int64) {
	be.RSpace.Label = 0
	be.PushBack.Add(physics.NewVector(60, 0))
	abilities.Produce(
		abilities.StartAt(floatgeom.Point2{be.X() + 8, be.Y() + 10}),
		abilities.LineTo(floatgeom.Point2{be.X() + 30, be.Y() + 30}),
		//abilities.FollowSpeed(ply.Delta.Xp(), ply.Delta.Yp()),
		abilities.WithParticles(vfx.WhiteSlash()),
		abilities.FrameLength(20),
	)
	be.CheckedBind(func(be *BasicEnemy, data interface{}) int {

		// wait for pushback to complete
		if be.PushBack.Magnitude() > 0.15 {
			return 0
		}
		event.Trigger("EnemyDeath", []int64{secid, idx})

		w, h := be.R.GetDims()

		// Create a visual effect, overwrite?
		source := vfx.RedRing().Generate(2)
		source.SetPos(be.X()+float64(w/2), be.Y()+float64(h/2))
		endSource := time.Now().Add(time.Millisecond * 30)
		source.CID.Bind(func(id int, data interface{}) int {
			eff, ok := event.GetEntity(id).(*particle.Source)
			if ok {
				eff.ShiftX(be.Delta.X() + 1)

				if endSource.Before(time.Now()) {
					eff.Stop()
					return 1
				}
			}
			return 0
		}, "EnterFrame")

		be.Destroy()
		return event.UnbindSingle
	}, "EnterFrame")

}

func (be *BasicEnemy) CheckedBind(bnd func(*BasicEnemy, interface{}) int, ev string) {
	be.Bind(func(id int, data interface{}) int {
		be, ok := event.GetEntity(id).(*BasicEnemy)
		if !ok {
			dlog.Error("Basic Enemy binding was called on non-basic enemy")
			return event.UnbindSingle
		}
		return bnd(be, data)
	}, ev)
}

// NewEnemy creates an enemy that will animate walking or standing appropriately,
// move according to its speed, flip its facing when the player picks up
// a chest, and die when a player attack hits it
func (ec *Constructor) NewEnemy(secid, idx int64) (*BasicEnemy, error) {
	if ec.Dimensions == (floatgeom.Point2{}) {
		return nil, errors.New("Dimensions must be provided")
	}
	for _, s := range requiredAnimations {
		if _, ok := ec.AnimationMap[s]; !ok {
			return nil, errors.New("Animation name " + s + " must be provided")
		}
	}
	be := &BasicEnemy{}
	be.PushBack = physics.NewVector(0, 0)
	newMp := map[string]render.Modifiable{}
	for animKey, anim := range ec.AnimationMap {
		newMp[animKey] = anim.Copy()
	}
	be.swtch = render.NewSwitch("standLT", newMp)
	if ec.SpaceOffset != (physics.Vector{}) {

		for animKey := range ec.AnimationMap {
			be.swtch.SetOffsets(animKey, ec.SpaceOffset)
		}
	}
	be.Interactive = entities.NewInteractive(
		ec.Position.X(),
		ec.Position.Y(),
		ec.Dimensions.X(),
		ec.Dimensions.Y(),
		be.swtch,
		nil,
		be.Init(),
		0,
	)
	// be.swtch.SetOffsets("walkLT", )
	be.Health = ec.Health
	be.Speed = physics.NewVector(ec.Speed.X(), ec.Speed.Y())
	be.baseSpeed = be.Speed.Copy()
	be.facing = "LT"
	be.RSpace.Label = labels.Enemy
	be.CheckedBind(func(be *BasicEnemy, _ interface{}) int {
		be.facing = "RT"
		be.Speed = be.Speed.Scale(-1)
		return 0
	}, "RunBack")
	be.CheckedBind(func(be *BasicEnemy, _ interface{}) int {
		// Enemies should only do anything if they are on screen
		// Todo: other things could effect delta temporarily

		push := be.PushBack.Copy()

		if be.facing == "RT" {
			push.Scale(-1)
		}
		be.Delta = be.Speed.Copy().Add(push)
		be.PushBack.Scale(0.86)
		if be.X() <= float64(oak.ScreenWidth+oak.ViewPos.X) &&
			be.X()+be.W >= float64(oak.ViewPos.X) {
			//be.RSpace.Label = labels.Enemy
			be.ShiftPos(be.Delta.X(), be.Delta.Y())
			// Default behavior is to flip when hitting the ceiling
			if be.Y() < float64(oak.ScreenHeight)*1/3 ||
				be.Y() > (float64(oak.ScreenHeight)-be.H) {
				be.Speed.SetY(be.Speed.Y() * -1)
				// Adjust so we don't exist in the wall for a frame
				be.ShiftPos(0, be.Speed.Y())
			}
		}
		if be.Delta.X() != 0 || be.Delta.Y() != 0 {
			be.swtch.Set("walk" + be.facing)
		} else {
			be.swtch.Set("stand" + be.facing)
		}
		<-be.RSpace.CallOnHits()
		return 0
	}, "EnterFrame")

	be.GetReactiveSpace().Add(labels.EffectsEnemy, func(s, bf *collision.Space) {
		be, ok := s.CID.E().(*BasicEnemy)
		if !ok {
			dlog.Error("Non-enemy affected??")
			fmt.Printf("%T\n", s.CID.E())
			return
		}

		fmt.Println("Consider moving this effect to trigger vie the attacked event", be)

		be.DeathEffect(secid, idx)
	})
	be.CheckedBind(func(be *BasicEnemy, data interface{}) int {

		effectMap, ok := data.(map[string]float64)
		if !ok {
			dlog.Warn("Data sent on attack was not in the right format")
			return 0
		}

		for k, v := range effectMap {
			switch k {
			case "pushback":
				be.PushBack.Add(physics.NewVector(v, 0))
			case "damage":
				be.Health -= int(v)
				if be.Health < 1 {
					event.Trigger("EnemyDeath", []int64{secid, idx})
					be.Destroy()
				}
			case "frost":
				endDebuff := time.Now().Add(time.Second * 3)
				be.CheckedBind(func(be *BasicEnemy, data interface{}) int {
					if !time.Now().After(endDebuff) {
						return 0
					}
					be.Speed = be.Speed.Scale(v)
					return event.UnbindSingle
				}, "EnterFrame")
				be.Speed = be.Speed.Scale(1 / v)
				dlog.Verb("BE speed is now", be.Speed)
			}
		}

		return 0
	}, "Attacked")
	// be.RSpace.Add(labels.PlayerAttack, func(s, _ *collision.Space) {
	// 	be, ok := s.CID.E().(*BasicEnemy)
	// 	if !ok {
	// 		dlog.Error("On hit for basic enemy called on non-basic enemy")
	// 		return
	// 	}
	// 	// TODO: track changes?
	// 	event.Trigger("EnemyDeath", []int64{secid, idx})
	// 	be.Destroy()
	// })
	for ev, b := range ec.Bindings {
		be.CheckedBind(b, ev)
	}
	return be, nil
}

func (be *BasicEnemy) RunBackwards() {
	be.facing = "RT"
	be.Speed = be.Speed.Scale(-1)
}

func (be *BasicEnemy) GetDims() (int, int) {
	return be.swtch.GetDims()
}
