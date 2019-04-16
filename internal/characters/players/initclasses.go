package players

import (
	"strings"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/mod"
)

func Init() {
	WarriorsInit()
	MageInit()
}
func filterCharMap(baseCharMap map[string]render.Modifiable, filter mod.Filter) map[string]render.Modifiable {
	outputMap := make(map[string]render.Modifiable)

	for k, v := range baseCharMap {
		outputMap[k] = v.Copy()
		if !strings.Contains(k, "dead") {
			outputMap[k].Filter(filter)
		}
	}

	return outputMap
}
