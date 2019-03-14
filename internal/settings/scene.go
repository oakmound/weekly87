package settings

import (
	"image/color"
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/scene"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/mods"
	"github.com/oakmound/weekly87/internal/menus"
)

var (
	SFXVolume     float64 = 1.0
	MusicVolume   float64 = 1.0
	MasterVolume  float64 = 1.0
	ShowFpsToggle bool
)
var (
	stayInMenu bool
)

var (
	musicLevel  = new(float64)
	sfxLevel    = new(float64)
	masterLevel = new(float64)
)

func init() {
	*musicLevel = 1.0
	*sfxLevel = 1.0
	*masterLevel = 1.0
}

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

		volBackground := render.NewVerticalGradientBox(150, 42, color.RGBA{0, 210, 90, 255}, color.RGBA{0, 170, 50, 255})
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

		x := 200.0
		y := 120.0

		infR1 := render.NewVerticalGradientBox(150, 32, color.RGBA{0, 120, 255, 255}, color.RGBA{0, 80, 230, 255})
		infR2 := render.NewCompositeM(infR1, checkMark).ToSprite()

		showFps := btn.And(
			menus.BtnCfgB,
			btn.Toggle(infR2, infR1,
				&ShowFpsToggle),
			btn.Pos(x, y),
			btn.Text("Show FPS"),
		)
		btn.New(showFps)

		sfxVolume := menus.NewSlider(0, x, y+50, 150, 32, 10, 10, nil,
			volBackground.Copy(), 0, 100, 100*(*sfxLevel),
			render.NewColorBox(5, 15, color.RGBA{255, 0, 0, 255}), 1, 1)

		sfxVolume.SetString("SFX Volume")
		sfxVolume.Callback = func(val float64) {
			*sfxLevel = val * 0.01
		}

		musicVolume := menus.NewSlider(0, x, y+100, 150, 32, 10, 10, nil,
			volBackground.Copy(), 0, 100, 100*(*musicLevel),
			render.NewColorBox(5, 15, color.RGBA{255, 0, 0, 255}), 1, 1)

		musicVolume.SetString("Music Volume")
		musicVolume.Callback = func(val float64) {
			*musicLevel = val * 0.01
		}

		masterVolume := menus.NewSlider(0, x, y+150, 150, 32, 10, 10, nil,
			volBackground.Copy(), 0, 100, 100*(*masterLevel),
			render.NewColorBox(5, 15, color.RGBA{255, 0, 0, 255}), 1, 1)

		masterVolume.SetString("Master Volume")
		masterVolume.Callback = func(val float64) {
			*masterLevel = val * 0.01
		}

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) * 3 / 4

		btn.New(menus.BtnCfgB, btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3), btn.Pos(menuX, menuY), btn.Text("Return To Menu"), btn.Binding(func(int, interface{}) int {
			stayInMenu = false
			return 0
		}))
	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End: func() (string, *scene.Result) {
		//if save == nil {
		//	save = &stat.Save{}
		//}
		SFXVolume = *sfxLevel
		MusicVolume = *musicLevel
		MasterVolume = *masterLevel
		//dlog.ErrorCheck(stat.EncodeSave(save, stat.Savefile))
		//dlog.Error("Savefile", save)
		return "startup", nil
	},
}
