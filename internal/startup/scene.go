package startup

import (
	"os"
	"path/filepath"

	"github.com/oakmound/oak/mouse"

	"golang.org/x/image/colornames"

	"github.com/oakmound/weekly87/internal/abilities"
	"github.com/oakmound/weekly87/internal/joys"
	"github.com/oakmound/weekly87/internal/menus/selector"
	"github.com/oakmound/weekly87/internal/run"
	"github.com/oakmound/weekly87/internal/sfx"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/entities/x/btn/grid"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/enemies"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/records"
	"github.com/oakmound/weekly87/internal/run/section"
	"github.com/oakmound/weekly87/internal/settingsmanagement/settings"
)

var stayInMenu bool
var nextscene string
var saveHistory records.Records

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
			saveHistory := records.Load()
			run.BaseSeed = saveHistory.BaseSeed
			joys.Init()
			abilities.Init()
			players.Init()
			section.Init()
			enemies.Init()
			settings.Load()
			sfx.Init()

			// The game has just started, maybe do some
			// intro visual stuff

			dlog.Info("Starting game")
		}
		if prevScene == "settings" {
			sfx.UpdateLevels()
		}
		// Todo: joystick mouse

		menuBackground, _ := render.LoadSprite("", filepath.Join("raw", "standard_placeholder.png"))
		render.Draw(menuBackground, 0)

		looters, _ := render.LoadSprite("", filepath.Join("raw", "title_card.png"))
		looters.SetPos(30, 60)
		render.Draw(looters, 0)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 6
		menuY := float64(oak.ScreenHeight) / 2.7

		// Render menu buttons
		// 1. Start game
		// 2. Select save file? <- don't worry about saving progress for first build
		// 3. Settings
		// 4. Credits
		// 5. Exit game
		//get the title

		selectors := grid.New(
			grid.Defaults(btn.And(menus.BtnCfgB, btn.Pos(menuX, menuY))),
			grid.YGap(menus.BtnHeightB*1.5),
			grid.Content(
				[][]btn.Option{
					{
						btn.And(btn.Color(colornames.Green), btn.Text("Start Game"), bindNewScene("inn")),
						// btn.And(btn.Color(colornames.Blueviolet), btn.Text("Load Game"), bindNewScene("load")),
						btn.And(btn.Color(colornames.Gold), btn.Text("Game History"), bindNewScene("history")),
						btn.And(btn.Color(colornames.Blue), btn.Text("Settings"), bindNewScene("settings")),
						btn.And(btn.Color(colornames.Blueviolet), btn.Text("Credits"), bindNewScene("credits")),
						btn.And(btn.Text("Exit Game"), btn.Binding(mouse.ClickOn, func(int, interface{}) int {
							os.Exit(3)
							return 0
						})),
					},
				},
			),
		)

		selector.New(
			menus.ButtonSelectorA(selectors),
			selector.MouseBindings(true),
		)

	},
	Loop: scene.BooleanLoop(&stayInMenu),
	End:  scene.GoToPtr(&nextscene),
}

func bindNewScene(newScene string) btn.Option {
	return btn.Binding(mouse.ClickOn, func(int, interface{}) int {
		nextscene = newScene
		stayInMenu = false
		return 0
	})
}
