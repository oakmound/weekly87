package music

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/oakmound/oak/audio"
	"github.com/oakmound/oak/dlog"

	klg "github.com/200sc/klangsynthese/audio"
	"github.com/200sc/klangsynthese/audio/filter"

	"github.com/oakmound/weekly87/internal/settings"
)

// Start plays the audio tracks in order, optionally looping the final track
func Start(loop bool, names ...string) (*klg.Audio, error) {
	if len(names) == 0 {
		return nil, errors.New("No names provided")
	}

	var music klg.Audio
	var err error
	music, err = audio.Load(filepath.Join("assets", "audio"), names[0])
	if err != nil {
		return nil, err
	}
	music, err = music.Copy()
	if err != nil {
		return nil, err
	}
	music = music.MustFilter(
		filter.Volume(0.5 * settings.Active.MusicVolume * settings.Active.MasterVolume),
	)

	music.Play()

	if len(names) > 0 {
		go func() {
			for i := 1; i < len(names); i++ {
				time.Sleep(music.PlayLength())
				music, err = audio.Load(filepath.Join("assets", "audio"), names[i])
				dlog.ErrorCheck(err)
				music, err = music.Copy()
				dlog.ErrorCheck(err)
				music = music.MustFilter(
					filter.Volume(0.5 * settings.Active.MusicVolume * settings.Active.MasterVolume),
				)
				if loop && i == len(names)-1 {
					music = music.MustFilter(
						filter.LoopOn(),
					)
				}
				music.Play()
			}
		}()
	}

	return &music, nil
}
