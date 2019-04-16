package run

import (
	"fmt"
	"image/color"
	"path/filepath"
	"sync"
	"time"

	"github.com/oakmound/oak/mouse"

	"github.com/oakmound/weekly87/internal/characters/enemies"

	"github.com/oakmound/weekly87/internal/settings"

	klg "github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/audio/filter"

	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/records"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/dtools"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/run/section"
)

var stayInGame bool
var nextscene string
var BaseSeed int64
var music klg.Audio

var runInfo records.RunInfo

// facing is whether is game is moving forward or backward,
// 1 means forward, -1 means backward
var facing = 1

// Scene  to display the run
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInGame = true
		nextscene = "endGame"
		facing = 1

		runInfo = records.RunInfo{
			Party:           []*players.Player{},
			SectionsCleared: 0,
		}

		// There should be some way to draw to a stack based
		// on layer name
		if settings.ShowFpsToggle {
			render.SetDrawStack(
				// ground
				render.NewCompositeR(),
				// maybe background / parallax?
				// wall
				render.NewCompositeR(),
				// entities
				render.NewHeap(false),
				// maybe effects?
				// ui
				render.NewHeap(true),
				render.NewDrawFPS(),
				render.NewLogicFPS(),
			)
		} else {
			render.SetDrawStack(
				render.NewCompositeR(),
				render.NewCompositeR(),
				render.NewHeap(false),
				render.NewHeap(true),
			)
		}

		debugTree := dtools.NewRTree(collision.DefTree)
		debugTree.ColorMap = map[collision.Label]color.RGBA{
			labels.Chest:        color.RGBA{255, 255, 0, 255},
			labels.Door:         color.RGBA{125, 125, 125, 255},
			labels.Enemy:        color.RGBA{0, 0, 255, 255},
			labels.EnemyAttack:  color.RGBA{255, 0, 0, 255},
			labels.PC:           color.RGBA{125, 0, 255, 255},
			labels.PlayerAttack: color.RGBA{255, 0, 125, 255},
		}
		render.Draw(debugTree, 2, 1000)

		ptycon := players.PartyConstructor{
			Players: []players.Constructor{
				*players.SpearmanConstructor.Copy(),
				*players.MageConstructor.Copy(),
				*players.MageConstructor.Copy(),
				*players.SwordsmanConstructor.Copy(),
			},
		}
		ptycon.Players[0].Position = floatgeom.Point2{players.WallOffset, float64(oak.ScreenHeight / 2)}
		pty, err := ptycon.NewParty()
		if err != nil {
			dlog.Error(err)
			return
		}

		endLock := sync.Mutex{}
		defeatedShowing := false

		tracker := section.NewTracker(BaseSeed)
		sct := tracker.Next()
		sct.Draw()
		nextSct := tracker.Next()
		nextSct.SetBackgroundX(sct.X() + sct.W())
		nextSct.Draw()
		nextSct.ActivateEntities()
		var oldSct *section.Section

		facingLock := sync.Mutex{}

		for _, p := range pty.Players {
			render.Draw(p.R, 2, 2)
			rs := p.GetReactiveSpace()

			// Interaction with Enemies
			rs.Add(labels.Enemy, func(s, e *collision.Space) {
				ply, ok := s.CID.E().(*players.Player)
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
				if ply.ForcedInvulnerable || !en.Active {
					return
				}
				ply.Alive = false
				for _, r := range ply.Chests {
					r.Undraw()
				}
				ply.ChestValues = []int64{}
				ply.Trigger("Kill", nil)

				endLock.Lock()
				defer endLock.Unlock()
				if pty.Defeated() && !defeatedShowing {
					for _, ply := range pty.Players {
						ply.RunSpeed = 0
					}
					pty.Acceleration = 0
					defeatedShowing = true
					// Show pop up to go to endgame scene
					menuX := (float64(oak.ScreenWidth) - 180) / 2
					menuY := float64(oak.ScreenHeight) / 4
					btn.New(menus.BtnCfgB, btn.Layers(3, 0),
						btn.Pos(menuX, menuY), btn.Text("Defeated! See Your Stats?"),
						btn.Width(180),
						btn.Binding(mouse.ClickOn, func(int, interface{}) int {
							nextscene = "endGame"
							stayInGame = false

							return 0
						}))
				}
			})

			rs.Add(labels.Chest, func(s, s2 *collision.Space) {
				p, ok := s.CID.E().(*players.Player)
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

				chestHeight := (len(p.ChestValues) + 1) * (h + 1)

				r.(*render.Sprite).Vector = r.Attach(p.Vector, -3, -float64(chestHeight))
				p.ChestValues = append(p.ChestValues, ch.Value)
				p.Chests = append(p.Chests, r)

				ch.Destroy()
				render.Draw(r, 2, 2)
				facingLock.Lock()
				if facing == 1 {

					facing = -1
					facingLock.Unlock()

					event.Trigger("RunBack", nil)

					// Shift sections
					if tracker.At() > 3 {
						tracker.ShiftDepth(-1)
					}
					oldSct = nextSct
					nextSct = tracker.Prev()
					nextSct.SetBackgroundX(sct.X() - sct.W())
				} else {
					facingLock.Unlock()
				}
			})

			// Player got back to the Inn!
			rs.Add(labels.Door, func(_, d *collision.Space) {
				d.CID.Trigger("RibbonCut", nil)
				go func() {
					time.Sleep(500 * time.Millisecond)
					nextscene = "endGame"
					stayInGame = false
				}()
			})
		}

		// Section creation bind to support infinite* hallway
		event.GlobalBind(func(int, interface{}) int {
			// This calculation needs to be modified based
			// on how much of the screen a section takes up.
			// If a section takes up more than one screen,
			// this is fine, otherwise it needs to change a little
			w := sct.W() * float64(facing)
			var offLeft int
			var shift bool
			if facing == 1 {
				offLeft = oak.ViewPos.X - int(w)
				shift = offLeft >= 0
			} else {
				offLeft = oak.ViewPos.X - int(w)
				shift = offLeft <= -int(w)
			}
			if shift && !tracker.AtStart() {
				if oldSct != nil {
					nextSct.Shift(-w)
					sct.Shift(-w)
					oldSct.Destroy()
					oldSct = nil
				} else {
					sct.Destroy()
					sct = nextSct
				}
				// We need a way to make these actions draw-level atomic
				// Or a way to fake it so there isn't a blip
				oak.ViewPos.X = offLeft
				pty.ShiftX(-w)
				// event.DefaultBus.Pause()
				// time.Sleep(5 * time.Second)
				// event.DefaultBus.Resume()
				nextSct.Shift(-w)
				go func() {
					nextSct = tracker.Produce(int64(facing))
					pty.SpeedUp()
					//fmt.Println("Sec", nextSct.GetID(), "total", tracker.SectionsDeep())
					nextSct.SetBackgroundX(sct.X() + w)
					nextSct.Draw()
					nextSct.ActivateEntities()
					if tracker.AtStart() {
						oak.SetViewportBounds(0, 0, 4000, 4000)
					}
					// fmt.Println("sections", runInfo.SectionsCleared)
					runInfo.SectionsCleared++
				}()
			}
			return 0
		}, "EnterFrame")

		music, err = audio.Load(filepath.Join("assets", "audio"), "runIntro.wav")
		dlog.ErrorCheck(err)
		music, err = music.Copy()
		dlog.ErrorCheck(err)
		music = music.MustFilter(
			filter.Volume(0.5 * settings.MusicVolume * settings.MasterVolume),
		)

		music.Play()
		go func() {
			time.Sleep(music.PlayLength())
			music, err = audio.Load(filepath.Join("assets", "audio"), "runLoop.wav")
			dlog.ErrorCheck(err)
			music, err = music.Copy()
			dlog.ErrorCheck(err)
			music = music.MustFilter(
				filter.Volume(0.5*settings.MusicVolume*settings.MasterVolume),
				filter.LoopOn(),
			)
			music.Play()
		}()

		// Maybe there's a countdown timer

		// The state of the game is generated based on combining a base seed
		// and the current section the player is in. When the game is first started
		// base seed is populated randomly and stored in a settings file, then
		// incremented as sections are cleared

		// We also need to keep track of changes to each section like enemies destroyed
		// This means map[int64][]int, where the slice is list of enemies destroyed
		// with enemies identified by order they are made in

		// Background should probably be very basic hallway with tile types
		// and different themes populate the tile types

		// Character types: 10 Sec cooldown ability, 30 sec
		// Spearman - Shove up - Attack in front
		// Warrior - Shove back - Shove all enemies back
		// Cleric - Slow down run speed for short period - Revive
		// Ranger - Y Speed boost for short period - Shoot arrow in front
		// Rogue - Invisible for short period - Blink / jump forward
		// Paladin - Invincible for short period -
		// Mage - Spawns Fire - Freeze all enemies in place

		// Enemy types:
		// 1. Stands Still or walks in a basic path
		// 2. Charges right to left and hurts if you touch it
		// 2a. As 2, but can slightly curve during charge

		// Treasure types:
		// bigger and fancier boxes with fancier colors
		// maybe 3 sizes
		// Color order: brown, then dim gray, then black,
		// white, yellow, green, blue, purple, orange, red, silver

		// for right now the one character can chain boxes behind them
		// infinitely, eventually there should be an upgrade thing
	},
	Loop: scene.BooleanLoop(&stayInGame),
	End: func() (string, *scene.Result) {
		music.Stop()
		return nextscene, nil
	},
}
