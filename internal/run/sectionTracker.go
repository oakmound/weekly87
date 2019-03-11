package run

import (
	"image/color"
	"math/rand"

	"github.com/oakmound/weekly87/internal/characters"

	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/render"
)

type SectionTracker struct {
	start        int64
	sectionsDeep int64
	rng          *rand.Rand
	*sectionGenerator
	changes map[int64][]SectionChange
}

func NewSectionTracker(baseSeed int64) *SectionTracker {
	return &SectionTracker{
		start:            baseSeed,
		rng:              rand.New(rand.NewSource(baseSeed)),
		sectionGenerator: &sectionGenerator{},
		changes:          make(map[int64][]SectionChange),
	}
}

// SetDepth is to be used when the chest is picked up, in case there are
// sections that have been generated but they are still off screen.
func (st *SectionTracker) SetDepth(depth int64) {
	if depth < 1 {
		dlog.Error(`SectionTracker cannot have depth set to`, depth, `otherwise
			it could fail to create the starting section.`)
		return
	}
	st.sectionsDeep = depth
}

// ShiftDepth is an alternative to SetDepth in case sections aren't being tracked
// outside of the section tracker, but the game knows how many sections ahead
// it has generated from the current one
func (st *SectionTracker) ShiftDepth(depth int64) {
	st.SetDepth(st.sectionsDeep + depth)
}

// Next produces another section.
func (st *SectionTracker) Next() *Section {
	st.sectionsDeep++
	st.rng.Seed(st.start + st.sectionsDeep)
	// This following section is test code
	// These initial rng calls should make these test sections more distinct
	rLimit := st.rng.Intn(255)
	gLimit := st.rng.Intn(255)
	bLimit := st.rng.Intn(255)
	for x := 0; x < len(st.ground); x++ {
		for y := 0; y < len(st.ground[x]); y++ {
			st.ground[x][y] = render.NewColorBox(
				16, 16, color.RGBA{
					uint8(st.rng.Intn(rLimit)),
					uint8(st.rng.Intn(gLimit)),
					uint8(st.rng.Intn(bLimit)),
					255},
			)
		}
	}
	rLimit = st.rng.Intn(255)
	gLimit = st.rng.Intn(255)
	bLimit = st.rng.Intn(255)
	for x := 0; x < len(st.wall); x++ {
		for y := 0; y < len(st.wall[x]); y++ {
			st.wall[x][y] = render.NewColorBox(
				16, 16, color.RGBA{
					uint8(st.rng.Intn(rLimit)),
					uint8(st.rng.Intn(gLimit)),
					uint8(st.rng.Intn(bLimit)),
					255},
			)
		}
	}
	testEC := &characters.EnemyConstructor{
		Position:   floatgeom.Point2{400, 400},
		Dimensions: floatgeom.Point2{32, 32},
		Speed:      floatgeom.Point2{-3 * rand.Float64(), -5 * (rand.Float64() - .5)},
		AnimationMap: map[string]render.Modifiable{
			"standRT": render.NewColorBox(32, 32, color.RGBA{255, 125, 0, 255}),
			"standLT": render.NewColorBox(32, 32, color.RGBA{125, 255, 0, 255}),
			"walkRT":  render.NewColorBox(32, 32, color.RGBA{0, 125, 255, 255}),
			"walkLT":  render.NewColorBox(32, 32, color.RGBA{255, 255, 255, 255}),
		},
	}
	e, err := testEC.NewEnemy()
	if err != nil {
		dlog.Error(err)
	} else {
		st.entities = append(st.entities, e)
	}
	return st.generate()
}
