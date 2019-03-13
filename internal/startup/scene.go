package startup

import (
	"os"
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/enemies"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/run/section"
	"golang.org/x/image/colornames"
)

var stayInMenu bool
var nextscene string

// Scene  to display at startup
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		stayInMenu = true
		nextscene = "startup"
		render.SetDrawStack(
			render.NewCompositeR(),
			render.NewHeap(false),
			render.NewHeap(true),
		)

		if prevScene == "loading" {
			players.Init()
			section.Init()
			enemies.Init()
			// The game has just started, maybe do some
			// intro visual stuff

			//Maybe load? Maybe set seed?
			dlog.Info("Starting game")

		}

		menuBackground, _ := render.LoadSprite("", filepath.Join("raw", "standard_placeholder.png"))
		render.Draw(menuBackground, 0)

		// Render menu buttons
		// 1. Start game
		// 2. Select save file? <- don't worry about saving progress for first build
		// 3. Settings
		// 4. Credits
		// 5. Exit game

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := float64(oak.ScreenHeight) / 4

		btn.New(menus.BtnCfgB, btn.Color(colornames.Green), btn.Pos(menuX, menuY), btn.Text("Start Game"), bindNewScene("inn"))
		menuY += menus.BtnHeightB * 1.5
		btn.New(menus.BtnCfgB, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Load Game"), bindNewScene("load"))
		menuY += menus.BtnHeightB * 1.5
		btn.New(menus.BtnCfgB, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Settings"), bindNewScene("settings"))
		menuY += menus.BtnHeightB * 1.5
		btn.New(menus.BtnCfgB, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Credits"), bindNewScene("credits"))
		menuY += menus.BtnHeightB * 1.5
		btn.New(menus.BtnCfgB, btn.Pos(menuX, menuY), btn.Text("Exit Game"), btn.Binding(func(int, interface{}) int {
			os.Exit(3)
			return 0
		}))
	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End:  scene.GoToPtr(&nextscene),
}

func bindNewScene(newScene string) btn.Option {
	return (btn.Binding(func(int, interface{}) int {
		nextscene = newScene
		stayInMenu = false
		return 0
	}))
}
