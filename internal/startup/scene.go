package startup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/run"
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
			characters.Init()
			run.Init()
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

		start := btn.New(menus.BtnCfgA, btn.Color(colornames.Green), btn.Pos(menuX, menuY), btn.Text("Start Game"), bindNewScene("inn"))
		menuY += menus.BtnHeightA * 1.5
		load := btn.New(menus.BtnCfgA, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Load Game"), bindNewScene("load"))
		menuY += menus.BtnHeightA * 1.5
		settings := btn.New(menus.BtnCfgA, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Settings"), bindNewScene("settings"))
		menuY += menus.BtnHeightA * 1.5
		credits := btn.New(menus.BtnCfgA, btn.Color(colornames.Blueviolet), btn.Pos(menuX, menuY), btn.Text("Credits"), bindNewScene("credits"))
		menuY += menus.BtnHeightA * 1.5
		exit := btn.New(menus.BtnCfgA, btn.Pos(menuX, menuY), btn.Text("Exit Game"), btn.Binding(func(int, interface{}) int {
			os.Exit(3)
			return 0
		}))
		// render.Draw()

		fmt.Println("How high are the buttons", start.Y(), load.Y(), settings.Y(), credits.Y(), exit.Y())

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
