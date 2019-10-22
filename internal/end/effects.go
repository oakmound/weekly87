package end

import (
	"image/color"
	"math"
	"path/filepath"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/layer"
	"github.com/oakmound/weekly87/internal/sfx"
)

func deathSprites(x, y float64) {
	ptGenLife := intrange.NewLinear(40, 60)
	ptColor := color.RGBA{200, 200, 200, 255}
	ptColorRand := color.RGBA{0, 0, 0, 0}
	newPf := floatrange.NewLinear(5, 10)
	ptLife := floatrange.NewLinear(100, 200)
	angle := floatrange.NewLinear(0, 360)
	speed := floatrange.NewLinear(1, 4)
	size := intrange.Constant(3)
	layerFn := func(v physics.Vector) int {
		return layer.Effect
	}
	particle.NewColorGenerator(
		particle.Pos(x, y),
		particle.Duration(ptGenLife),
		particle.LifeSpan(ptLife),
		particle.Angle(angle),
		particle.Speed(speed),
		particle.Layer(layerFn),
		particle.Shape(shape.Square),
		particle.Size(size),
		particle.Color(ptColor, ptColorRand, ptColor, ptColorRand),
		particle.NewPerFrame(newPf)).Generate(0)
}

func explodeChest(x, y float64) {
	sp, err := render.GetSprite(filepath.Join("raw", "wood_junk.png"))
	if err != nil {
		dlog.Error(err)
		return
	}
	sfx.Play("chestExplode")
	explodeSprite(x, y, sp)

}
func explodeSprite(x, y float64, sprite *render.Sprite) {
	layerFn := func(v physics.Vector) int {
		return layer.Effect
	}
	ptGenLife := intrange.NewLinear(40, 60)
	sg := particle.NewSpriteGenerator(
		particle.NewPerFrame(floatrange.NewSpread(8, 0)),
		particle.Pos(x, y),
		particle.LifeSpan(floatrange.NewSpread(20, 5)),
		particle.Angle(floatrange.NewSpread(0, 360)),
		particle.Speed(floatrange.NewSpread(2, .5)),
		particle.Spread(3, 2),
		particle.Duration(ptGenLife),
		particle.Layer(layerFn),
		particle.Sprite(sprite),
		particle.SpriteRotation(floatrange.Constant(10)),
	)
	sg.Generate(layer.Effect)
}

// makeGoldParticles creates the appropriate amount of collision particles within the given location
func makeGoldParticles(goldCount int, location floatgeom.Rect2) {
	debug := collision.NewRect2Space(location, 0)
	debug.UpdateLabel(collision.Label(labels.Ornament))
	collision.Add(debug)

	center := location.Center()

	//TODO: make this an actual fxn probably making it a log of goldCount
	particleCount := int(math.Log(float64(10.0 * goldCount)))

	colorOpts := particle.And(
		particle.NewPerFrame(floatrange.NewConstant(float64(particleCount))),
		particle.Limit(particleCount),
		particle.InfiniteLifeSpan(),
		particle.Spread(location.W()/2+8, location.H()/2),
		particle.Shape(shape.Diamond),
		particle.Size(intrange.NewConstant(4)),
		particle.Speed(floatrange.NewConstant(0)),
		particle.Pos(center.X(), center.Y()),
		particle.Color(color.RGBA{200, 200, 0, 255}, color.RGBA{0, 0, 0, 0},
			color.RGBA{200, 200, 0, 255}, color.RGBA{0, 0, 0, 0}),
	)
	shiftFactor := floatrange.NewLinear(0, 6)
	pg := particle.NewCollisionGenerator(
		particle.NewColorGenerator(colorOpts),
		particle.Fragile(false),
		particle.HitMap(map[collision.Label]collision.OnHit{
			labels.PC: func(a, b *collision.Space) {
				// b.CID.Trigger("Attacked", hitEffects)
				p, ok := event.GetEntity(int(b.CID)).(*entities.Interactive)
				if !ok {
					dlog.Error("A non player is colliding with gold?")
					return
				}
				goldPiece := particle.Lookup(int(a.CID))
				d := p.Delta.Copy().Scale(shiftFactor.Poll())
				goldPiece.ShiftX(d.X())
				goldPiece.ShiftY(d.Y())
			},
		}),
	)
	pg.Generate(layer.Play)

}
