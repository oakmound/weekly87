package layer

import (
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/settings"
)

// Layer constants
const (
	Ground = iota
	Background
	Play
	Effect
	Overlay
	UI
)

// Get returns a slice of all common layers, including fps if active
func Get() []render.Stackable {
	layers := []render.Stackable{
		// ground groundLayer
		render.NewCompositeR(),
		// wall backgroundLayer
		render.NewCompositeR(),
		// entities 	playLayer
		render.NewHeap(false),
		// effects effectLayer
		render.NewHeap(false),
		// overlay Level
		render.NewHeap(false),
		// ui uiLayer
		render.NewHeap(true),
	} 

	if settings.Active.ShowFpsToggle {
		layers = append(layers, render.NewDrawFPS(), render.NewLogicFPS())
	}
	return layers
}
