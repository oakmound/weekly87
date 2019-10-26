package players

import (
	"sort"
	"sync"
	"time"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/abilities/buff"
	"github.com/oakmound/weekly87/internal/layer"

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
	"deadRT",
	"deadLT",
	"walkHold",
	"standHold",
}

type Constructor struct {
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
	Special1     abilities.Ability
	Special2     abilities.Ability
	RunSpeed     float64
	Name         string
	AccruedValue int
}

// Copy returns a shallow copy of the constructor.
// Don't modify the animation map or bindings on a copy of a constructor, that part isn't
// deep copied.
func (pc *Constructor) Copy() *Constructor {
	return &Constructor{
		Position:     pc.Position,
		Speed:        pc.Speed,
		Dimensions:   pc.Dimensions,
		AnimationMap: pc.AnimationMap,
		Bindings:     pc.Bindings,
		Special1:     pc.Special1,
		Special2:     pc.Special2,
		RunSpeed:     pc.RunSpeed,
		Name:         pc.Name,
		AccruedValue: pc.AccruedValue,
	}
}

// Player contains all the information for the pcs
type Player struct {
	*entities.Interactive
	facing       string
	Name         string
	AccruedValue int
	Swtch        *render.Switch
	Special1     abilities.Ability
	Special2     abilities.Ability
	Alive        bool
	RunSpeed     float64
	PartyIndex   int
	ChestValues  []int64
	Chests       []render.Renderable
	ChestsHeight float64
	buffR        []render.Renderable
	BuffLock     sync.Mutex
	Buffs        []buff.Buff
	*buff.Status
	Party *Party
}

func (p *Player) GetDelta() physics.Vector {
	return p.Party.Players[0].Delta
}

func (p *Player) AddBuff(b buff.Buff) {
	if !p.Alive {
		return
	}
	p.BuffLock.Lock()
	b.ExpireAt = time.Now().Add(b.Duration)
	b.R = buff.BasicBuffSwitch(b.RGen())
	p.Buffs = append(p.Buffs, b)
	render.Draw(b.R, layer.UI, 10)
	sort.Slice(p.Buffs, func(i, j int) bool {
		return p.Buffs[i].ExpireAt.Before(p.Buffs[j].ExpireAt)
	})
	p.BuffLock.Unlock()
	b.Enable(p.Status)
	p.ReorderBuffs()

}

// DropChest if the player has one
func (p *Player) DropChest() {
	if len(p.ChestValues) == 0 {
		return
	}

	_, h := p.Chests[len(p.Chests)].GetDims()
	p.ChestsHeight -= float64(h)
	p.ChestValues = p.ChestValues[:len(p.ChestValues)-1]
	p.Chests = p.Chests[:len(p.Chests)-1]

	if len(p.ChestValues) > 0 {
		return
	}
	p.Swtch.Set("walkLT")
	p.Special1.Enable(true)
	p.Special2.Enable(true)
}

func (p *Player) AddChest(h int, r render.Modifiable, contents int64) {
	p.ChestsHeight += float64(h)
	chestHeight := p.ChestsHeight

	r.(*render.Sprite).Vector = r.Attach(p.Vector, -3, -chestHeight)
	p.ChestValues = append(p.ChestValues, contents)
	p.Chests = append(p.Chests, r)
	render.Draw(r, layer.Play, 2)

	if len(p.ChestValues) == 1 {
		p.Swtch.Set("walkHold")
		p.Special1.Enable(false)
		p.Special2.Enable(false)
	}
}

func (p *Player) ReorderBuffs() {
	xOffset := abilities.BuffIconSize + 4
	yOffset := abilities.BuffIconSize + 4
	x := float64(p.PartyIndex*xOffset + oak.ScreenWidth/4*3)

	for i, b := range p.Buffs {
		b.R.SetPos(x, float64(yOffset*i+16))
	}
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

func (p *Player) Direction() string {
	return p.facing
}

func (p *Player) Ready() bool {
	return p.Alive
}

func (p *Player) Kill() {
	p.Special1.Enable(false)
	p.Special2.Enable(false)
	dlog.ErrorCheck(p.Swtch.Set("dead" + p.facing))
	for _, r := range p.Chests {
		r.Undraw()
	}
	p.ChestValues = []int64{}
	p.ChestsHeight = 0
	p.Chests = []render.Renderable{}
	p.BuffLock.Lock()
	for _, b := range p.Buffs {
		b.Disable(p.Status)
		b.R.Undraw()
	}
	p.Buffs = []buff.Buff{}
	p.BuffLock.Unlock()
	// Consider: Drop the chests?
	p.Alive = false
}

func (p *Player) Revive() {
	p.Alive = true
	p.Special1.Enable(true)
	p.Special2.Enable(true)
	dlog.ErrorCheck(p.Swtch.Set("walk" + p.facing))
}

const WallOffset = 50
