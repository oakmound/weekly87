package inn

import (
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/joys"
	"github.com/oakmound/weekly87/internal/layer"
)

const frameLag = 10
const maxPartySize = 4

const (
	prePlay = iota

	inMenu
	playing
)

type innWalker struct {
	front     *entities.Interactive
	followers []*entities.Interactive
	scale     float64
	lagDeltas [frameLag * maxPartySize]floatgeom.Point2
	lagIdx    int
	gameState int
}

// setParty for the innWalker. Create any nessecary players and update if not
// Is responsible for inital location!
func (iw *innWalker) setParty(plys []*players.Player) {
	if len(plys) == 0 {
		dlog.Error("Need at least one party member")
		return
	}

	if iw.front == nil {
		iw.front = entities.NewInteractive(
			float64(oak.ScreenWidth/2),
			float64(170),
			16*iw.scale,
			32*iw.scale,
			plys[0].Swtch.Copy().Modify(mod.Scale(iw.scale, iw.scale)),
			nil,
			0,
			0,
		)
		iw.bindFront()
	} else {
		old := iw.front.R
		old.Undraw()
		iw.front.R = plys[0].Swtch.Copy().Modify(mod.Scale(iw.scale, iw.scale))
		iw.front.R.SetPos(old.X(), old.Y())
		iw.front.R.(*render.Switch).Set(old.(*render.Switch).Get())
	}
	render.Draw(iw.front.R, layer.Play, maxPartySize)

	for i := 1; i < len(plys); i++ {
		if i >= len(iw.followers) {
			// make a new one for this position
			iw.followers = append(iw.followers, entities.NewInteractive(
				iw.front.X(),
				iw.front.Y(),
				16*iw.scale,
				32*iw.scale,
				plys[i].Swtch.Copy().Modify(mod.Scale(iw.scale, iw.scale)),
				nil,
				0,
				0,
			))

		} else {
			old := iw.followers[i-1].R
			old.Undraw()
			iw.followers[i-1].R = plys[i].Swtch.Copy().Modify(mod.Scale(iw.scale, iw.scale))
			iw.followers[i-1].SetPos(old.X(), old.Y())
			iw.followers[i-1].R.(*render.Switch).Set(old.(*render.Switch).Get())

		}
		render.Draw(iw.followers[i-1].R, layer.Play, maxPartySize-1)
	}
}

// newInnWalker creates a special character for the inn
func newInnWalker(scale float64, plys []*players.Player) *innWalker {

	iw := &innWalker{
		scale: scale,
	}
	iw.setParty(plys)

	return iw
}

func (iw *innWalker) bindFront() {
	lowestID := joys.LowestID()
	iw.front.Bind(func(id int, _ interface{}) int {
		p, ok := event.GetEntity(id).(*entities.Interactive)
		if !ok {
			dlog.Error("Non-player sent to player binding")
		}

		switch iw.gameState {
		case inMenu:
			p.Delta.Zero()
			return 0
		case playing:
			p.Delta.Zero()
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
		p.Vector.Add(p.Delta)

		_, h := p.R.GetDims()
		hf := float64(h)
		// if p.Vector.Y() < float64(oak.ScreenHeight)*1/3 {
		// 	p.Vector.SetY(float64(oak.ScreenHeight) * 1 / 3)
		// }
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
		<-iw.front.RSpace.CallOnHits()
		if collision.HitLabel(iw.front.RSpace.Space, labels.Blocking, labels.NPC) != nil {
			p.Vector.Sub(p.Delta)
			p.Delta.Zero()
			p.R.SetPos(p.Vector.X(), p.Vector.Y())
			p.RSpace.Update(p.Vector.X(), p.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
		}

		if p.Delta.Magnitude() != 0 {
			//fmt.Println("delta x, y", p.Delta.X(), p.Delta.Y())
			//fmt.Println("lag index", iw.lagIdx)
			// Store this event in our frame delay
			iw.lagDeltas[iw.lagIdx] = floatgeom.Point2{p.Delta.X(), p.Delta.Y()}

			// Access stored frame deltas, move followers
			for i, fw := range iw.followers {
				delta := iw.lagDeltas[(iw.lagIdx+frameLag*(i+1))%len(iw.lagDeltas)]
				fw.Vector.Add(physics.NewVector(delta.X(), delta.Y()))
				fw.R.SetPos(fw.Vector.X(), fw.Vector.Y())

				swch := fw.R.(*render.Switch)
				if delta.X() > 0 {
					swch.Set("walkRT")
				} else {
					swch.Set("walkLT")
				}
			}

			// Shift the contents of the frame deltas
			iw.lagIdx--
			if iw.lagIdx < 0 {
				iw.lagIdx = len(iw.lagDeltas) - 1
			}
		} else {
			for _, fw := range iw.followers {
				swch := fw.R.(*render.Switch)
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
		return 0
	}, "EnterFrame")
	iw.front.Speed = physics.NewVector(5, 5) // We actually allow players to move around in the inn!

	iw.front.RSpace.Add(collision.Label(labels.Door), (func(s1, s2 *collision.Space) {
		d, ok := s2.CID.E().(*doodads.InnDoor)
		if !ok {
			dlog.Error("Non-door sent to inndoor binding")
			return
		}
		nextscene = d.NextScene
		stayInMenu = false
	}))
}

// NPC is a inn only construct
type NPC struct {
	*entities.Interactive
	Swtch          *render.Switch
	Class          int
	Button         render.Renderable
	UndrawButtonAt time.Time
	*AI
}

// Init the npc so it has a CID!
func (n *NPC) Init() event.CID {
	return event.NextID(n)
}

// FaceLeft take wether the npc should face left
// Set the facing renderable appropriately
func (n *NPC) FaceLeft(shouldFaceLeft bool) *NPC {
	if shouldFaceLeft {
		n.R.(*render.Switch).Set("standLT")
	} else {
		n.R.(*render.Switch).Set("standRT")
	}
	return n
}

// NewInnNPC creates a npc to interact with for setting up party
func NewInnNPC(class int, scale, x, y float64) *NPC {
	pcon := players.ClassConstructor([]players.PartyMember{{class, 0, "NPC How did you find me"}})[0]
	n := &NPC{}
	n.Class = class
	n.Swtch = render.NewSwitch("standRT", pcon.AnimationMap).Copy().(*render.Switch)
	n.Swtch.Modify(mod.Scale(scale, scale))
	n.Interactive = entities.NewInteractive(
		x,
		y,
		pcon.Dimensions.X()*scale,
		pcon.Dimensions.Y()*scale,
		n.Swtch,
		nil,
		n.Init(),
		0,
	)

	return n
}

// Activate draws the npc and makes it collidable
func (n NPC) Activate() {
	n.RSpace.UpdateLabel(labels.NPC)
	render.Draw(n.R, layer.Play, 1)
}

// Destroy removes the npc from the collision tree
func (n NPC) Destroy() {
	n.Tree.Remove(n.RSpace.Space)
	n.Interactive.Destroy()
}

// NewInnkeeper creates the innkeeper who moves around behind the bar
func NewInnkeeper(scale, x, y float64) *NPC {
	keeper := NewInnNPC(players.InnKeeper, scale, x, y)
	keeper.AI = NewAI(
		[]aiAction{
			aiWalkUpDownBar{
				top:      300,
				bottom:   700,
				speed:    floatrange.NewConstant(1),
				duration: intrange.NewConstant(2000),
			},
		},
		[]float64{
			1,
		},
	)
	keeper.Bind(func(id int, frame interface{}) int {

		kpr, ok := event.GetEntity(id).(*NPC)
		if !ok {
			dlog.Error("Got non NPC in Innkeeper bindings")
			return 1
		}

		if kpr.inAction {
			status := kpr.curAction(id)
			if status == aiComplete {
				kpr.curCancel(id)
				kpr.inAction = false
			}
			return 0
		}

		action := keeper.AI.Choose()
		kpr.curAction, kpr.curCancel = action.start()
		kpr.inAction = true

		return 0
	}, "EnterFrame")
	return keeper
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

			// Todo: pull out into "move"?
			keeper.ShiftPos(keeper.Delta.X(), keeper.Delta.Y())
			if keeper.Delta.X() != 0 || keeper.Delta.Y() != 0 {
				if keeper.Delta.X() < 0 {
					keeper.Swtch.Set("walkLT")
				} else {
					keeper.Swtch.Set("walkRT")
				}
			} else {
				cur := keeper.Swtch.Get()
				err := keeper.Swtch.Set("stand" + string(cur[len(cur)-2:]))
				dlog.ErrorCheck(err)
			}

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

type aiStatus int

const (
	aiContinue aiStatus = iota
	aiComplete aiStatus = iota
)

type AI struct {
	inAction      bool
	curAction     func(int) aiStatus
	curCancel     func(int)
	actions       []aiAction
	actionWeights []float64
}

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

func (ai *AI) Choose() aiAction {
	i := alg.WeightedChooseOne(ai.actionWeights)
	return ai.actions[i]
}
