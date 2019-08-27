package abilities

import (
	"fmt"
	"image/color"
	"time"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/weekly87/internal/characters"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

var BuffIconSize = 16

// User is something that can use abilities
type User interface {
	Vec() physics.Vector      //Position
	GetDelta() physics.Vector // Speed
	Direction() string        //Facing
}

// Ability is an action with an associated UI element that can be invoked
type Ability interface {
	Renderable() render.Modifiable
	Trigger()
	Cooldown() time.Duration
	SetUser(User) Ability
	Enable(bool)
	SetButton(btn.Btn)
}

type ability struct {
	renderable *render.Switch
	cooldown   *cooldown
	button     btn.Btn

	disabled bool

	user    User
	trigger func(User) []characters.Character
}

func (a *ability) SetButton(b btn.Btn) {
	a.button = b
}

// Renderable gets the renderable underlyting the ability
func (a *ability) Renderable() render.Modifiable {
	return a.renderable
}

// Trigger checks if the user is ready and if the ability is off cooldown and then performs the ability if so
func (a *ability) Trigger() {

	if !a.disabled && a.cooldown.Trigger() {
		artifacts := a.trigger(a.user)
		dlog.Verb("Trigger ability firing")
		event.Trigger("AbilityFired", artifacts)
	}
}

// Cooldown gets the total cooldown time  for the ability
func (a *ability) Cooldown() time.Duration {
	return a.cooldown.totalTime
}

// SetUser copies the ability and sets the user on it making a nice unique instance
func (a *ability) SetUser(newUser User) Ability {
	r := a.renderable.Copy().(*render.Switch)
	cool := r.GetSub("active").(*render.CompositeM).Get(1).(*cooldown)
	cool.ResetTiming()
	return &ability{
		renderable: r,
		cooldown:   cool,
		user:       newUser,
		trigger:    a.trigger,
	}
}

func (a *ability) Enable(enabled bool) {
	a.disabled = !enabled
	if a.disabled {
		fmt.Println("Disabled ability")
		a.button.SetMetadata("inactive", "yes")
		rvt := a.button.GetRenderable().(*render.Reverting)
		rvt.Set("inactive")
	} else {
		// Todo: When reviving happens, there's going to be a bug
		// where the button could still be highlighted, because
		// we can't check for highlighting when we disable
		fmt.Println("Enabled ability")
		a.button.SetMetadata("inactive", "")
		rvt := a.button.GetRenderable().(*render.Reverting)
		rvt.Set("active")
	}
}

// NewAbility creates an ability
func NewAbility(r render.Modifiable, c time.Duration, t func(User) []characters.Character) *ability {

	inactive := r.Copy()
	inactive.Filter(mod.ConformToPallete(color.GrayModel))
	w, h := r.GetDims()
	//inactive := render.NewColorBox(w, h, color.RGBA{10, 10, 10, 255})
	cool := NewCooldown(w, h, c)
	composite := render.NewCompositeM(r, cool)
	swith := render.NewSwitch("active", map[string]render.Modifiable{
		"active":   composite,
		"inactive": inactive,
	})

	return &ability{renderable: swith, cooldown: cool, trigger: t}
}

//Make produced ability type that captures a created ability
