package main

import (
	"github.com/oakmound/oak"
	"github.com/oakmound/weekly87/internal/startup"
)

func main() {
	// Add scenes
	oak.AddScene("startup", startup.Scene)
	//oak.AddScene("inn", inn.Scene)
	//oak.AddScene(run, run.Scene)
	oak.Init("startup")
}
