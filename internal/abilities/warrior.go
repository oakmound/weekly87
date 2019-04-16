package abilities

import (
	"fmt"
	"image/color"
	"time"

	"github.com/oakmound/oak/render"
)

var (
	SpearStab = &ability{
		renderable: render.NewColorBox(64, 64, color.RGBA{200, 200, 0, 255}),
		cooldown:   time.Second * 5,
		trigger:    func(u User) { fmt.Println("Just tried to stab a guy ", u) },
	}
)
