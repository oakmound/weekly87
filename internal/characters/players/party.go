package players

import (
	"errors"
	"fmt"
	"image/color"
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/abilities/buff"
	"github.com/oakmound/weekly87/internal/sfx"

	"github.com/oakmound/oak/key"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/enemies"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/joys"
	"github.com/oakmound/weekly87/internal/vfx"
)

// Party handles having multiple pcs and the info to manage and display them!
type Party struct {
	event.CID
	Players      []*Player
	Acceleration float64
	speedUps     float64
	joystickID   uint32
	Debug        bool
}

// Init the party giving them a CID
func (p *Party) Init() event.CID {
	return event.NextID(p)
}

// SpeedUp the party based on the contained alg
// We use a log approach to make sure that speedup falls off
func (p *Party) SpeedUp(n float64) {
	p.speedUps += n
	p.Acceleration = math.Log10(math.Pow(
		math.Log10(p.speedUps+10), 2)) * 15
	if p.Players[0].RunSpeed == 0 {
		p.Acceleration = 0
	}

}

// CheckedBind wraps binding to the party performing our standard checks
func (p *Party) CheckedBind(bnd func(*Party, interface{}) int, ev string) {
	p.Bind(func(id int, data interface{}) int {
		be, ok := event.GetEntity(id).(*Party)
		if !ok {
			dlog.Error("Party binding was called on non-party")
			return event.UnbindSingle
		}
		return bnd(be, data)
	}, ev)
}

// RunSpeed retrieves the current speed for the party to run at
func (p *Party) RunSpeed() int {
	if p.Players[0].facing == "LT" {
		return int(p.Players[0].RunSpeed - p.Acceleration)
	}
	return int(p.Players[len(p.Players)-1].RunSpeed + p.Acceleration)
}

// Speed returns the party's speed vector
func (p *Party) Speed() physics.Vector {
	if p.Players[0].facing == "LT" {
		return p.Players[0].Speed
	}
	return p.Players[len(p.Players)-1].Speed
}

// Defeated checks if the party is in a defeated state
func (p *Party) Defeated() bool {
	for _, pl := range p.Players {
		if pl.Alive {
			return false
		}
	}
	return true
}

// ShiftX for the entire party
func (p *Party) ShiftX(f float64) {
	for _, pl := range p.Players {
		pl.ShiftX(f)
	}
}

//PartyConstructor helps set up a party
type PartyConstructor struct {
	Players    []Constructor
	Bindings   map[string]func(*Party, interface{}) int
	MaxPlayers int
}

// NewRunningParty creates a party for the run scene
func (pc *PartyConstructor) NewRunningParty() (*Party, error) {
	return pc.NewParty(false)
}

// PlayerGap is the xgap imbetween each member of the party
const PlayerGap = 50

// NewParty sets up the party from a constructor
// The party may be a moving or unmoving party representation for interaction or for display
func (pc *PartyConstructor) NewParty(unmoving bool) (*Party, error) {
	if len(pc.Players) == 0 {
		return nil, errors.New("At least one player must be in a party")
	}
	for len(pc.Players) < pc.MaxPlayers && unmoving {
		// Buffer with invisible players
		// Note this is an assumption about future classes
		pc.Players = append(pc.Players, *EmptyConstructor)
	}

	pty := &Party{}

	for i, pcon := range pc.Players {
		if !unmoving && pcon.RunSpeed == -1 {
			continue
		}
		if pcon.Dimensions == (floatgeom.Point2{}) {
			return nil, errors.New("Dimensions must be provided for player " + strconv.Itoa(i))
		}
		for _, s := range requiredAnimations {
			if _, ok := pcon.AnimationMap[s]; !ok {
				return nil, errors.New("Animation name " + s + " must be provided for player " + strconv.Itoa(i))
			}
		}
		p := Player{}
		p.PartyIndex = i
		p.Status = &buff.Status{}

		if pcon.Special1 != nil {
			p.Special1 = pcon.Special1.SetUser(&p)

		}
		if pcon.Special2 != nil {
			p.Special2 = pcon.Special2.SetUser(&p)
		}

		p.Swtch = render.NewSwitch("walkRT", pcon.AnimationMap)
		p.Interactive = entities.NewInteractive(
			pc.Players[0].Position.X()+float64(i)*PlayerGap,
			pc.Players[0].Position.Y(),
			pcon.Dimensions.X(),
			pcon.Dimensions.Y(),
			p.Swtch,
			nil,
			p.Init(),
			0,
		)
		p.facing = "RT"
		p.Alive = true
		p.Speed = physics.NewVector(pcon.Speed.X(), pcon.Speed.Y())
		p.RunSpeed = pcon.RunSpeed
		p.RSpace.UpdateLabel(labels.PC)
		if i != 0 {
			p.Delta = pty.Players[0].Delta
		}
		p.Name = pcon.Name
		p.AccruedValue = pcon.AccruedValue
		// Interaction with Enemies
		p.RSpace.Add(labels.Enemy, func(s, e *collision.Space) {
			ply, ok := s.CID.E().(*Player)
			if !ok {
				dlog.Error("Non-player sent to player binding")
				return
			}
			en, ok := e.CID.E().(*enemies.BasicEnemy)
			if !ok {
				dlog.Error("Non-enemy sent to enemy binding")
				fmt.Printf("%T\n", s.CID.E())
				return
			}
			if ply.Invulnerable > 0 || !en.Active {
				return
			}

			if ply.Shield > 0 {

				facingRight := (p.facing == "RT")
				direction := 1.0
				if facingRight {
					direction = -1.0
				}
				pushD := -1*direction + p.RunSpeed*5
				dlog.Info("Enemy hit us be we were shielded so we pushed them back by", pushD)

				vfx.VerySmallShaker.Shake(time.Duration(400) * time.Millisecond)
				sfx.Play("bounced1")
				// Affect the enemy
				en.PushBack(physics.NewVector(pushD, 0))

				source := vfx.PushBack1().Generate(2)
				plyX := ply.X() - 5
				if facingRight {
					plyX += 26
				}
				source.SetPos(plyX, ply.Y()+16)
				endSource := time.Now().Add(time.Millisecond * 300)
				source.CID.Bind(func(id int, data interface{}) int {
					eff, ok := event.GetEntity(id).(*particle.Source)
					if ok {
						eff.ShiftX(ply.Delta.X() + 1)

						if endSource.Before(time.Now()) {
							eff.Stop()
							return 1
						}
					}

					return 0
				}, "EnterFrame")

				// Remove the charge from our buffs
				for buffIdx, b := range ply.Buffs {
					if b.Name == buff.NameShield {
						b.Charges--
						if b.Charges <= 0 {
							b.ExpireAt = time.Now()
						}
						ply.Buffs[buffIdx] = b

						//TODO: Consider have shields create different pushbacks

						return
					}
				}
				dlog.Warn("We thought we had shield but we could not find a buff with such a name")
				return
			}

			abilities.Produce(
				abilities.StartAt(floatgeom.Point2{ply.X() + 8, ply.Y() + 10}),
				//abilities.FollowSpeed(ply.Delta.Xp(), ply.Delta.Yp()),
				abilities.WithParticles(vfx.WhiteRing()),
				abilities.Duration(time.Millisecond*20),
			)
			vfx.SmallShaker.Shake(time.Duration(1000) * time.Millisecond)
			sfx.Play("playerHit1")

			ply.Kill()
			event.Trigger("PlayerDeath", nil)
		})

		// Hitting Chests
		p.RSpace.Add(labels.Chest, func(s, s2 *collision.Space) {
			p, ok := s.CID.E().(*Player)
			if !ok {
				dlog.Error("Non-player sent to player binding")
				return
			}
			ch, ok := s2.CID.E().(*doodads.Chest)
			if !ok {
				dlog.Error("Non-chest sent to chest binding")
				return
			}
			if !ch.Active {
				return
			}
			r := ch.R.(render.Modifiable).Copy()
			_, h := r.GetDims()

			p.AddChest(h, r, ch.Value)

			// p.ChestsHeight += float64(h)
			// chestHeight := p.ChestsHeight

			// r.(*render.Sprite).Vector = r.Attach(p.Vector, -3, -chestHeight)
			// p.ChestValues = append(p.ChestValues, ch.Value)
			// p.Chests = append(p.Chests, r)
			// render.Draw(r, layer.Play, 2)

			ch.Destroy()

			event.Trigger("RunBackOnce", nil)
		})

		// Hitting buffs
		p.RSpace.Add(labels.EffectsPlayer, func(s, bf *collision.Space) {
			p, ok := s.CID.E().(*Player)
			if !ok {
				dlog.Error("Non-player sent to player binding")
				return
			}
			bfr, ok := bf.CID.E().(Buffer)
			if !ok {
				dlog.Error("EffectsPlayer label on non-Effecter")
				return
			}
			// Todo: How do we know if the buff is a party wide buff or not
			pty := p.Party
			if pty == nil {
				dlog.Error("Player had no party")
				return
			}
			bfs := bfr.Buffs()
			for _, b := range bfs {
				if b.Name == buff.NameRez {
					for _, ply := range pty.Players {
						if !ply.Alive {
							ply.Revive()
							break
						}
					}
				} else {
					if b.SinglePlayer {
						p.AddBuff(b)
						continue
					}
					for _, ply := range pty.Players {
						if ply.Alive {
							ply.AddBuff(b)
						}
					}
				}
			}
			if dstr, ok := bfr.(Destroyable); ok {
				dstr.Destroy()
			}
			//bf.CID.Trigger("Hit", nil)
		})

		p.CheckedBind(func(p *Player, _ interface{}) int {
			p.facing = "LT"
			if len(p.ChestValues) > 0 {
				p.Swtch.Set("walkHold")
			} else {
				if !p.Alive {
					p.Swtch.Set("dead" + p.facing)
				} else {
					p.Swtch.Set("walk" + p.facing)
				}
			}
			return 0
		}, "RunBack")

		for ev, b := range pcon.Bindings {
			p.CheckedBind(b, ev)
		}
		p.Party = pty
		pty.Players = append(pty.Players, &p)
	}

	pty.CID = pty.Init()

	lowestID := joys.LowestID()
	if lowestID != math.MaxInt32 {
		pty.joystickID = lowestID
	}
	if unmoving {
		return pty, nil
	}
	pty.CheckedBind(func(pty *Party, _ interface{}) int {
		for i, p := range pty.Players {
			// Lean towards being generous
			p.AddBuff(buff.Invulnerable(render.NewColorBox(8, 8, color.RGBA{255, 255, 0, 255}), 5*time.Second))
			i := i
			p.RunSpeed *= -1
			p.CheckedBind(func(p *Player, _ interface{}) int {
				// Shift the player back until against the right wall
				if int(p.X())-oak.ViewPos.X >= oak.ScreenWidth-(WallOffset+(len(pty.Players)-1-i)*PlayerGap) {
					return event.UnbindSingle
				}
				p.ShiftX(float64(-pty.RunSpeed()) * 2)
				return 0
			}, "EnterFrame")
		}
		return event.UnbindSingle
	}, "RunBack")

	pty.CheckedBind(func(pty *Party, _ interface{}) int {
		// Find the first player that's dead
		for _, p := range pty.Players {
			if !p.Alive {
				p.Revive()
				break
			}
		}
		return 0
	}, "Rez")

	buffIcon, err := render.LoadSprite(filepath.Join("assets/images", "16x16"), "place_holder_buff.png")
	dlog.ErrorCheck(err)

	pty.CheckedBind(func(pty *Party, _ interface{}) int {
		for _, p := range pty.Players {
			p.AddBuff(buff.Rage(buffIcon, 2*time.Second))
		}
		return 0
	}, "RageStart")

	pty.CheckedBind(func(pty *Party, _ interface{}) int {
		p0 := pty.Players[0]
		p0.Delta.Zero()

		js := joys.StickState(pty.joystickID)

		p0.Delta.SetX(float64(pty.RunSpeed()))
		if p0.Status.Rage <= 0 {
			if oak.IsDown(key.UpArrow) || js.StickLY > 8000 {
				p0.Delta.ShiftY(-pty.Speed().Y())
			}
			if oak.IsDown(key.DownArrow) || js.StickLY < -8000 {
				p0.Delta.ShiftY(pty.Speed().Y())
			}
		}

		p0.Vector.Add(p0.Delta)

		_, h := p0.Swtch.GetDims()
		hf := float64(h)
		if p0.Vector.Y() < float64(oak.ScreenHeight)*1/3 {
			p0.Vector.SetY(float64(oak.ScreenHeight) * 1 / 3)
		} else if p0.Vector.Y() > (float64(oak.ScreenHeight) - hf) {
			p0.Vector.SetY((float64(oak.ScreenHeight) - hf))
		}

		for i := 1; i < len(pty.Players); i++ {
			p := pty.Players[i]
			p.Vector.ShiftX(p0.Delta.X())
			p.Vector.SetY(p0.Vector.Y())
		}
		flashStartTime := time.Now().Add(time.Second * 5)
		flashCounter := 5
		for _, p := range pty.Players {
			// The idea behind splitting up the move functions is
			// flawed when they're all working together--we only want
			// to shift everything -once-, otherwise there are jitters
			// or other awkward bits to moving around.
			p.R.SetPos(p.Vector.X(), p0.Vector.Y())

			for len(p.Buffs) > 0 {
				if p.Buffs[0].ExpireAt.Before(time.Now()) {
					p.BuffLock.Lock()
					p.Buffs[0].Disable(p.Status)
					p.Buffs[0].R.Undraw()
					p.Buffs = p.Buffs[1:]
					p.BuffLock.Unlock()
					p.ReorderBuffs()
				} else {
					break
				}
			}
			for bIndex := 0; bIndex < len(p.Buffs); bIndex++ {
				if p.Buffs[bIndex].PreExpireCounter == 0 && p.Buffs[bIndex].ExpireAt.Before(flashStartTime) {
					p.Buffs[bIndex].PreExpireCounter++
					switchBuffR(&p.Buffs[bIndex])
					//p.Buffs[bIndex].R.SetRGBA(p.Buffs[bIndex].AltRenders.GetRGBA())

				} else if p.Buffs[bIndex].PreExpireCounter > 0 {
					p.Buffs[bIndex].PreExpireCounter++
					if p.Buffs[bIndex].PreExpireCounter > flashCounter {
						p.Buffs[bIndex].PreExpireCounter = 1
						switchBuffR(&p.Buffs[bIndex])
						//fmt.Println(p.Buffs[bIndex].R.(*render.Switch).Get())
					}

				} else {
					break
				}
			}

			if !p.Alive {
				continue
			}
			p.RSpace.Update(p.Vector.X(), p0.Vector.Y(), p.RSpace.GetW(), p.RSpace.GetH())
			<-p.RSpace.CallOnHits()
		}

		oak.ShiftScreen(int(pty.RunSpeed()), 0)

		return 0
	}, "EnterFrame")

	for ev, b := range pc.Bindings {
		pty.CheckedBind(b, ev)
	}

	return pty, nil
}

// switchBuffR is a utility fxn for buff update
func switchBuffR(b *buff.Buff) *buff.Buff {
	keyProgression := b.R.Get()
	if keyProgression == "base" {
		b.R.Set("flicker")
	} else {
		b.R.Set("base")
	}
	return b
}
