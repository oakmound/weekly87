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

// Renderable gets the renderable underlyting the ability
func (a *ability) Renderable() render.Modifiable {
	return a.renderable
}

// Trigger checks if the user is ready and if the ability is off cooldown and then performs the ability if so
func (a *ability) Trigger() {

	if a.user.Ready() && a.renderable.Get(1).(*cooldown).Trigger() {
		a.trigger(a.user)
	}
}

// Cooldown gets the total cooldown time  for the ability
func (a *ability) Cooldown() time.Duration {
	return a.renderable.Get(1).(*cooldown).totalTime
}

// SetUser copies the ability and sets the user on it making a nice unique instance
func (a *ability) SetUser(newUser User) Ability {
	composite := a.renderable.Copy().(*render.CompositeM)
	composite.Get(1).(*cooldown).ResetTiming()
	return &ability{
		renderable: composite,
		user:       newUser,
		trigger:    a.trigger,
	}
}

// NewAbility creates an ability
func NewAbility(r render.Modifiable, c time.Duration, t func(User)) *ability {

	w, h := r.GetDims()
	cool := NewCooldown(w, h, c)
	cr := render.NewCompositeM(r, cool)

	return &ability{renderable: cr, trigger: t}

}
