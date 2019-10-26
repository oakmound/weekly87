package run

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"sync"
	"time"

	klg "github.com/200sc/klangsynthese/audio"

	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/enemies"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/records"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/restrictor"
	"github.com/oakmound/weekly87/internal/dtools"
	"github.com/oakmound/weekly87/internal/layer"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/music"
	"github.com/oakmound/weekly87/internal/run/section"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/timing"
	"github.com/oakmound/oak/scene"
)

var stayInGame bool
var nextscene string
var BaseSeed int64
var bkgMusic *klg.Audio

var runInfo records.RunInfo

// facing is whether is game is moving forward or backward,
// 1 means forward, -1 means backward
var facing = 1

// Scene  to display the run
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {

		// Reset to start of run if not first run
		stayInGame = true
		nextscene = "endGame"
		facing = 1

		render.SetDrawStack(layer.Get()...)

		debugTree := dtools.NewRTree(collision.DefTree)
		debugTree.ColorMap = labels.ColorMap
		render.Draw(debugTree, layer.Debug, 20)

		restrictor.ResetDefault()
		restrictor.Start(1)

		ptycon := players.PartyConstructor{
			Players: players.ClassConstructor(records.Load().PartyComp),
		}
		ptycon.Players[0].Position = floatgeom.Point2{players.WallOffset, float64(oak.ScreenHeight / 2)}
		pty, err := ptycon.NewRunningParty()
		dlog.ErrorCheck(err)

		runInfo = records.RunInfo{
			Party:           pty,
			SectionsCleared: 0,
			EnemiesDefeated: 0,
		}

		tracker := section.NewTracker(BaseSeed)

		for i, p := range pty.Players {
			render.Draw(p.R, layer.Play, 2)
			rs := p.GetReactiveSpace()
			// Player got back to the Inn!
			rs.Add(labels.Door, func(_, d *collision.Space) {
				d.CID.Trigger("RibbonCut", nil)
				go func() {
					time.Sleep(500 * time.Millisecond)
					nextscene = "endGame"
					stayInGame = false
				}()
			})

			// Ability icon rendering / binding
			const aRendDims = 64.0
			const aPad = aRendDims + 12.0 //Size of ability image plus padding
			const cornerPad = 20
			abilityKeys := []string{
				"Q", "W", "E", "R", "T",
			}
			abilityX := float64(cornerPad + i*aPad)

			btnOpts := btn.And(menus.BtnCfgB, btn.Layers(layer.UI, 0),
				btn.Pos(abilityX, cornerPad),
				btn.Height(aRendDims),
				btn.Width(aRendDims),
				btn.Text(""),
			)

			joyBtns := [4]string{
				"A",
				"X",
				"B",
				"Y",
			}

			// Set up abilities
			if p.Special1 != nil {

				trg := p.Special1.Trigger

				p.Special1.Renderable().SetPos(abilityX, cornerPad)
				btnOpts := btn.And(btnOpts, btn.Renderable(p.Special1.Renderable()),
					btn.Binding(mouse.ClickOn, func(int, interface{}) int {
						trg()
						return 0
					}))
				keyToBind := key.Down + strconv.Itoa(i+1)
				if i < 10 {
					// Todo joystick triggers
					btnOpts = btn.And(btnOpts, btn.Binding(keyToBind, func(int, interface{}) int {
						trg()
						return 0
					}))
				}

				if i < 4 {
					btnOpts = btn.And(btnOpts, btn.Binding(joyBtns[i]+joystick.ButtonUp, func(_ int, state interface{}) int {
						jState := state.(*joystick.State)
						fmt.Println("Joystick TriggerR", jState.TriggerR)
						if jState.TriggerR < 100 {
							fmt.Println("Triggering ability 1")
							trg()
						}
						return 0
					}))
				}
				newBtn := btn.New(btnOpts)
				p.Special1.SetButton(newBtn)
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

				if i < 10 {
					btnOpts = btn.And(btnOpts, btn.Binding(key.Down+abilityKeys[i], func(int, interface{}) int {
						trg()
						return 0
					}))
				}
				if i < 4 {
					btnOpts = btn.And(btnOpts, btn.Binding(joyBtns[i]+joystick.ButtonUp, func(_ int, state interface{}) int {
						jState := state.(*joystick.State)
						fmt.Println("Joystick TriggerR", jState.TriggerR)
						if jState.TriggerR > 100 {
							fmt.Println("Triggering ability 1")
							trg()
						}
						return 0
					}))
				}

				newBtn := btn.New(btnOpts)
				p.Special2.SetButton(newBtn)

			}
		}

		sec1 := tracker.Next()
		sec2 := tracker.Next()
		sec3 := sec1.Copy()

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

		// Create a debug for Section drawing
		secDebugHeight := 20
		secDebug1 := render.NewColorBox(oak.ScreenWidth/3, secDebugHeight, color.RGBA{100, 2, 2, 100})
		secDebug2 := render.NewColorBox(oak.ScreenWidth/3, secDebugHeight, color.RGBA{2, 100, 2, 100})
		secDebug3 := render.NewColorBox(oak.ScreenWidth/3, secDebugHeight, color.RGBA{2, 2, 100, 100})

		secDebug2.SetPos(float64(oak.ScreenWidth)/3, 0)
		secDebug3.SetPos(2*float64(oak.ScreenWidth)/3, 0)

		render.Draw(secDebug1, layer.UI, 3)
		render.Draw(secDebug2, layer.UI, 3)
		render.Draw(secDebug3, layer.UI, 3)

		pSecDebug := render.NewColorBox(4, 4, color.RGBA{10, 10, 10, 255})
		render.Draw(pSecDebug, layer.UI, 5)
		pSecXNormalizer := sec1.W() * 3 / float64(oak.ScreenWidth)
		pSecYNormalizer := float64(oak.ScreenHeight) / 3 * 2 / float64(secDebugHeight-4)

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
						oak.ShiftScreen(-int(sec1.W())*2, 0)

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

						oak.ShiftScreen(int(sec1.W())*2, 0)

					}
				}
			}
			pSecDebug.SetPos((x-2)/pSecXNormalizer, (pty.Players[0].Y()/pSecYNormalizer)-float64(secDebugHeight)/3-2)
			lastX = x
			return 0
		}, "EnterFrame")

		runbackDisabled := false
		runbackOnce := sync.Once{}

		event.GlobalBind(func(int, interface{}) int {
			if runbackDisabled {
				return 0
			}
			runbackOnce.Do(func() {
				facing = -1
				event.Trigger("RunBack", nil)
				tracker.Prev()
			})
			return event.UnbindSingle
		}, "RunBackOnce")

		endLock := sync.Mutex{}
		defeatedShowing := false

		event.GlobalBind(func(int, interface{}) int {
			endLock.Lock()
			defer endLock.Unlock()
			if pty.Defeated() && !defeatedShowing {
				pty.UnbindAll()
				defeatedShowing = true
				// Show pop up to go to endgame scene
				menuX := (float64(oak.ScreenWidth) - 180) / 2
				menuY := float64(oak.ScreenHeight) / 4
				btn.New(menus.BtnCfgB, btn.Layers(layer.UI, 0),
					btn.Pos(menuX, menuY), btn.Text("Defeated! See Your Stats?"),
					btn.Width(180),
					btn.Binding(mouse.ClickOn, func(int, interface{}) int {
						nextscene = "endGame"
						stayInGame = false

						return 0
					}))
			}
			return 0
		}, "PlayerDeath")

		event.GlobalBind(func(cid int, data interface{}) int {
			dlog.Info("An Enemy Died")

			info := data.([]int64)
			tracker.UpdateHistory(info[0],
				section.Change{
					Typ: section.EntityDestroyed,
					Val: int(info[1])})
			runInfo.EnemiesDefeated += info[0]

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

		bkgMusic, err = music.Start(true, "run2.wav")
		dlog.ErrorCheck(err)

		// Enemy types:
		// 1. Stands Still or walks in a basic path
		// 2. Charges right to left and hurts if you touch it
		// 2a. As 2, but can slightly curve during charge

		// Treasure types:
		// bigger and fancier boxes with fancier colors
		// maybe 3 sizes
		// Color order: brown, then dim gray, then black,
		// white, yellow, green, blue, purple, orange, red, silver

		oak.ResetCommands()
		oak.AddCommand("invuln", func(args []string) {
			dlog.Error("Cheating to set the invulnerability toggle")
			for _, ply := range pty.Players {
				ply.Invulnerable += 100
			}
			if len(args) > 0 && (args[0][0:0] == "f" || args[0][0:0] == "F") {
				for _, ply := range pty.Players {
					ply.Invulnerable -= 100
				}
			}
		})
		oak.AddCommand("stopRunback", func(args []string) {

			runbackDisabled = !runbackDisabled
			dlog.Warn("Cheating to toggle runbackDisabled to:", runbackDisabled)
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
			dlog.Warn("Cheating to toggle speedup party by ", up)
		})
		oak.AddCommand("killme", func(args []string) {
			be := &enemies.BasicEnemy{
				Active: true,
			}
			cid := be.Init()
			sp := collision.NewFullSpace(float64(oak.ViewPos.X), float64(oak.ViewPos.Y),
				1000, 500, labels.Enemy, event.CID(cid))
			collision.Add(sp)
		})
		oak.AddCommand("kill", func(args []string) {
			if len(args) < 1 {
				fmt.Println("Require one argument to kill")
				return
			}
			idx, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Expected integer argument to kill", err)
				return
			}
			x := players.WallOffset + idx * players.PlayerGap

			be := &enemies.BasicEnemy{
				Active: true,
			}
			cid := be.Init()
			sp := collision.NewFullSpace(float64(x+oak.ViewPos.X), float64(oak.ViewPos.Y),
				10, 500, labels.Enemy, event.CID(cid))
			collision.Add(sp)
			timing.DoAfter(30 * time.Millisecond, func() {
				collision.Remove(sp)
			})
		})

		oak.AddCommand("shake", func(args []string) {
			//TODO: determine if default shaker is even noticable

			ss := oak.ScreenShaker{
				Random: true,
				Magnitude: floatgeom.Point2{
					3,
					10,
				},
			}
			ss.Shake(time.Duration(1000) * time.Millisecond)
		})

		oak.AddCommand("grantchest", func(args []string) {
			dlog.Warn("Cheating to grant a chest to a player")

			c := doodads.NewChest(10)
			_, h := c.R.GetDims()
			pty.Players[0].AddChest(h, c.R.(render.Modifiable), c.Value)

		})
		oak.AddCommand("win", func(args []string) {
			nextscene = "endGame"
			stayInGame = false
		})

		dlog.Info("Current Debug Commands are: ", strings.Join(oak.GetDebugKeys(), " , "))
		oak.AddCommand("help", func(args []string) {
			dlog.Info("Current Debug Commands are: ", strings.Join(oak.GetDebugKeys(), " , "))
		})
	},
	Loop: scene.BooleanLoop(&stayInGame),
	End: func() (string, *scene.Result) {
		(*bkgMusic).Stop()
		restrictor.Stop()
		restrictor.Clear()
		return nextscene, &scene.Result{NextSceneInput: Outcome{runInfo}}
	},
}

// Outcome is returned by the run scene
type Outcome struct {
	R records.RunInfo
}
