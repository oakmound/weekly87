package main

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/weekly87/internal/inn"
	"github.com/oakmound/weekly87/internal/run"
	"github.com/oakmound/weekly87/internal/settings"
	"github.com/oakmound/weekly87/internal/startup"
)

func main() {
	oak.SetupConfig = oak.Config{
		Screen: oak.Screen{
			Width:  1024,
			Height: 576,
		},
	}
	// Add scenes
	oak.AddScene("startup", startup.Scene)
	oak.AddScene("inn", inn.Scene)
	oak.AddScene("settings", settings.Scene)
	oak.AddScene("run", run.Scene)
	oak.Init("startup")
}
