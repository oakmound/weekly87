package abilities

import (
	"fmt"
	"image/color"
	"time"

	"github.com/oakmound/oak/render"
	"github.com/oakmound/weekly87/internal/characters"
)

var (
	//Fireball tries to cast a magical fire ball in front of the mage
	Fireball = NewAbility(
		render.NewColorBox(64, 64, color.RGBA{200, 10, 0, 255}),
		time.Second*10,
		func(u User) []characters.Character {
			fmt.Println("Just tried to burn a guy ", u)
			return nil
		},
	)
	// Invulnerability = NewAbility(
	// 	render.NewColorBox(64, 64, color.RGBA{255, 255, 125, 255}),
	// 	time.Second*10,
	// 	func(u User) {
	// 		c := Constructor{}
	// 		err := c.StartAt(u.Pos()).
	// 			ArcTo(pos2).
	// 			WithParticles(ps).
	// 			//WithRenderable().
	// 			//WithLabel().
	// 			ThenDrop(otherThing).
	// 			Fire()
	// 		dlog.ErrorCheck(err)

	// 		//c.StartAt(u.Pos()).LineTo(pos2)
	// 	},
	// )
)
