package joys

import (
	"time"

	"github.com/oakmound/oak/joystick"
)

var (
	JoyStickStates map[uint32]joystick.State = make(map[uint32]joystick.State)
)

type handler struct{}

func (h *handler) Trigger(ev string, state interface{}) {
	if ev == joystick.Disconnected {
		id, ok := state.(uint32)
		if ok {
			JoyStickStates[id] = joystick.State{}
		}
		return
	}
	st, ok := state.(*joystick.State)
	if !ok {
		return
	}
	JoyStickStates[st.ID] = *st
}

func Init() {
	joystick.Init()

	go func() {
		jCh, _ := joystick.WaitForJoysticks(3 * time.Second)
		for {
			select {
			case j := <-jCh:
				j.Handler = &handler{}
				j.Listen(nil)
			}
		}
	}()
}
