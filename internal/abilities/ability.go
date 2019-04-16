package abilities

import (
	"time"

	"github.com/oakmound/oak/physics"

	"github.com/oakmound/oak/render"
)

type User interface {
	Vec() physics.Vector //Position
	Direction() string   //Facing
	Ready() bool         // Currently are you alive
}

type Ability interface {
	Renderable() render.Modifiable
	Trigger()
	Cooldown() time.Duration
	SetUser(User) Ability
}

type ability struct {
	renderable render.Modifiable
	cooldown   time.Duration
	user       User
	trigger    func(User)
}

func (a *ability) Renderable() render.Modifiable {
	return a.renderable
}
func (a *ability) Trigger() {
	a.trigger(a.user)
}

func (a *ability) Cooldown() time.Duration {
	return a.cooldown
}
func (a *ability) SetUser(newUser User) Ability {
	return &ability{
		renderable: a.renderable.Copy(),
		cooldown:   a.cooldown,
		user:       newUser,
		trigger:    a.trigger,
	}
}
