package inn

import (
	"math/rand"
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
)

type aiStatus int

const (
	aiContinue aiStatus = iota
	aiComplete aiStatus = iota
)

// AI can be added to NPCs to dictate their behaviour
type AI struct {
	inAction      bool
	curAction     func(int) aiStatus
	curCancel     func(int)
	actions       []aiAction
	actionWeights []float64
}

// NewAI creates a set of aiactions and their weights that can be acted upon to perform desired actions
func NewAI(actions []aiAction, weights []float64) *AI {
	remainingWeights := alg.RemainingWeights(weights)
	return &AI{
		actions:       actions,
		actionWeights: remainingWeights,
	}
}

type aiAction interface {
	start() (func(int) aiStatus, func(int))
}

// Choose what action the ai should take
func (ai *AI) Choose() aiAction {
	i := alg.WeightedChooseOne(ai.actionWeights)
	return ai.actions[i]
}

type aiWalkUpDownBar struct {
	top, bottom float64
	speed       floatrange.Range
	// in milliseconds
	duration intrange.Range
}

func (a aiWalkUpDownBar) start() (func(int) aiStatus, func(int)) {
	start := time.Now()
	end := start.Add(time.Duration(a.duration.Poll()) * time.Millisecond)
	downSpeed := physics.NewVector(0, a.speed.Poll())
	upSpeed := downSpeed.Copy().Scale(-1)
	speed := upSpeed
	return func(id int) aiStatus {
			keeper, ok := event.GetEntity(id).(*NPC)
			if !ok {
				dlog.Error("Got non NPC in Innkeeper bindings")
				return 1
			}
			if time.Now().After(end) {
				return aiComplete
			}
			if keeper.Y() > a.bottom {
				speed = upSpeed
			} else if keeper.Y() < a.top {
				speed = downSpeed
			}
			keeper.Delta.SetPos(speed.X(), speed.Y())

			keeper.ShiftPos(keeper.Delta.X(), keeper.Delta.Y())

			return aiContinue
		}, func(id int) {
			keeper, ok := event.GetEntity(id).(*NPC)
			if !ok {
				dlog.Error("Got non NPC in Innkeeper bindings")
				return
			}
			keeper.Delta.SetPos(0, 0)
		}
}

type aiServeDrinkLocation struct {
	rect      floatgeom.Rect2
	drinkImg  *render.Sprite
	sOverride *string
}

// start trying to serve a drink
// Navigate to a location and throw a drink at it
func (a aiServeDrinkLocation) start() (func(int) aiStatus, func(int)) {
	// where to serve drink
	barX := (rand.Float64() * (a.rect.Max.X() - a.rect.Min.X())) + a.rect.Min.X()
	barY := (rand.Float64() * (a.rect.Max.Y() - a.rect.Min.Y())) + a.rect.Min.Y()
	var drink *doodads.Consumable
	var dDelta physics.Vector

	reached := false
	return func(id int) aiStatus {

			keeper, ok := event.GetEntity(id).(*NPC)
			if !ok {
				dlog.Error("Got non NPC in Innkeeper bindings")
				return 1
			}

			// if gotten to serving location TODO: actually make pretty
			if reached {
				// CONSIDER: should innkeeper be responsible for this movement?
				drink.ShiftPos(dDelta.X(), dDelta.Y())
				if drink.X()-barX < 1 { //complete serving within acceptable bounds
					return aiContinue
				}
				drink.Activate()
				return aiComplete
			}

			dy := 1.0
			if keeper.Y() > barY {
				dy = -1.0
			}

			keeper.Delta.SetPos(0, dy)
			keeper.ShiftPos(0, dy)
			if keeper.Y() > barY-1 && keeper.Y() < barY+1 {
				reached = true
				*a.sOverride = "consume"
				drinkStartY := keeper.Y() - 16 // Unintentionally funny name!
				drink = doodads.NewDrinkable(keeper.X(), drinkStartY, a.drinkImg)
				dDelta = physics.NewVector(barX-keeper.X(), barY-drinkStartY).Normalize()

			}

			return aiContinue
		}, func(id int) {
			*a.sOverride = ""
		}
}

type aiDrinker struct {
	duration  intrange.Range
	solid     *entities.Interactive
	consuming int
}

// start aidrinkerbindings
// currently unsafe and could end up with people in bad states
func (a aiDrinker) start() (func(int) aiStatus, func(int)) {
	start := time.Now()
	nextCheck := start.Add(time.Duration(a.duration.Poll()) * time.Millisecond)

	pSpace := a.solid.GetReactiveSpace()
	// x, y := pSpace.GetPos()
	drinkSpace := collision.NewEmptyReactiveSpace(
		collision.NewUnassignedSpace(pSpace.X()-50, pSpace.Y(), pSpace.GetW()+50, pSpace.GetH()))

	swtch := a.solid.R.(*render.Switch)

	a.solid.CID.Bind(func(id int, _ interface{}) int {

		// TODO: reconsider structure here and if we want mutexes
		a.consuming--
		if a.consuming <= 0 {
			swtch.Set("standLT") // hacky as we only support bar drinkers for now
		}

		return 0
	}, "consumeCompleted")

	drinkSpace.Add(labels.Drinkable, func(_, d *collision.Space) {
		d.CID.Trigger("Consume", pSpace.Space)
		swtch.Set("consume")
		a.consuming++

		dlog.Info("Trying to consume a drink")
	})
	a.solid.Tree.Add(drinkSpace.Space)

	return func(id int) aiStatus {
			if time.Now().After(nextCheck) {
				nextCheck = nextCheck.Add(time.Duration(a.duration.Poll()) * time.Millisecond)

				<-drinkSpace.CallOnHits()
			}
			return aiContinue
		}, func(id int) {
			drinker, ok := event.GetEntity(id).(*NPC)
			if !ok {
				dlog.Error("Got non NPC in drinker bindings")
				return
			}
			drinker.Delta.SetPos(0, 0)
		}
}

type aiIdler struct {
	duration intrange.Range
}

func (a aiIdler) start() (func(int) aiStatus, func(int)) {
	start := time.Now()
	end := start.Add(time.Duration(a.duration.Poll()) * time.Millisecond)
	return func(id int) aiStatus {
			if time.Now().After(end) {
				return aiComplete
			}
			return aiContinue
		}, func(id int) {

		}
}
