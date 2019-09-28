package sfx

import (
	"github.com/200sc/klangsynthese/audio/filter"
	"github.com/200sc/klangsynthese/font"
	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/dlog"
)

// Sfx and audio fonts
var (
	LoudSFX = font.New()
	SoftSFX = font.New()
	Music   = font.New()
	Audios  = map[string]*audio.Audio{}
	inited  bool
)

// Init the sfx files
func Init() {
	LoudSFX.Filter(filter.Volume(.5))
	SoftSFX.Filter(filter.Volume(.25))
	Music.Filter(filter.Volume(.4), filter.LoopOn())

	files := map[string]*font.Font{

		"playerHit1":    LoudSFX,
		"stormEffect":   LoudSFX,
		"bannerPlaced1": SoftSFX,
		"slashHeavy":    SoftSFX,
		"slashLight":    SoftSFX,
		"bounced1":      LoudSFX,
		"warriorCast1":  LoudSFX,
	}
	for s, f := range files {
		a, err := audio.Get(s + ".wav")
		if err != nil {
			dlog.Error(err)

			continue
		}
		Audios[s] = audio.New(f, a)
	}
	dlog.Info("SFX load completed")
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
