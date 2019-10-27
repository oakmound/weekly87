package abilities

import (
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

// newCooldown creates a new cooldown
func newCooldown(w, h int, totalTime time.Duration) *cooldown {
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
	// Asset based variables
	cooldownColor := color.RGBA{125, 125, 125, 125}
	w, h := c.GetDims()
	c.Sprite.SetRGBA(image.NewRGBA(image.Rect(0, 0, w, h)))
	centerX := w / 2
	centerY := h / 2
	// Time based variables
	percentRecovered := float64(time.Since(*c.triggeredTime)) / float64(c.totalTime)
	cooldownPerimPoints := int((float64(w)*2 + float64(h)*2) * (1 - percentRecovered))
	pEvaluated := 0

	// Draw each octant as a set of lines to display the cooldown

	// octant 8
	for x := w / 2; x > 0; x-- {
		if pEvaluated < cooldownPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, x, 0, cooldownColor)
		}
		pEvaluated++

	}
	// octants 7 and 6
	if pEvaluated > cooldownPerimPoints {
		goto End
	}
	for y := 0; y < h; y++ {

		if pEvaluated < cooldownPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, 0, y, cooldownColor)
		}
		pEvaluated++
	}

	// octants 5 and 4
	if pEvaluated > cooldownPerimPoints {
		goto End
	}
	for x := 0; x < w; x++ {
		if pEvaluated < cooldownPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, x, h, cooldownColor)
		}
		pEvaluated++
	}

	// octants 3 and 2
	if pEvaluated > cooldownPerimPoints {
		goto End
	}
	for y := h; y > 0; y-- {
		if pEvaluated < cooldownPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, w, y, cooldownColor)
		}
		pEvaluated++
	}

	// octant 1
	if pEvaluated > cooldownPerimPoints {
		goto End
	}
	for x := w; x >= w/2; x-- {
		if pEvaluated < cooldownPerimPoints {
			render.DrawLine(c.GetRGBA(), centerX, centerY, x, 0, cooldownColor)
		}
		pEvaluated++
	}
End:
	c.Sprite.DrawOffset(buff, xOff, yOff)

}

// Copy gets a deep copy of the cooldown
func (c *cooldown) Copy() render.Modifiable {
	return &cooldown{c.Sprite.Copy().(*render.Sprite), c.triggeredTime, c.totalTime}
}
