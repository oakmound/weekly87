package abilities

import (
	"fmt"
	"image/color"
	"image/draw"
	"time"

	"github.com/oakmound/oak/render"
)

type cooldown struct {
	*render.Sprite
	triggeredTime time.Time
	totalTime     time.Duration
}

// NewCooldown creates a new cooldown
func NewCooldown(w, h int, totalTime time.Duration) *cooldown {
	s := render.NewEmptySprite(0, 0, w, h)

	return &cooldown{s, time.Time{}, totalTime}
}

func (c *cooldown) Trigger() bool {
	fmt.Println("time is at ", c.triggeredTime)
	if time.Since(c.triggeredTime) < c.totalTime {
		return false
	}
	// Start the cooldown
	c.triggeredTime = time.Now()
	fmt.Println("triggered at ", c.triggeredTime)
	return true
}

func (c *cooldown) Draw(buff draw.Image) {
	c.DrawOffset(buff, 0, 0)
}

func (c *cooldown) DrawOffset(buff draw.Image, xOff, yOff float64) {

	if time.Since(c.triggeredTime) >= c.totalTime {
		fmt.Println("Triggered at ", c.triggeredTime, "which has been ", time.Since(c.triggeredTime), c.totalTime)
		return
	}
	fmt.Println("Drawing a cooldown")

	cooldownColor := color.RGBA{155, 155, 155, 255}

	percentRecovered := 1 - (time.Since(c.triggeredTime) / c.totalTime)

	w, h := c.GetDims()
	recoveredPerimPoints := (w*2 + h*2) * int(percentRecovered)

	pEvaluated := 0

	centerX := w / 2
	centerY := h / 2

	// O1
	for x := w / 2; x < w; x++ {
		if pEvaluated > recoveredPerimPoints {
			fmt.Println("Drawingt Line")
			render.DrawLine(c.GetRGBA(), centerX, centerY, x, 0, cooldownColor)
		}
		pEvaluated++
	}

	// o2  o3
	for y := 0; y < h; y++ {
		if pEvaluated > recoveredPerimPoints {
			fmt.Println("Drawingt Line")
			render.DrawLine(c.GetRGBA(), centerX, centerY, w, y, cooldownColor)
		}
		pEvaluated++
	}

	// o4  o5
	for x := w; w > 0; x-- {
		if pEvaluated > recoveredPerimPoints {
			fmt.Println("Drawingt Line")
			render.DrawLine(c.GetRGBA(), centerX, centerY, x, h, cooldownColor)
		}
		pEvaluated++
	}

	//o6 o7
	for y := h; y > 0; y-- {

		if pEvaluated > recoveredPerimPoints {
			fmt.Println("Drawingt Line")
			render.DrawLine(c.GetRGBA(), centerX, centerY, 0, y, cooldownColor)
		}
		pEvaluated++
	}

	// o8
	for x := 0; x < w/2; x++ {

		if pEvaluated > recoveredPerimPoints {
			fmt.Println("Drawingt Line")
			render.DrawLine(c.GetRGBA(), centerX, centerY, x, 0, cooldownColor)
		}
		pEvaluated++
	}
	fmt.Println("Eval ", pEvaluated, " with remeainingPerc", recoveredPerimPoints)

	c.Sprite.DrawOffset(buff, xOff, yOff)

}

func (c *cooldown) Copy() render.Modifiable {
	return &cooldown{c.Sprite.Copy().(*render.Sprite), c.triggeredTime, c.totalTime}
}
