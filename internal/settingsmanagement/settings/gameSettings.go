package settings

// Settings serves as our safe file and all variables we track across
// multiple runs
type Settings struct {
	SFXVolume     float64 `json:"sfxVolume"`
	MusicVolume   float64 `json:"musicVolume"`
	MasterVolume  float64 `json:"masterVolume"`
	ShowFpsToggle bool    `json:"showFpsToggle"`
	Debug         bool    `json:"debugOn,omitempty"`
}

// Active settings for the game
var Active Settings
