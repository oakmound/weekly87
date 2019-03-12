package characters

import (
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/physics"

	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
)

type Pc struct {
	*entities.Interactive
	SpecialBind1 func()
	SpecialBind2 func()
	animationMap map[string]render.Modifiable
	Job          int
}

func (s *Pc) Init() event.CID {
	return event.NextID(s)
}

func NewPc(job int, x, y float64) *Pc {
	p := &Pc{}
	// r := render.NewColorBox(playerWidth, playerHeight, color.RGBA{255, 0, 0, 255})

	p.setJob(job)

	r := render.NewSwitch("walkRT", p.animationMap)
	p.Interactive = entities.NewInteractive(x, y, playerWidth, playerHeight, r, nil, p.Init(), 0)
	collision.Add(p.RSpace.Space)
	p.Speed = physics.NewVector(0, 5)

	// p.R = render.NewCompoundR("walkRT", p.loadAnimationMap())
	// h.animation = ch.R.(*render.Compound)
	return p
}

func (p *Pc) SetJob(job int) {
	p.setJob(job)
	p.R = render.NewSwitch(p.R.(*render.Switch).Get(), p.animationMap)
}

func (p *Pc) setJob(job int) {
	p.Job = job
	switch job {
	case JobArcher:
		job := &Archer{}
		p.animationMap = job.loadAnimationMap()
		p.SpecialBind1 = job.Special1
		p.SpecialBind2 = job.Special2
	case JobMedic:
		job := &Medic{}
		p.animationMap = job.loadAnimationMap()
		p.SpecialBind1 = job.Special1
		p.SpecialBind2 = job.Special2
	default:
		job := &Swordsman{}
		p.animationMap = job.loadAnimationMap()
		p.SpecialBind1 = job.Special1
		p.SpecialBind2 = job.Special2
	}

	return
}

func (p *Pc) Special1() {

}

func (p *Pc) Special2() {

}
