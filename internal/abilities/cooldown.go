package abilities

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"time"

	"github.com/oakmound/oak/render"
)

type cooldown struct {
	*render.Sprite
	triggeredTime *time.Time
	totalTime     time.Duration
}

// NewCooldown creates a new cooldown
func NewCooldown(w, h int, totalTime time.Duration) *cooldown {
	s := render.NewEmptySprite(0, 0, w, h)

	return &cooldown{s, &time.Time{}, totalTime}
}

// ResetTiming clears the triggered time for a cooldown
func (c *cooldown) ResetTiming() {
	c.triggeredTime = &time.Time{}
}

// Trigger tries to trigger the cooldown and returns whether it was succesful
func (c *cooldown) Trigger() bool {
	if time.Since(*c.triggeredTime) < c.totalTime {
		return false
	}
	// Start the cooldown
	*c.triggeredTime = time.Now()
	return true
}

// Draw the cooldown
func (c *cooldown) Draw(buff draw.Image) {
	c.DrawOffset(buff, 0, 0)
}

// DrawOffset draws the cooldown with the given offset
func (c *cooldown) DrawOffset(buff draw.Image, xOff, yOff float64) {

	if time.Since(*c.triggeredTime) >= c.totalTime {
		return
	}

	cooldownColor := color.RGBA{125, 125, 125, 125}

	percentRecovered := float64(time.Since(*c.triggeredTime)) / float64(c.totalTime)

	w, h := c.GetDims()
	c.Sprite.SetRGBA(image.NewRGBA(image.Rect(0, 0, w, h)))

	recoveredPerimPoints := int((float64(w)*2 + float64(h)*2) * percentRecovered)

	pEvaluated := 0

	centerX := w / 2
	centerY := h / 2

	// O1
	for x := w / 2; x < w; x++ {
		if pEvaluated > recoveredPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, x, 0, cooldownColor)
		}
		pEvaluated++
	}

	// o2  o3
	for y := 0; y < h; y++ {
		if pEvaluated > recoveredPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, w, y, cooldownColor)
		}
		pEvaluated++
	}

	// o4  o5
	for x := w; x > 0; x-- {
		if pEvaluated > recoveredPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, x, h, cooldownColor)
		}
		pEvaluated++
	}

	//o6 o7
	for y := h; y > 0; y-- {

		if pEvaluated > recoveredPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, 0, y, cooldownColor)
		}
		pEvaluated++
	}

	// o8
	for x := 0; x < w/2; x++ {

		if pEvaluated > recoveredPerimPoints {
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
