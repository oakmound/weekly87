package enemies

import (
	"errors"
	"fmt"

	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/restrictor"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/physics"

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
	Speed      floatgeom.Point2
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

func (c *Constructor) Copy() *Constructor {
	c2 := &Constructor{
		Position:   c.Position,
		Dimensions: c.Dimensions,
		Speed:      c.Speed,
		// Todo: Assuming right now that the bindings map never gets modified (by a variant)
		Bindings:     c.Bindings,
		Health:       c.Health,
		AnimationMap: make(map[string]render.Modifiable, len(c.AnimationMap)),
	}
	for k, v := range c.AnimationMap {
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
	be.Speed = physics.NewVector(ec.Speed.X(), ec.Speed.Y())
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

		push := be.PushBack

		if be.facing == "RT" {
			push.Scale(-1)
		}
		be.Delta = be.Speed.Copy().Add(push)
		be.PushBack.Scale(0.8)
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

		fmt.Println("Trying to effect an enemy but for now we will be lazy and murder the enemy ", be)
		event.Trigger("EnemyDeath", []int64{secid, idx})
		be.Destroy()
	})
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
