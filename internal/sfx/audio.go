package sfx

import (
	"github.com/200sc/klangsynthese/audio/filter"
	"github.com/200sc/klangsynthese/font"
	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/weekly87/internal/settingsmanagement/settings"
)

// Sfx and audio fonts
var (
	loudSFX, softSFX *font.Font
	Audios           = map[string]*audio.Audio{}

	inited bool
)

// Init the sfx files
func Init() {
	loadAudio()
	dlog.Info("SFX load completed")
}

// UpdateLevels to play sfx at
// TODO: (currently re loads from files shouldnt)
func UpdateLevels() {
	loadAudio()
	dlog.Info("SFX volume adjusting completed")
}

func loadAudio() {
	loudSFX = font.New()
	softSFX = font.New()
	loudSFX.Filter(filter.Volume(.5 * settings.Active.SFXVolume * settings.Active.MasterVolume))
	softSFX.Filter(filter.Volume(.25 * settings.Active.SFXVolume * settings.Active.MasterVolume))

	files := map[string]*font.Font{

		"playerHit1":    loudSFX,
		"stormEffect":   loudSFX,
		"bannerPlaced1": softSFX,
		"slashHeavy":    softSFX,
		"slashLight":    softSFX,
		"bounced1":      loudSFX,
		"warriorCast1":  loudSFX,
		"mageCast1":     loudSFX,
		"fireball1":     loudSFX,
		"chestHop1":     loudSFX,
		"ohWell":        loudSFX,
		"chestExplode":  softSFX,
		"dissappear1":   softSFX,
		"selected":      loudSFX,
		"cooldown":      loudSFX,
		"nope1":         loudSFX,
		"abilityReady1": softSFX,
	}
	for s, f := range files {
		a, err := audio.Get(s + ".wav")
		if err != nil {
			dlog.Error(err)

			continue
		}
		Audios[s] = audio.New(f, a)
	}
}

// Play a copy of the sfx requested
func Play(s string) {
	audOrigin, ok := Audios[s]
	if !ok {
		dlog.Error("Tried to play unloaded Audio:", s)
		return
	}
	aud, err := audOrigin.Copy()
	if err != nil {
		dlog.Error(err)
		return
	}
	a := aud.(*audio.Audio)
	a.Play()
}
