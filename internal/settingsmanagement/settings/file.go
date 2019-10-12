package settings

import (
	"encoding/json"
	"os"

	"github.com/oakmound/oak/dlog"
)

const settingsFile = "settings.json"

// Store the settings into a file
func (s *Settings) Store() {
	f, err := os.Create(settingsFile)
	data, err := json.Marshal(s)
	dlog.ErrorCheck(err)
	_, err = f.Write(data)
	dlog.ErrorCheck(err)
	dlog.ErrorCheck(f.Close())
}

// Load the settings from the filesystem
func Load() {
	s := &Settings{}

	f, err := os.Open(settingsFile)
	if err != nil {
		f, err := os.Create(settingsFile)
		dlog.ErrorCheck(err)
		s.SFXVolume = 1.0
		s.MusicVolume = 1.0
		s.MasterVolume = 1.0
		data, err := json.Marshal(s)
		dlog.ErrorCheck(err)
		_, err = f.Write(data)
		dlog.ErrorCheck(err)
	} else {
		dc := json.NewDecoder(f)
		dlog.ErrorCheck(dc.Decode(s))
	}
	if f != nil {
		dlog.ErrorCheck(f.Close())
	}

	Active = *s

}
