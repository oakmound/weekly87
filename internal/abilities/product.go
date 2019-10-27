package abilities

import (
	"time"

	"github.com/oakmound/oak/shape"

	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/weekly87/internal/abilities/buff"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/layer"
	"github.com/oakmound/weekly87/internal/sfx"
)

//Producer of ability affects
type Producer struct {
	Start     floatgeom.Point2
	End       floatgeom.Point2
	ArcPoints []float64
	Frames    int

	W float64
	H float64

	FollowX   *float64
	FollowY   *float64
	Generator particle.Generator

	HitEffects map[collision.Label]collision.OnHit
	Label      collision.Label
	R          render.Renderable

	ToPlay string

	ThenFn    DoOption
	WhileFn   DoOption
	Interval  time.Duration
	TotalLife time.Duration

	Arc bool

	Buffs []buff.Buff
}

// Option to set on the producer
type Option func(Producer) Producer

// FrameLength overwrites the default of 100 for frame length with the provided int
func FrameLength(frames int) Option {
	return func(p Producer) Producer {
		p.Frames = frames
		return p
	}
}

// FollowSpeed to set on the producer. Modifies the base speed of the ability to keep pace weith these float pointers
func FollowSpeed(xFollow, yFollow *float64) Option {
	return func(p Producer) Producer {
		p.FollowX = xFollow
		p.FollowY = yFollow
		return p
	}
}

// StartAt location for the ability
func StartAt(pt floatgeom.Point2) Option {
	return func(p Producer) Producer {
		p.Start = pt
		return p
	}
}

// ArcTo denotes a path to set for the ability
func ArcTo(pts ...floatgeom.Point2) Option {
	return func(p Producer) Producer {
		p.ArcPoints = make([]float64, len(pts)*2)
		for i, pt := range pts {
			p.ArcPoints[i*2] = pt.X()
			p.ArcPoints[i*2+1] = pt.Y()
		}

		return p
	}
}

// LineTo denotes a simple point to go to for the ability
func LineTo(pt floatgeom.Point2) Option {
	return func(p Producer) Producer {
		p.End = pt
		return p
	}
}

// Duration sets the total life of the ability
func Duration(dur time.Duration) Option {
	return func(p Producer) Producer {
		p.TotalLife = dur
		return p
	}
}

// WithParticles sets a particle generator to couple with the ability
func WithParticles(pg particle.Generator) Option {
	return func(p Producer) Producer {
		p.Generator = pg
		return p
	}
}

// WithLabel sets the collision label  for the ability
func WithLabel(l collision.Label) Option {
	return func(p Producer) Producer {
		p.Label = l
		return p
	}
}

// WithHitEffects sets the hitmap on the ability
func WithHitEffects(he map[collision.Label]collision.OnHit) Option {
	return func(p Producer) Producer {
		p.HitEffects = he
		return p
	}
}

// WithRenderable sets the renderable for the ability
func WithRenderable(r render.Renderable) Option {
	return func(p Producer) Producer {
		p.R = r
		return p
	}
}

// DoOption is an option that will be performed at a given passed in location
// Often used for the DoAfter function
type DoOption func(floatgeom.Point2)

// Drop an ability effect at the passed in location
func Drop(p Producer) DoOption {
	return func(pt floatgeom.Point2) {
		dlog.Info("An ability dropped something")
		p.Start = pt
		chrs, err := p.Produce()
		if err != nil {
			dlog.Error(err)
			return
		}

		event.Trigger("AbilityFired", chrs)
	}
}

// DoPlay a sound as a DoOption
func DoPlay(s string) DoOption {
	return func(_ floatgeom.Point2) {
		sfx.Play(s)
	}
}

// PlaySFX on Produce() of the ability
func PlaySFX(s string) Option {
	return func(p Producer) Producer {
		p.ToPlay = s
		return p
	}
}

// Chain from where you left off and perform the action of the given producer
func Chain(p Producer) DoOption {
	return func(pt floatgeom.Point2) {
		dlog.Info("Chaining next ability")
		p.Start = p.Start.Add(pt)
		p.End = p.End.Add(pt)
		chrs, err := p.Produce()
		if err != nil {
			dlog.Error(err)
			return
		}

		event.Trigger("AbilityFired", chrs)
	}
}

// AndDo chains DoOptions
func AndDo(dos ...DoOption) DoOption {
	return func(pt floatgeom.Point2) {

		for _, o := range dos {
			o(pt)
		}
	}
}

// Then sets the action to take at the end of the producers life
func Then(do DoOption) Option {
	return func(p Producer) Producer {
		p.ThenFn = do
		return p
	}
}

// While a producer is alive do something
// Todo: implement while effects on product
func While(do DoOption, interval time.Duration) Option {
	return func(p Producer) Producer {
		p.WhileFn = do
		p.Interval = interval
		return p
	}
}

// WithBuff sets the buff on the ability
func WithBuff(b buff.Buff) Option {
	return func(p Producer) Producer {
		old := p.Buffs
		p.Buffs = make([]buff.Buff, len(old))
		copy(p.Buffs, old)
		p.Buffs = append(p.Buffs, b)
		return p
	}
}

// And concatenates options
func And(opts ...Option) Option {
	return func(p Producer) Producer {
		for _, o := range opts {
			p = o(p)
		}
		return p
	}
}

func defProducer() Producer {
	return Producer{
		W:      1,
		H:      1,
		Frames: 100,
	}
}

// Produce a set of outcomes given a set of options on the default producer
func Produce(opts ...Option) ([]characters.Character, error) {
	prd := defProducer()
	return prd.Produce(opts...)
}

// Produce a set of outcomes given a producer and a set of options
func (p Producer) Produce(opts ...Option) ([]characters.Character, error) {
	for _, o := range opts {
		p = o(p)
	}

	prd := &Product{
		Interactive: &entities.Interactive{},
		next:        p.ThenFn,
	}

	prd.Init()

	if p.Generator != nil {
		// Todo: what layer?
		particle.Layer(func(physics.Vector) int {
			return layer.Play
		})(p.Generator)
		prd.source = p.Generator.Generate(layer.Play)
	}

	// Todo: and label?
	if p.R != nil && p.W <= 1 && p.H <= 1 {
		w, h := p.R.GetDims()
		p.W = float64(w)
		p.H = float64(h)
	}

	prd.Interactive = entities.NewInteractive(
		p.Start.X(), p.Start.Y(),
		p.W, p.H,
		p.R, nil,
		prd.CID, 0,
	)
	prd.RSpace.Space.Label = p.Label
	for l, ef := range p.HitEffects {
		prd.RSpace.Add(l, ef)
	}

	if prd.R != nil {
		prd.R.SetPos(p.Start.X(), p.Start.Y())
		render.Draw(prd.R, layer.Effect)
	} else {
		prd.R = render.NewEmptySprite(0, 0, 1, 1) //Safety for Mover functionality
	}

	if prd.source != nil {
		x, y := p.Generator.GetPos()
		prd.source.SetPos(p.Start.X()+x, p.Start.Y()+y)
	}
	prd.FollowX = p.FollowX
	prd.FollowY = p.FollowY
	if prd.FollowX == nil {
		prd.FollowX = new(float64)
	}
	if prd.FollowY == nil {
		prd.FollowY = new(float64)
	}

	// If there's no end point, we shouldn't try to move the product
	if p.End != (floatgeom.Point2{}) || len(p.ArcPoints) > 0 {

		var curve shape.Bezier
		var err error
		if len(p.ArcPoints) > 0 {
			tempPoints := []float64{p.Start.X(), p.Start.Y()}
			tempPoints = append(tempPoints, p.ArcPoints...)
			curve, err = shape.BezierCurve(tempPoints...)
			if err != nil {
				dlog.Error("error making bezier curve", err)
				return nil, err
			}
		} else {
			curve, err = shape.BezierCurve(p.Start.X(), p.Start.Y(), p.End.X(), p.End.Y())
			if err != nil {
				dlog.Error("error making bezier curve", err)
				return nil, err
			}
		}
		positions := make([]floatgeom.Point2, p.Frames)
		delta := 1 / float64(p.Frames)
		j := 0
		for i := 0.0; j < len(positions); i += delta {
			x, y := curve.Pos(i)
			positions[j] = floatgeom.Point2{x, y}
			j++
		}

		deltas := make([]floatgeom.Point2, len(positions)-1)
		for i := 0; i < len(positions)-1; i++ {
			deltas[i] = positions[i+1].Sub(positions[i])
		}

		prd.Bind(func(id int, _ interface{}) int {
			prd, ok := event.GetEntity(id).(*Product)
			if !ok {
				dlog.Error("Non product sent to product enter frame")
				return 0
			}
			prd.position++
			if prd.position >= len(deltas) {
				prd.Destroy()
				return event.UnbindSingle
			}
			nextDelta := deltas[prd.position]
			prd.Interactive.ShiftPos(nextDelta.X()+*prd.FollowX, nextDelta.Y()+*prd.FollowY)
			if prd.source != nil {
				prd.source.ShiftX(nextDelta.X() + *prd.FollowX)
				prd.source.ShiftY(nextDelta.Y() + *prd.FollowY)
			}
			<-prd.Interactive.RSpace.CallOnHits()
			return 0
		}, "EnterFrame")
	}
	if p.TotalLife != 0 {
		endTime := time.Now().Add(p.TotalLife)
		prd.Bind(func(id int, _ interface{}) int {
			if time.Now().After(endTime) {
				prd.Destroy()
				return event.UnbindSingle
			}

			return 0
		}, "EnterFrame")
	}

	// This might expand later on if things have time limits
	if p.ThenFn == nil {
		prd.shouldPersist = true
	}

	prd.buffs = make([]buff.Buff, len(p.Buffs))
	copy(prd.buffs, p.Buffs)

	chrs := make([]characters.Character, 1)
	chrs[0] = prd

	if p.ToPlay != "" {
		sfx.Play(p.ToPlay)
	}

	return chrs, nil
}

// Init a product by getting it a CID
func (p *Product) Init() event.CID {
	p.CID = event.NextID(p)
	return p.CID
}

//Product of a ability producer
type Product struct {
	*entities.Interactive
	shouldPersist bool
	position      int
	TotalLife     time.Duration
	FollowX       *float64
	FollowY       *float64

	source *particle.Source
	next   func(floatgeom.Point2)
	buffs  []buff.Buff
}

// MoveParticles updates the location of the particle source on a product if it exists
func (p *Product) MoveParticles(nextDelta floatgeom.Point2) {
	if p.source != nil {
		p.source.ShiftX(nextDelta.X())
		p.source.ShiftY(nextDelta.Y())
	}
}

// Destroy cleans up a product
func (p *Product) Destroy() {
	// Note: this assumes that destroys aren't happening simultaneously
	if p.next != nil {
		p.next(floatgeom.Point2{p.X(), p.Y()})
		p.next = nil
	}
	p.Interactive.Destroy()
	if p.source != nil {
		p.source.Stop()
		p.source = nil
	}
}

// Activate a product, needed to fulfill some interfaces
func (p *Product) Activate() {}

// ShouldPersist to our records
func (p *Product) ShouldPersist() bool {
	return p.shouldPersist
}

// Buffs that the product gives
func (p *Product) Buffs() []buff.Buff {
	return p.buffs
}
