package abilities

import (
	"fmt"
	"image/color"
	"path/filepath"
	"time"

	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/recolor"
	"github.com/oakmound/weekly87/internal/sfx"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

// Init to be run after oak setup to get our assets set up
func Init() {
	var err error
	blastIcon, err = render.LoadSprite("", filepath.Join("64x64", "BlastIcon.png"))
	dlog.ErrorCheck(err)
	shieldAuraIcon, err = render.LoadSprite("", filepath.Join("64x64", "ShieldAuraIcon.png"))
	dlog.ErrorCheck(err)
	shieldIcon, err = render.LoadSprite("", filepath.Join("64x64", "ShieldIcon.png"))
	dlog.ErrorCheck(err)
	slashIcon, err = render.LoadSprite("", filepath.Join("64x64", "SlashIcon.png"))
	dlog.ErrorCheck(err)
	hammerIcon, err = render.LoadSprite("", filepath.Join("64x64", "HammerIcon.png"))
	dlog.ErrorCheck(err)
	rezIcon, err = render.LoadSprite("", filepath.Join("64x64", "RezIcon.png"))
	dlog.ErrorCheck(err)

	red := color.RGBA{200, 100, 100, 255}
	blue := color.RGBA{100, 100, 200, 255}

	redBlastIcon = blastIcon.Copy().(*render.Sprite)
	redBlastIcon.Filter(recolor.WithStrategy(recolor.ColorMix(red)))
	blueBlastIcon = blastIcon.Copy().(*render.Sprite)
	blueBlastIcon.Filter(recolor.WithStrategy(recolor.ColorMix(blue)))

	redBlastDIcon = redBlastIcon.Copy().(*render.Sprite)
	redBlastDIcon.Modify(mod.Rotate(270))

	blueBlastDIcon = blueBlastIcon.Copy().(*render.Sprite)
	blueBlastDIcon.Modify(mod.Rotate(270))

	blueBlastDIcon = blueBlastIcon.Copy().(*render.Sprite)
	blueBlastDIcon.Modify(mod.Rotate(270))

	upSlashIcon = slashIcon.Copy().(*render.Sprite)
	upSlashIcon.Modify(mod.Rotate(90))

	downSlashIcon = slashIcon.Copy().(*render.Sprite)
	downSlashIcon.Modify(mod.Transpose, mod.Rotate(90))
	// downSLash.Modify()

	bannerSeq, err = render.LoadSheetSequence(
		filepath.Join("16x32", "banner.png"),
		16, 32, 0, 5, []int{0, 0, 1, 0, 2, 0, 3, 0, 0, 1, 1, 1, 2, 1}...)
	dlog.ErrorCheck(err)

	placeHolderBuff, err = render.LoadSprite(
		filepath.Join("assets/images", "16x16"),
		"place_holder_buff.png")
	dlog.ErrorCheck(err)

	mageInit()
	WarriorInit()
}

var (
	blastIcon, shieldAuraIcon, shieldIcon, slashIcon, hammerIcon *render.Sprite
	redBlastIcon, blueBlastIcon, redBlastDIcon, blueBlastDIcon   *render.Sprite
	upSlashIcon, downSlashIcon, rezIcon, placeHolderBuff         *render.Sprite
	bannerSeq                                                    *render.Sequence
	iconW                                                        = 64
	iconH                                                        = 64

	// BuffIconSize is used to determine how to display the buffs icons
	BuffIconSize = 16

	dmg     = map[string]float64{"damage": 1.0}
	baseHit = map[collision.Label]collision.OnHit{
		labels.Enemy: func(a, b *collision.Space) {
			b.CID.Trigger("Attacked", dmg)
		},
	}
)

// User is something that can use abilities
type User interface {
	Vec() physics.Vector      //Position
	GetDelta() physics.Vector // Speed
	Direction() string        //Facing
	DebugEnabled() bool       // If Debug is on
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
	if a.disabled {
		sfx.Play("nope1")
		return
	}
	if !a.cooldown.Trigger() && !a.user.DebugEnabled() {
		sfx.Play("cooldown")
		return
	}

	artifacts := a.trigger(a.user)
	dlog.Verb("Trigger ability firing")
	event.Trigger("AbilityFired", artifacts)

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

// newAbility creates an ability
func newAbility(r render.Modifiable, c time.Duration, t func(User) []characters.Character) *ability {

	inactive := r.Copy()
	inactive.Filter(mod.ConformToPallete(color.GrayModel))
	// iconW, iconH := r.GetDims()

	cool := newCooldown(iconW, iconH, c)
	composite := render.NewCompositeM(r, cool)
	swith := render.NewSwitch("active", map[string]render.Modifiable{
		"active":   composite,
		"inactive": inactive,
	})

	return &ability{renderable: swith, cooldown: cool, trigger: t}
}

//Make produced ability type that captures a created ability
