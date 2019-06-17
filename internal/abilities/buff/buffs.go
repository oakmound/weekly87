package buff

import (
	"time"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

type Buff struct {
	Duration         time.Duration
	ExpireAt         time.Time
	Charges          int
	PreExpireCounter int
	Name             Name
	Enable           func(*Status)
	Disable          func(*Status)
	R                render.Switch
	RGen             func() render.Modifiable
}

type Name int

const (
	unnamed = iota
	NameShield
)

type Status struct {
	Invulnerable int
	Shield       int
}

// BasicBuffSwitch is a utlity that creates our standard flicker setup
func BasicBuffSwitch(base render.Modifiable) render.Switch {
	flick := base.Copy()
	flick.Filter(mod.Fade(120))
	return *render.NewSwitch("base", map[string]render.Modifiable{"base": base.Copy(), "flicker": flick})
}

func Invulnerable(r render.Modifiable, dur time.Duration) Buff {
	return Buff{
		Duration: dur,
		Enable: func(s *Status) {
			s.Invulnerable++
		},
		Disable: func(s *Status) {
			s.Invulnerable--
		},

		RGen: func() render.Modifiable {
			return r.Copy()
		},
	}
}
func Shield(r render.Modifiable, dur time.Duration, charges int) Buff {
	return Buff{
		Duration: dur,
		Enable: func(s *Status) {
			s.Shield++
		},
		Disable: func(s *Status) {
			s.Shield--
		},
		RGen: func() render.Modifiable {
			return r.Copy()
		},
		Name:    NameShield,
		Charges: charges,
	}
}
