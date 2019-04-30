package abilities

import (
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/weekly87/internal/characters/enemies"

	"time"
)

var (
	_ Effect = &ParticleEffect{}
)

type Effect interface {
	Create() (cancel func())
	OnEnemy(*enemies.BasicEnemy)
	OnPlayer(User)
	GetPos() (float64, float64)
}

type ParticleEffect struct {
	event.CID
	Start, End  floatgeom.Point2
	Speed       floatgeom.Point2
	Generator   particle.Generator
	BaseLayer   int
	OnPlayerHit func(User)
	OnEnemyHit  func(*enemies.BasicEnemy)
}

func (pe *ParticleEffect) Create() func() {
	// Set reasonable defaults
	if pe.Speed == (floatgeom.Point2{0, 0}) {
		pe.Speed = floatgeom.Point2{1, 1}
	}
	if pe.BaseLayer == 0 {
		pe.BaseLayer = 3
	}
	quit := make(chan struct{})
	// Start the generator
	src := pe.Generator.Generate(pe.BaseLayer)
	// Move the generator from start to end
	go func() {
		tick := time.NewTicker(30 * time.Millisecond)
		for {
			select {
			case <-tick.C:
			case <-quit:
				src.Undraw()
				return
			}
			src.ShiftX(pe.Speed.X())
			src.ShiftY(pe.Speed.Y())
		}
	}()
	return func() {
		close(quit)
	}
}

func (pe *ParticleEffect) OnEnemy(e *enemies.BasicEnemy) {
	pe.OnEnemyHit(e)
}

func (pe *ParticleEffect) OnPlayer(u User) {
	pe.OnPlayerHit(u)
}

func (pe *ParticleEffect) GetPos() (float64, float64) {
	return pe.Generator.GetPos()
}
