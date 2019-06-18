package abilities

import (
	"time"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/weekly87/internal/characters"

	"github.com/oakmound/oak/render"
)

var BuffIconSize = 16

// User is something that can use abilities
type User interface {
	Vec() physics.Vector      //Position
	GetDelta() physics.Vector // Speed
	Direction() string        //Facing
	Ready() bool              // Currently are you alive
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
	trigger func(User) []characters.Character
}

// Renderable gets the renderable underlyting the ability
func (a *ability) Renderable() render.Modifiable {
	return a.renderable
}

// Trigger checks if the user is ready and if the ability is off cooldown and then performs the ability if so
func (a *ability) Trigger() {

	if a.user.Ready() && a.renderable.Get(1).(*cooldown).Trigger() {
		artifacts := a.trigger(a.user)
		dlog.Verb("Trigger ability firing")
		event.Trigger("AbilityFired", artifacts)
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
func NewAbility(r render.Modifiable, c time.Duration, t func(User) []characters.Character) *ability {

	w, h := r.GetDims()
	cool := NewCooldown(w, h, c)
	cr := render.NewCompositeM(r, cool)

	return &ability{renderable: cr, trigger: t}

}

//Make produced ability type that captures a created ability
