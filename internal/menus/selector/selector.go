package selector

import (
	"errors"
	"image/color"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/key"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/render"
)

type Option func(*Constructor)

func Spaces(sps ...*collision.Space) Option {
	return func(sc *Constructor) {
		sc.Spots = make([]floatgeom.Rect2, len(sps))
		for i, s := range sps {
			sc.Spots[i] = floatgeom.NewRect2WH(s.X(), s.Y(), s.W(), s.H())
		}
	}
}

// Todo: more spot constructors

// And combines a variadic number of options
func And(opts ...Option) Option {
	return func(sc *Constructor) {
		for _, opt := range opts {
			opt(sc)
		}
	}
}

// Control sets the given strings to control movement in the selector
func Control(prev, next string) Option {
	return func(sc *Constructor) {
		if sc.Bindings == nil {
			sc.Bindings = make(map[string]func(*Selector))
		}
		sc.Bindings[prev] = func(s *Selector) {
			s.MoveTo(s.Pos - 1)
		}
		sc.Bindings[next] = func(s *Selector) {
			s.MoveTo(s.Pos + 1)
		}
	}
}

// VertWASDControl sets W and S as selection controls
func VertWASDControl() Option {
	return Control(key.Down+key.W, key.Down+key.S)
}

// VertArrowControl sets uparrow and downarrow as selection controls
func VertArrowControl() Option {
	return Control(key.Down+key.UpArrow, key.Down+key.DownArrow)
}

// HorzWASDControl sets A and D as selection controls
func HorzWASDControl() Option {
	return Control(key.Down+key.A, key.Down+key.D)
}

// HorzArrowControl sets leftarrow and r as selection controls
func HorzArrowControl() Option {
	return Control(key.Down+key.LeftArrow, key.Down+key.RightArrow)
}

func JoystickHorzDpadControl() Option {
	return Control("Left"+joystick.ButtonUp, "Right"+joystick.ButtonUp)
}

func JoystickVertDpadControl() Option {
	return Control("Up"+joystick.ButtonUp, "Down"+joystick.ButtonUp)
}

// Layers sets the layer for drawing the selector
func Layers(lys ...int) Option {
	return func(sc *Constructor) {
		sc.Layers = lys
	}
}

// Callback determines what shoud happen on a select event
func Callback(cb func(i int, data ...interface{})) Option {
	return func(sc *Constructor) {
		sc.Callback = cb
	}
}

// Cleanup determines what shoud happen on a select event
func Cleanup(cb func(i int)) Option {
	return func(sc *Constructor) {
		sc.Cleanup = cb
	}
}

// Display sets how to display the selector
func Display(display func(floatgeom.Point2) render.Renderable) Option {
	return func(sc *Constructor) {
		sc.Display = display
	}
}

// Wraps sets whether selection should wrap around
func Wraps(shouldWrap bool) Option {
	return func(sc *Constructor) {
		sc.Wraps = true
	}
}

// SelectTrigger sets the input/event to trigger selection with
func SelectTrigger(trigger string) Option {
	return func(sc *Constructor) {
		if sc.Bindings == nil {
			sc.Bindings = make(map[string]func(*Selector))
		}
		sc.Bindings[trigger] = func(s *Selector) {
			s.Select()
		}
	}
}

// InteractTrigger sets the input/event to trigger an interaction
func InteractTrigger(trigger string, data ...interface{}) Option {
	return func(sc *Constructor) {

		if sc.Bindings == nil {
			sc.Bindings = make(map[string]func(*Selector))
		}
		sc.Bindings[trigger] = func(s *Selector) {
			s.Interact(data...)
		}
	}
}

// DestroyTrigger sets the input/event to destroy the selector
func DestroyTrigger(trigger string) Option {
	return func(sc *Constructor) {
		if sc.Bindings == nil {
			sc.Bindings = make(map[string]func(*Selector))
		}
		sc.Bindings[trigger] = func(s *Selector) {
			s.Destroy()
		}
	}
}

// MouseBindings sets whether selection options should react to mouseover and clicking
func MouseBindings(on bool) Option {
	return func(sc *Constructor) {
		if !on {
			sc.MouseLeft = nil
			sc.MouseRight = nil
			return
		}

		if sc.MouseLeft == nil {
			sc.MouseLeft = defaultMouseSelect

		}
		if sc.MouseRight == nil {
			sc.MouseRight = defaultMouseSelect
		}
	}
}

var defaultMouseSelect = func(s *Selector, i int) int {
	s.MoveTo(i)
	s.Select()
	return 0
}

func MouseLeft(op func(s *Selector, i int) int) Option {
	return func(sc *Constructor) {
		sc.MouseLeft = op
	}
}
func MouseRight(op func(s *Selector, i int) int) Option {
	return func(sc *Constructor) {
		sc.MouseRight = op
	}
}

func MouseInteract(data ...interface{}) func(s *Selector, i int) int {
	return func(s *Selector, i int) int {
		dlog.Info("Interacting with ", i)
		s.MoveTo(i)
		s.Interact(data...)
		return 0
	}
}

type mousePhaser struct {
	id event.CID
	mouse.CollisionPhase
}

func (mp *mousePhaser) Init() event.CID {
	mp.id = event.NextID(mp)
	return mp.id
}

type Constructor struct {
	Display func(size floatgeom.Point2) render.Renderable
	// Step/Limit/Size should be option A; Slice of floatgeom.Rect2 should be option B;
	// StepFn/Size should be option C;
	Spots      []floatgeom.Rect2
	Layers     []int
	Callback   func(i int, data ...interface{})
	Cleanup    func(i int)
	MouseLeft  func(*Selector, int) int
	MouseRight func(*Selector, int) int
	Bindings   map[string]func(*Selector)
	Wraps      bool
}

type Selector struct {
	event.CID
	*Constructor
	Pos int
	R   render.Renderable
}

func New(opts ...Option) (*Selector, error) {
	c := &Constructor{}
	for _, opt := range opts {
		opt(c)
	}
	return c.Generate()
}

func (sc *Constructor) Generate() (*Selector, error) {
	if sc.Display == nil {
		sc.Display = func(size floatgeom.Point2) render.Renderable {
			return render.NewColorBox(int(size.X()), int(size.Y()), color.RGBA{125, 125, 0, 125})
		}
	}

	s := &Selector{
		Constructor: sc,
	}
	s.CID = s.Init()
	for ev, bnd := range sc.Bindings {
		bnd := bnd
		s.Bind(func(id int, _ interface{}) int {
			s, ok := event.GetEntity(id).(*Selector)
			if !ok {
				dlog.Error("Failed to get selector from id")
				return 0
			}
			bnd(s)
			return 0
		}, ev)
	}

	mouseSpaces := []*collision.Space{}
	if sc.MouseLeft != nil || sc.MouseRight != nil {
		// Make collisionPhase objects for each spot
		for i, spt := range sc.Spots {
			i := i
			mp := &mousePhaser{}
			mp.Init()
			sp := collision.NewSpace(spt.Min.X(), spt.Min.Y(), spt.W(), spt.H(), mp.id)
			mouseSpaces = append(mouseSpaces, sp)
			mouse.Add(sp)
			mouse.PhaseCollision(sp)
			mp.id.Bind(func(int, interface{}) int {
				s.MoveTo(i)
				return 0
			}, mouse.Start)
			mp.id.Bind(func(_ int, data interface{}) int {
				m := data.(mouse.Event)
				switch m.Button {
				case "LeftMouse":
					return sc.MouseLeft(s, i)
				case "RightMouse":
					return sc.MouseRight(s, i)
				}
				return 0
			}, mouse.ClickOn)
		}
	}

	s.MoveTo(0)
	oldCleanup := s.Cleanup
	s.Cleanup = func(i int) {
		for _, ms := range mouseSpaces {
			mouse.Remove(ms)
		}
		if oldCleanup != nil {
			oldCleanup(i)
		}

	}

	return s, nil
}

func (s *Selector) Init() event.CID {
	return event.NextID(s)
}

func (s *Selector) MoveTo(i int) error {
	if i < 0 || i >= len(s.Spots) {
		if !s.Wraps {
			return errors.New("Index to move to exceeds limit")
		}
		if i < 0 {
			i = len(s.Spots) - 1
		} else {
			i = 0
		}
	}
	spot := s.Spots[i]
	draw := false
	if s.R == nil || s.Spots[s.Pos].W() != spot.W() || s.Spots[s.Pos].H() != spot.H() {
		draw = true
		if s.R != nil {
			s.R.Undraw()
		}
		s.R = s.Display(spot.Max.Sub(spot.Min))
	}

	s.R.SetPos(spot.Min.X(), spot.Min.Y())
	if draw {
		render.Draw(s.R, s.Layers...)
	}
	s.Pos = i
	return nil
}

func (s *Selector) Interact(data ...interface{}) {
	if s.Callback != nil {
		s.Callback(s.Pos, data...)
	}
}

func (s *Selector) Select() {
	if s.Callback != nil {
		s.Callback(s.Pos)
	}
	if s.Cleanup != nil {
		s.Cleanup(s.Pos)
	}
}

func (s *Selector) Destroy() {
	if s.Cleanup != nil {
		s.Cleanup(s.Pos)
	}
	event.DestroyEntity(int(s.CID))
	s.R.Undraw()
}
