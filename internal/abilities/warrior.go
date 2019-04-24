package abilities

import (
	"fmt"
	"image/color"
	"time"

	"github.com/oakmound/oak/render"
)

var (
	SpearStab = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{1, 1, 0, 1}), //color.RGBA{200, 200, 0, 255}
		time.Second*5,
		func(u User) { fmt.Println("Just tried to stab a guy ", u) },
	)
)
