package run

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oakmound/oak/key"

	"github.com/oakmound/oak/mouse"

	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/characters/enemies"
	"github.com/oakmound/weekly87/internal/restrictor"

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

		runLayers := []render.Stackable{
			// ground groundLayer
			render.NewCompositeR(),
			// wall backgroundLayer
			render.NewCompositeR(),
			// entities 	playLayer
			render.NewHeap(false),
			// effects effectLayer
			render.NewHeap(false),
			// ui uiLayer
			render.NewHeap(true),
		}

		// There should be some way to draw to a stack based
		// on layer name
		if settings.Active.ShowFpsToggle {
			runLayers = append(runLayers, render.NewDrawFPS(), render.NewLogicFPS())
		}

		render.SetDrawStack(runLayers...)

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

		debugInvuln := false

		restrictor.ResetDefault()
		restrictor.Start(1)
		r := records.Load()

		ptycon := players.PartyConstructor{
			Players: players.ClassConstructor(
				r.PartyComp),
			// []int{players.Spearman, players.Mage, players.Mage, players.Swordsman}),

			// []int{players.Spearman, players.Spearman, players.Spearman, players.Spearman}),
		}
		ptycon.Players[0].Position = floatgeom.Point2{players.WallOffset, float64(oak.ScreenHeight / 2)}
		pty, err := ptycon.NewRunningParty()
		if err != nil {
			dlog.Error(err)
			return
		}
		runInfo = records.RunInfo{
			Party:           pty,
			SectionsCleared: 0,
			EnemiesDefeated: 0,
		}
		endLock := sync.Mutex{}
		defeatedShowing := false
		runbackOnce := sync.Once{}

		tracker := section.NewTracker(BaseSeed)

		for i, p := range pty.Players {
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
				if debugInvuln || ply.ForcedInvulnerable || !en.Active {
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
				runbackOnce.Do(func() {
					facing = -1
					event.Trigger("RunBack", nil)
					tracker.Prev()
				})
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

			const aRendDims = 64.0
			const aPad = aRendDims + 12.0 //Size of ability image plus padding
			const cornerPad = 20
			abilityKeys := []string{
				"Q", "W", "E", "R", "T",
			}
			abilityX := float64(cornerPad + i*aPad)

			btnOpts := btn.And(menus.BtnCfgB, btn.Layers(4, 0),
				btn.Pos(abilityX, cornerPad),
				btn.Height(aRendDims),
				btn.Width(aRendDims),
				btn.Text(""),
			)

			// Set up abilities
			if p.Special1 != nil {

				trg := p.Special1.Trigger

				p.Special1.Renderable().SetPos(abilityX, cornerPad)
				btnOpts := btn.And(btnOpts, btn.Renderable(p.Special1.Renderable()),
					btn.Binding(mouse.ClickOn, func(int, interface{}) int {
						trg()
						return 0
					}))
				keyToBind := key.Down + strconv.Itoa(i)
				if i < 10 {
					btnOpts = btn.And(btn.Binding(keyToBind, func(int, interface{}) int {
						trg()
						return 0
					}), btnOpts)
				}
				btn.New(btnOpts)
			}
			if p.Special2 != nil {
				p.Special2.Renderable().SetPos(abilityX, cornerPad+aPad)

				trg := p.Special2.Trigger

				btnOpts := btn.And(btnOpts, btn.Renderable(p.Special2.Renderable()),
					btn.Pos(abilityX, cornerPad+aPad),
					btn.Binding(mouse.ClickOn, func(int, interface{}) int {
						trg()
						return 0
					}))

				btnOpts = btn.And(btn.Binding(key.Down+abilityKeys[i], func(int, interface{}) int {
					trg()
					return 0
				}), btnOpts)

				btn.New(btnOpts)
			}
		}

		var sec1, sec2, sec3 *section.Section

		sec1 = tracker.Next()
		sec2 = tracker.Next()
		sec3 = sec1.Copy()

		sec2.SetBackgroundX(sec1.W())
		sec3.SetBackgroundX(sec1.W() * 2)

		sec1.Draw()
		sec1.ActivateEntities()
		sec2.Draw()
		sec2.ActivateEntities()

		const (
			sec1Mid = 1120
			sec2Mid = 1120 * 3
			sec3Mid = 1120 * 5
		)

		var lastX float64

		// Section creation bind to support infinite* hallway
		event.GlobalBind(func(int, interface{}) int {
			x := pty.Players[0].X()

			if facing == 1 {
				if lastX <= sec2Mid {
					if x > sec2Mid {
						go func() {
							// - A+=2
							// - C+=2
							sec1.Destroy()
							sec3.Destroy()

							sec3 = tracker.Next()
							sec1 = sec3.Copy()
							sec3.SetBackgroundX(sec1.W() * 2)

							sec3.Draw()
							sec3.ActivateEntities()
							sec1.Draw()

							runInfo.SectionsCleared++
						}()
					}
				} else if lastX <= sec3Mid {
					if x > sec3Mid {
						// - Teleport all entities two section widths back (Including the viewport)
						// - B+=2

						pty.ShiftX(-sec1.W() * 2)
						sec3.ShiftEntities(-sec1.W() * 2)
						oak.ViewPos.X -= int(sec1.W()) * 2

						go func() {
							sec2.Destroy()
							sec2 = tracker.Next()
							sec2.SetBackgroundX(sec1.W())
							sec2.Draw()
							sec2.ActivateEntities()

							pty.SpeedUp(1)
							runInfo.SectionsCleared++
						}()
					}
				}
			} else {
				if lastX >= sec2Mid {
					if x < sec2Mid {
						go func() {
							// - A-=2
							// - C-=2
							sec1.Destroy()
							sec3.Destroy()

							sec1 = tracker.Prev()
							sec3 = sec1.Copy()
							sec3.SetBackgroundX(sec1.W() * 2)

							sec3.Draw()
							sec1.Draw()
							sec1.ActivateEntities()

							if tracker.AtStart() {
								oak.SetViewportBounds(0, 0, 8000, 8000)
							}
							runInfo.SectionsCleared++
						}()
					}
				} else if lastX >= sec1Mid {
					if x < sec1Mid && !tracker.AtStart() {
						go func() {
							// - Teleport all entities two section widths forward (Including the viewport)
							// - B-=2
							sec2.Destroy()
							sec2 = tracker.Prev()
							sec2.SetBackgroundX(sec1.W())
							sec2.Draw()
							sec2.ActivateEntities()

							pty.SpeedUp(1)
							runInfo.SectionsCleared++
						}()

						pty.ShiftX(sec1.W() * 2)
						sec1.ShiftEntities(sec1.W() * 2)
						oak.ViewPos.X += int(sec1.W()) * 2
					}
				}
			}

			lastX = x
			return 0
		}, "EnterFrame")

		event.GlobalBind(func(cid int, data interface{}) int {
			dlog.Info("An Enemy Died")

			info := data.([]int64)
			tracker.UpdateHistory(info[0],
				section.Change{
					Typ: section.EntityDestroyed,
					Val: int(info[1])})
			return 0
		}, "EnemyDeath")

		event.GlobalBind(func(cid int, data interface{}) int {
			dlog.Info("A character fired an ability")
			artifacts := data.([]characters.Character)

			abilitySection := sec3
			//Add to appropriate section and potentially to the tracker's changelog for the given section
			if pty.Players[0].X() < (sec1.W()) { //In section 1
				abilitySection = sec1
			} else if pty.Players[0].X() < (2 * sec1.W()) { //In section 2
				abilitySection = sec2
			}
			abilitySection.AppendEntities(artifacts...)

			for _, a := range artifacts {
				if p, ok := a.(Persistable); ok && p.ShouldPersist() {
					tracker.UpdateHistory(abilitySection.GetId(), section.Change{})
				}
			}

			return 0
		}, "AbilityFired")

		music, err = audio.Load(filepath.Join("assets", "audio"), "runIntro.wav")
		dlog.ErrorCheck(err)
		music, err = music.Copy()
		dlog.ErrorCheck(err)
		music = music.MustFilter(
			filter.Volume(0.5 * settings.Active.MusicVolume * settings.Active.MasterVolume),
		)

		music.Play()
		go func() {
			time.Sleep(music.PlayLength())
			music, err = audio.Load(filepath.Join("assets", "audio"), "runLoop.wav")
			dlog.ErrorCheck(err)
			music, err = music.Copy()
			dlog.ErrorCheck(err)
			music = music.MustFilter(
				filter.Volume(0.5*settings.Active.MusicVolume*settings.Active.MasterVolume),
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

		oak.ResetCommands()
		oak.AddCommand("invuln", func(args []string) {
			dlog.Error("Cheating to set the invulnerability toggle")
			debugInvuln = true
			if len(args) > 0 && (args[0][0:0] == "f" || args[0][0:0] == "F") {
				debugInvuln = false
			}
		})
		oak.AddCommand("debug", func(args []string) {
			dlog.Warn("Cheating to toggle debug mode")
			if debugTree.DrawDisabled {
				debugTree.DrawDisabled = false
				return
			}
			debugTree.DrawDisabled = true
		})
		oak.AddCommand("speedup", func(args []string) {
			up := 5.0
			if len(args) > 0 {
				var err error
				up, err = strconv.ParseFloat(args[0], 64)
				dlog.ErrorCheck(err)
			}
			pty.SpeedUp(up)
		})

		dlog.Info("Current Debug Commands are: ", strings.Join(oak.GetDebugKeys(), " , "))
	},
	Loop: scene.BooleanLoop(&stayInGame),
	End: func() (string, *scene.Result) {
		music.Stop()
		restrictor.Stop()
		restrictor.Clear()
		return nextscene, nil
	},
}
