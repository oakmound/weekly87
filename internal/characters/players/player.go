package players

import (
	"sort"
	"sync"
	"time"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/abilities/buff"

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
	}
}

type Player struct {
	*entities.Interactive
	facing      string
	Swtch       *render.Switch
	Special1    abilities.Ability
	Special2    abilities.Ability
	Alive       bool
	RunSpeed    float64
	ChestValues []int64
	Chests      []render.Renderable
	BuffLock    sync.Mutex
	Buffs       []buff.Buff
	*buff.Status
	Party *Party
}

func (p *Player) AddBuff(b buff.Buff) {
	p.BuffLock.Lock()
	b.ExpireAt = time.Now().Add(b.Duration)
	p.Buffs = append(p.Buffs, b)
	sort.Slice(p.Buffs, func(i, j int) bool {
		return p.Buffs[i].ExpireAt.Before(p.Buffs[j].ExpireAt)
	})
	p.BuffLock.Unlock()
	b.Enable(p.Status)
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

const WallOffset = 50
