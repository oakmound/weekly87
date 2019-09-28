package settings

import (
	"fmt"
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak/mouse"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/scene"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/mods"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/menus/selector"
)

var (
	stayInMenu bool
	Active     Settings
)

var (
	musicLevel  = new(float64)
	sfxLevel    = new(float64)
	masterLevel = new(float64)
)

var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInMenu = true
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)

		menuBackground, _ := render.LoadSprite("", filepath.Join("raw", "standard_placeholder.png"))
		render.Draw(menuBackground, 0)

		showFPS := render.NewVerticalGradientBox(100, 20, color.RGBA{0, 210, 90, 255}, color.RGBA{0, 170, 50, 255})
		showFPS.Modify(mod.CutRound(.05, .25),
			mods.Highlight(color.RGBA{170, 170, 170, 200}, 1),
			mods.HighlightOff(color.RGBA{0, 0, 0, 100}, 1, 2, 1))

		sliderWidth := 150.0
		sliderHeight := 42.0
		volBackground := render.NewVerticalGradientBox(int(sliderWidth), int(sliderHeight), color.RGBA{0, 210, 90, 255}, color.RGBA{0, 170, 50, 255})
		volBackground.Modify(mod.CutRound(.05, .25),
			mods.Highlight(color.RGBA{170, 170, 170, 200}, 1),
			mods.HighlightOff(color.RGBA{0, 0, 0, 100}, 1, 2, 1))

		checkMark, err := render.NewPolygon(
			floatgeom.Point2{0, 16},
			floatgeom.Point2{16, 32},
			floatgeom.Point2{32, 0},
			floatgeom.Point2{27, 0},
			floatgeom.Point2{16, 26},
			floatgeom.Point2{0, 11},
		)
		dlog.ErrorCheck(err)
		checkMark.Fill(color.RGBA{100, 255, 100, 255})
		checkMark.ShiftX(110)

		x := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 6
		y := float64(oak.ScreenHeight) / 2.7

		fnt := render.DefFontGenerator.Copy()
		fnt.Color = render.FontColor("Green")
		fnt.Size = 60
		titleFnt := fnt.Generate()

		title := titleFnt.NewStrText("Settings", x-20, y-40)
		render.Draw(title, 2, 12)

		infR1 := render.NewVerticalGradientBox(150, 32, color.RGBA{0, 120, 255, 255}, color.RGBA{0, 80, 230, 255})
		infR2 := render.NewCompositeM(infR1, checkMark).ToSprite()

		showFps := btn.And(
			menus.BtnCfgB,
			btn.Width(150),
			btn.Height(32),
			btn.Toggle(infR2, infR1,
				&Active.ShowFpsToggle),
			btn.Pos(x, y),
			btn.Text("Show FPS"),
		)
		fpsBtn := btn.New(showFps)

		y += 50

		sfxVolume := menus.NewSlider(0, x, y, sliderWidth, sliderHeight, 10, 20, nil,
			volBackground.Copy(), 0, 100, 100*(*sfxLevel),
			render.NewColorBox(5, 15, color.RGBA{255, 0, 0, 255}), 1, 1)

		sfxVolume.SetString("SFX Volume")
		sfxVolume.Callback = func(val float64) {
			*sfxLevel = val * 0.01
		}

		y += 50
		musicVolume := menus.NewSlider(0, x, y, sliderWidth, sliderHeight, 10, 20, nil,
			volBackground.Copy(), 0, 100, 100*(*musicLevel),
			render.NewColorBox(5, 15, color.RGBA{255, 0, 0, 255}), 1, 1)

		musicVolume.SetString("Music Volume")
		musicVolume.Callback = func(val float64) {
			*musicLevel = val * 0.01
		}
		y += 50
		masterVolume := menus.NewSlider(0, x, y, sliderWidth, sliderHeight, 10, 20, nil,
			volBackground.Copy(), 0, 100, 100*(*masterLevel),
			render.NewColorBox(5, 15, color.RGBA{255, 0, 0, 255}), 1, 1)

		masterVolume.SetString("Master Volume")
		masterVolume.Callback = func(val float64) {
			*masterLevel = val * 0.01
		}
		y += 100
		returnBtn := btn.New(menus.BtnCfgB,
			btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
			btn.Pos(x+10, y),
			btn.Text("Return To Menu"),
			btn.Binding(mouse.ClickOn, func(int, interface{}) int {
				stayInMenu = false
				return 0
			}))

		selector.New(
			selector.MouseBindings(true),
			selector.JoystickVertDpadControl(),
			selector.VertArrowControl(),
			selector.Spaces(
				fpsBtn.GetSpace(),
				sfxVolume.Space,
				musicVolume.Space,
				masterVolume.Space,
				returnBtn.GetSpace(),
			),

			selector.SelectTrigger(key.Down+key.Spacebar),
			selector.SelectTrigger("A"+joystick.ButtonUp),

			selector.InteractTrigger(key.Down+key.LeftArrow, -10.0),
			selector.InteractTrigger("Left"+joystick.ButtonUp, -10.0),
			selector.InteractTrigger(key.Down+key.RightArrow, 10.0),
			selector.InteractTrigger("Right"+joystick.ButtonUp, 10.0),

			selector.Callback(func(i int, inc ...interface{}) {
				if len(inc) == 0 {
					if i == 0 {
						fpsBtn.Trigger("MouseClickOn", nil)
					}
					if i == 4 {
						stayInMenu = false
					}
					return
				}
				change, ok := inc[0].(float64)
				if !ok {
					dlog.Error("Expected a float increment")
					return
				}
				switch i {
				case 1:
					sfxVolume.Slide(change)
				case 2:
					musicVolume.Slide(change)
				case 3:
					masterVolume.Slide(change)
				}
			}),
			selector.Layers(2, 20),
		)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End: func() (string, *scene.Result) {

		Active.SFXVolume = *sfxLevel
		Active.MusicVolume = *musicLevel
		Active.MasterVolume = *masterLevel

		fmt.Println(Active)

		Active.Store()

		return "startup", nil
	},
}
