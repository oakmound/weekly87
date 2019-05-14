package restrictor

import (
	"fmt"
	"sync"

	"github.com/oakmound/oak"

	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
)

type Screen struct {
	event.CID
	runningLeft bool
	sliceLock   sync.Mutex
	restricted  []Restrictable
	inc         int
}

func (s *Screen) Init() event.CID {
	s.CID = event.NextID(s)
	return s.CID
}

func NewScreen() *Screen {
	s := &Screen{}
	s.Init()
	s.Bind(screenFlip, "RunBack")
	s.restricted = []Restrictable{}
	return s
}

func (s *Screen) Add(r Restrictable) {
	s.sliceLock.Lock()
	s.restricted = append(s.restricted, r)
	s.sliceLock.Unlock()
}

func (s *Screen) Start(inc int) {
	s.inc = inc
	s.Bind(screenEnter, "EnterFrame")
}

func (s *Screen) Stop() {
	event.UnbindAll(
		event.BindingOption{
			Event: event.Event{
				Name:     "EnterFrame",
				CallerID: int(s.CID),
			},
		},
	)
}

// Clear destroys all restrictables in the screen
func (s *Screen) Clear() {
	s.sliceLock.Lock()
	s.restricted = []Restrictable{}
	s.sliceLock.Unlock()
}

func screenFlip(id int, frame interface{}) int {
	s, ok := event.GetEntity(id).(*Screen)
	if !ok {
		dlog.Warn("Screen function triggered on non-screen")
		return event.UnbindSingle
	}
	s.runningLeft = true
	return 0
}

func screenEnter(id int, frame interface{}) int {
	f, ok := frame.(int)
	if !ok {
		dlog.Error("EnterFrame frame was not int")
		return event.UnbindSingle
	}

	s, ok := event.GetEntity(id).(*Screen)
	if !ok {
		dlog.Error("Screen function triggered on non-screen")
	}

	// If performance problems arise, switch to the new slice
	// model that the render heap uses

	i := f % s.inc
	for {
		s.sliceLock.Lock()
		if i >= len(s.restricted) {
			s.sliceLock.Unlock()
			break
		}
		res := s.restricted[i]
		s.sliceLock.Unlock()

		xf, _ := res.GetPos()
		x := int(xf)
		w, _ := res.GetDims()

		offScreen := false
		if !s.runningLeft {
			if x+w < oak.ViewPos.X {
				fmt.Println("Offscreen A", x, w, oak.ViewPos.X)
				offScreen = true
			}
		} else {
			if x > oak.ViewPos.X+oak.ScreenWidth {
				fmt.Println("Offscreen B", x, oak.ScreenWidth, oak.ViewPos.X)
				offScreen = true
			}
		}
		if offScreen {
			res.Destroy()
			s.sliceLock.Lock()
			copy(s.restricted[i:], s.restricted[i+1:])
			s.restricted[len(s.restricted)-1] = nil
			s.restricted = s.restricted[:len(s.restricted)-1]
			s.sliceLock.Unlock()
		}
		i += s.inc
	}
	return 0
}

type Restrictable interface {
	Destroy()
	GetPos() (float64, float64)
	GetDims() (int, int)
}

/*   A                B                 C
 S E C T I O N 1  | S E C T I O N 2 | S E C T I O N 1 |
		z	  			   y	              x
*/

// Going Right
// When player reaches y:


// When player reaches x:


// Going Left
// When player reaches z:


// When player reaches y:

