package abilities

import (
	"fmt"
	"image/color"
	"time"

	"github.com/oakmound/oak/render"
)

var (
	//Fireball tries to cast a magical fire ball in front of the mage
	Fireball = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 255}),
		time.Second*10,
		func(u User) { fmt.Println("Just tried to burn a guy ", u) },
	)
)
