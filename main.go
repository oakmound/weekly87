package main

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/weekly87/internal/assets"
	"github.com/oakmound/weekly87/internal/credits"
	"github.com/oakmound/weekly87/internal/end"
	"github.com/oakmound/weekly87/internal/history"
	"github.com/oakmound/weekly87/internal/inn"
	"github.com/oakmound/weekly87/internal/run"
	"github.com/oakmound/weekly87/internal/settingsmanagement"
	"github.com/oakmound/weekly87/internal/startup"
)

func main() {
	oak.SetupConfig = oak.Config{
		Screen: oak.Screen{
			Width:  1024,
			Height: 576,
		},
		Title:     "Hiring Now: Looters",
		BatchLoad: true,
		Debug: oak.Debug{

			Level: "INFO",
		},
		TrackInputChanges: true,
	}
	oak.SetBinaryPayload(assets.Asset, assets.AssetDir)

	oak.SetupTopMost = true

	// Add scenes
	oak.AddScene("startup", startup.Scene)
	oak.AddScene("inn", inn.Scene)
	oak.AddScene("settings", settingsmanagement.Scene)
	oak.AddScene("credits", credits.Scene)
	oak.AddScene("history", history.Scene)
	oak.AddScene("run", run.Scene)
	oak.AddScene("endGame", end.Scene) // At end if there is time break this into its own package and export the correct stats
	oak.Init("startup")
}
