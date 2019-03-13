package menus

// This was stolen from agent
// This should be made easier to use and moved to btn?

import (
	"image/color"
	"time"

	"github.com/oakmound/oak/entities/x/btn"

	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
)

// Slider is a little UI element with a movable slider
// representing a range. It has a specific purpose of
// setting visualization delay, but eventually it might
// be expanded to a more generalized structure.
type Slider struct {
	btn.TextBox
	min, val, max float64
	interval      time.Duration
	knub          render.Renderable
	knubLine      render.Renderable
	sliding       bool
	textY         float64
	Callback      func(float64)
}

// Init returns a CID for the Slider.
//
// Note on engine internals:
// All entities as defined by the engine
// need to have this function defined on them.
// this is because an entity is only a meaningful
// concept in terms of the engine for an entity's
// ability to have events bound to it and triggered
// on it, which this CID (caller ID) represents.
//
// Its literal meaning is, in our event bus, the value
// passed into NextID (which is the only way to get a
// legitimate CID), is stored at the array index of the
// returned CID.
func (sl *Slider) Init() event.CID {
	sl.CID = event.NextID(sl)
	return sl.CID
}

// NewSlider returns a slider with initialized values
// using the given font to render its text.
func NewSlider(cid event.CID, x, y, w, h, txtX, txtY float64,
	f *render.Font, r render.Renderable, min, max, defval float64,
	knub render.Renderable, layers ...int) *Slider {

	sl := new(Slider)
	if cid == 0 {
		cid = sl.Init()
	}
	sl.TextBox = *btn.NewTextBox(cid, x, y, w, h, txtX, txtY, f, r, layers...)
	sl.min = min
	sl.max = max
	sl.val = defval
	sl.textY = txtY

	sl.CID.Bind(sliderDragStart, "MousePressOn")

	sl.knub = knub
	sl.knubLine = render.NewLine(0, 0, max, 0, color.RGBA{255, 255, 255, 255})

	knubLayers := make([]int, len(layers))
	copy(knubLayers, layers)
	knubLayers[len(knubLayers)-1] += 2
	render.Draw(sl.knub, knubLayers...)
	knubLayers[len(knubLayers)-1]--
	render.Draw(sl.knubLine, knubLayers...)

	sl.SetPos(x, y)
	return sl
}

func (sl *Slider) SetPos(x float64, y float64) {
	sl.SetLogicPos(x, y)
	if sl.R != nil {
		sl.R.SetPos(x, y)
	}

	rwidth, rheight := sl.knub.GetDims()

	if sl.knub != nil {
		sl.knub.SetPos(x+float64(rwidth)+sl.val, y+float64(rheight)+sl.textY)
	}
	if sl.knubLine != nil {
		sl.knubLine.SetPos(x+float64(rwidth)+5, y+float64(rheight)+sl.textY*2)
	}

	if sl.Space != nil {
		mouse.UpdateSpace(sl.X(), sl.Y(), sl.W, sl.H, sl.Space)
	}
}

// sliderDragStart tells this demo to ignore some mouse events
// until sliding = false, and binds to every frame sliderDrag.
func sliderDragStart(sl int, nothing interface{}) int {
	slider := event.GetEntity(sl).(*Slider)
	if slider.sliding != true {
		slider.sliding = true
		slider.Bind(sliderDrag, "EnterFrame")
	}
	return 0
}

// sliderDrag updates the position and value of this slider's
// knub, within a defined range. Once the mouse is let go,
// it allows other mouse operations to resume and updates
// the visualizaton delay to the value it was left at.
func sliderDrag(sl int, nothing interface{}) int {
	slider := event.GetEntity(sl).(*Slider)
	me := mouse.LastEvent
	if me.Event == "MouseRelease" {
		slider.sliding = false
		if slider.Callback != nil {
			slider.Callback(slider.val)
		}
		return event.UnbindEvent
	}
	w, _ := slider.knub.GetDims()
	x := float64(me.X()) - (slider.X() + float64(w))
	if x <= slider.min {
		if slider.val == slider.min {
			return 0
		}
		slider.val = slider.min
	} else if x >= slider.max {
		if slider.val == slider.max {
			return 0
		}
		slider.val = slider.max
	} else {
		slider.val = x
	}
	slider.knub.SetPos(slider.X()+slider.val, slider.knub.Y())
	return 0
}
