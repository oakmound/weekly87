package restrictor

var DefScreen = NewScreen()

func ResetDefault() {
	DefScreen = NewScreen()
}

func Add(r Restrictable) {
	DefScreen.Add(r)
}

func Start(inc int) {
	DefScreen.Start(inc)
}

func Stop() {
	DefScreen.Stop()
}

// Clear destroys all restrictables in the screen
func Clear() {
	DefScreen.Clear()
}
