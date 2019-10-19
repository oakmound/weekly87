package inn

import (
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak"
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

const (
	prePlay = iota

	inMenu
	playing
)

type innWalker struct {
	*players.FreeWalker
}

// setParty for the innWalker. Create any nessecary players and update if not
// Is responsible for inital location!
func (iw *innWalker) setParty(plys []*players.Player) {
	if len(plys) == 0 {
		dlog.Error("Need at least one party member")
		return
	}

	if iw.Front == nil {
		iw.Front = entities.NewInteractive(
			float64(oak.ScreenWidth/2),
			float64(170),
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

// newInnWalker creates a special character for the inn
func newInnWalker(scale float64, plys []*players.Player) *innWalker {

	iw := &innWalker{
		&players.FreeWalker{Scale: scale},
	}
	iw.setParty(plys)

	return iw
}

func (iw *innWalker) bindFront() {

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
		default:
			lowestID := joys.LowestID()
			js := joys.StickState(lowestID)
			if oak.IsDown(key.UpArrow) || js.StickLY > 8000 ||
				oak.IsDown(key.DownArrow) || js.StickLY < -8000 ||
				oak.IsDown(key.LeftArrow) || js.StickLX < -8000 ||
				oak.IsDown(key.RightArrow) || js.StickLX > 8000 {
				iw.State = playing
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

	iw.Front.RSpace.Add(collision.Label(labels.Door), (func(s1, s2 *collision.Space) {
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
	n.AI = NewAI(
		[]aiAction{
			aiDrinker{
				duration: intrange.NewLinear(2000, 4000),
				solid:    n.Interactive,
			},
		},
		[]float64{
			1,
		},
	)

	n.Bind(func(id int, frame interface{}) int {

		drinkr, ok := event.GetEntity(id).(*NPC)
		if !ok {
			dlog.Error("Got non NPC in Drinker bindings")
			return 1
		}

		if drinkr.inAction {
			status := drinkr.curAction(id)
			if status == aiComplete {
				drinkr.curCancel(id)
				drinkr.inAction = false
			}
			return 0
		}

		action := drinkr.AI.Choose()
		drinkr.curAction, drinkr.curCancel = action.start()
		drinkr.inAction = true

		return 0
	}, "EnterFrame")

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
func NewInnkeeper(img *render.Sprite, scale, x, y float64) *NPC {
	keeper := NewInnNPC(players.InnKeeper, scale, x, y)
	keeper.AI = NewAI(
		[]aiAction{
			aiWalkUpDownBar{
				top:      300,
				bottom:   700,
				speed:    floatrange.NewConstant(1),
				duration: intrange.NewConstant(2000),
			},
			aiServeDrinkLocation{
				rect:     floatgeom.NewRect2WH(180, 150, 25, float64(oak.ScreenHeight)/2),
				drinkImg: img,
			},
		},
		[]float64{
			.5, .5,
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
