package layer

import (
	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/settingsmanagement/settings"
)

// Layer constants
const (
	Ground = iota
	Background
	Play
	Effect
	Overlay
	UI
	Debug
)

// Get returns a slice of all common layers, including fps if active
func Get() []render.Stackable {
	layers := []render.Stackable{
		// ground groundLayer
		render.NewCompositeR(),
		// wall backgroundLayer
		render.NewCompositeR(),
		// entities 	playLayer
		render.NewDynamicHeap(),
		// effects effectLayer
		render.NewDynamicHeap(),
		// overlay Level
		render.NewDynamicHeap(),
		// ui uiLayer
		render.NewStaticHeap(),
		// debug
		render.NewStaticHeap(),
	}

	if settings.Active.ShowFpsToggle {
		layers = append(layers, render.NewDrawFPS(), render.NewLogicFPS())
	}
	return layers
}
