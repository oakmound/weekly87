package joys

import (
	"math"
	"sync"
	"time"
	"strings"

	"github.com/oakmound/oak/joystick"
	"github.com/oakmound/oak/event"
)

var (
	joyStickStates    map[uint32]joystick.State = make(map[uint32]joystick.State)
	joyStickStateLock sync.RWMutex
)

func LowestID() uint32 {
	lowestID := uint32(math.MaxInt32)
	joyStickStateLock.RLock()
	for id := range joyStickStates {
		if id < lowestID {
			lowestID = id
		}
	}
	joyStickStateLock.RUnlock()
	return lowestID
}

func StickState(v uint32) joystick.State {
	joyStickStateLock.RLock()
	st := joyStickStates[v]
	joyStickStateLock.RUnlock()
	return st
}

func SetStickState(k uint32, v joystick.State) {
	joyStickStateLock.Lock()
	joyStickStates[k] = v
	joyStickStateLock.Unlock()
}

type handler struct{}

func (h *handler) Trigger(ev string, state interface{}) {
	if ev == joystick.Disconnected {
		id, ok := state.(uint32)
		if ok {
			SetStickState(id, joystick.State{})
		}
		return
	}
	if strings.HasSuffix(ev, joystick.ButtonUp) {
		event.Trigger(ev, state)
	}
	st, ok := state.(*joystick.State)
	if !ok {
		return
	}
	SetStickState(st.ID, *st)
}

var initOnce = sync.Once{}

func Init() {
	initOnce.Do(func() {
		joystick.Init()

		go func() {
			jCh, _ := joystick.WaitForJoysticks(3 * time.Second)
			for {
				select {
				case j := <-jCh:
					j.Handler = &handler{}
					j.Listen(&joystick.ListenOptions{
						JoystickChanges: true,
						ButtonPresses:   true,
					})
				}
			}
		}()
	})
}
