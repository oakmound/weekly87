package abilities

import (
	"time"

	"github.com/oakmound/oak/physics"

	"github.com/oakmound/oak/render"
)

// User is something that can use abilities
type User interface {
	Vec() physics.Vector //Position
	Direction() string   //Facing
	Ready() bool         // Currently are you alive
}

// Ability is an action with an associated UI element that can be invoked
type Ability interface {
	Renderable() render.Modifiable
	Trigger()
	Cooldown() time.Duration
	SetUser(User) Ability
}

type ability struct {
	renderable *render.CompositeM

	user    User
	trigger func(User)
}

func (a *ability) Renderable() render.Modifiable {
	return a.renderable
}
func (a *ability) Trigger() {

	if a.user.Ready() && a.renderable.Get(1).(*cooldown).Trigger() {
		a.trigger(a.user)
	}
}

func (a *ability) Cooldown() time.Duration {
	return a.renderable.Get(1).(*cooldown).totalTime
}
func (a *ability) SetUser(newUser User) Ability {
	composite := a.renderable.Copy().(*render.CompositeM)
	composite.Get(1).(*cooldown).ResetTiming()
	return &ability{
		renderable: composite,
		user:       newUser,
		trigger:    a.trigger,
	}
}
func NewAbility(r render.Modifiable, c time.Duration, t func(User)) *ability {

	w, h := r.GetDims()
	cool := NewCooldown(w, h, c)
	cr := render.NewCompositeM(r, cool)

	return &ability{renderable: cr, trigger: t}

}
