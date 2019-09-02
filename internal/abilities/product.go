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
)

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

	Label collision.Label
	R     render.Renderable

	ThenFn   DoOption
	WhileFn  DoOption
	Interval time.Duration

	Arc bool

	Buffs []buff.Buff
}

type Option func(Producer) Producer

func FrameLength(frames int) Option {
	return func(p Producer) Producer {
		p.Frames = frames
		return p
	}
}

func FollowSpeed(xFollow, yFollow *float64) Option {

	return func(p Producer) Producer {
		p.FollowX = xFollow
		p.FollowY = yFollow
		return p
	}
}

func StartAt(pt floatgeom.Point2) Option {
	return func(p Producer) Producer {
		p.Start = pt
		return p
	}
}

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

func LineTo(pt floatgeom.Point2) Option {
	return func(p Producer) Producer {
		p.End = pt

		return p
	}
}

func WithParticles(pg particle.Generator) Option {
	return func(p Producer) Producer {
		p.Generator = pg
		return p
	}
}

func WithLabel(l collision.Label) Option {
	return func(p Producer) Producer {
		p.Label = l
		return p
	}
}

func WithRenderable(r render.Renderable) Option {
	return func(p Producer) Producer {
		p.R = r
		return p
	}
}

type DoOption func(floatgeom.Point2)

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

func Perform(p Producer) DoOption {
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

func Then(do DoOption) Option {
	return func(p Producer) Producer {
		p.ThenFn = do
		return p
	}
}

// Todo: implement while effects on product
func While(do DoOption, interval time.Duration) Option {
	return func(p Producer) Producer {
		p.WhileFn = do
		p.Interval = interval
		return p
	}
}

func WithBuff(b buff.Buff) Option {
	return func(p Producer) Producer {
		old := p.Buffs
		p.Buffs = make([]buff.Buff, len(old))
		copy(p.Buffs, old)
		p.Buffs = append(p.Buffs, b)
		return p
	}
}

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

func Produce(opts ...Option) ([]characters.Character, error) {
	prd := defProducer()
	return prd.Produce(opts...)
}

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
			return 3
		})(p.Generator)
		prd.source = p.Generator.Generate(3)
	}

	if p.R != nil && p.W == 1 && p.H == 1 {
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

	if prd.R != nil {
		prd.R.SetPos(p.Start.X(), p.Start.Y())
		render.Draw(prd.R, layer.Effect)
	} else {
		prd.R = render.NewEmptySprite(0, 0, 1, 1) //Safety for Mover functionality
	}

	if prd.source != nil {
		prd.source.SetPos(p.Start.X(), p.Start.Y())
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

	return chrs, nil
}

func (p *Product) Init() event.CID {
	p.CID = event.NextID(p)
	return p.CID
}

type Product struct {
	*entities.Interactive
	shouldPersist bool
	position      int
	FollowX       *float64
	FollowY       *float64
	source        *particle.Source
	next          func(floatgeom.Point2)
	buffs         []buff.Buff
}

func (p *Product) MoveParticles(nextDelta floatgeom.Point2) {
	if p.source != nil {
		p.source.ShiftX(nextDelta.X())
		p.source.ShiftY(nextDelta.Y())
	}
}

func (p *Product) Destroy() {
	if p.next != nil {
		p.next(floatgeom.Point2{p.X(), p.Y()})
	}
	p.Interactive.Destroy()
	if p.source != nil {
		p.source.Stop()
	}
}

func (p *Product) Activate() {}

func (p *Product) ShouldPersist() bool {
	return p.shouldPersist
}

func (p *Product) Buffs() []buff.Buff {
	return p.buffs
}
