package section

import (
	"math/rand"
	"path/filepath"

	"github.com/200sc/go-dist/intrange"

	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak/render"

	"github.com/oakmound/weekly87/internal/characters/enemies"
)

type entityPlan struct {
	chestCount        intrange.Range
	chestRange        intrange.Range
	enemyCount        intrange.Range
	enemyDistribution [enemies.TypeLimit]float64
	enemyVariantRange intrange.Range
}

type tilePlan struct {
	groundTiles  []render.Modifiable
	skyTiles     []render.Modifiable
	surfaceTiles []render.Modifiable
}

type tileWeight struct {
	groundTileWeights  []float64
	skyTileWeights     []float64
	surfaceTileWeights []float64
}

type sectionPlan struct {
	tilePlan
	tileWeight
	entityPlan
	effects []render.Modifiable
}

func (sp *sectionPlan) setRng(rng *rand.Rand) {
	sp.entityPlan.chestCount.SetRand(rng)
	sp.entityPlan.chestRange.SetRand(rng)
	sp.entityPlan.enemyCount.SetRand(rng)
	sp.entityPlan.enemyVariantRange.SetRand(rng)

}

var tilePlans = map[string]tilePlan{}
var tileWeights = [13]tileWeight{}
var sectionPlans [13]sectionPlan

func Init() {
	dir := filepath.Join("assets", "images")
	wallSheet, err := render.LoadSprites(dir, filepath.Join("16x16", "walltiles.png"), 16, 16, 0)
	dlog.ErrorCheck(err)
	groundSheet, err := render.LoadSprites(dir, filepath.Join("16x16", "floortiles.png"), 16, 16, 0)
	dlog.ErrorCheck(err)

	entityPlanA := entityPlan{
		chestCount: intrange.NewLinear(0, 5),
		chestRange: intrange.NewLinear(1, 5),
		enemyCount: intrange.NewLinear(4, 9),
		enemyDistribution: [...]float64{
			enemies.Hare:   .5,
			enemies.Mantis: .5,
			enemies.Tree:   1,
		},
		enemyVariantRange: intrange.NewLinear(0, enemies.VariantCount-1),
	}

	aPlanWeight := tileWeight{
		groundTileWeights: []float64{
			5, 5, 5, 5,
			1, 1, 1, 1,
		},
		skyTileWeights: []float64{
			9, 9,
			1, 1,
		},
		surfaceTileWeights: []float64{
			9, 9,
			1, 1,
		},
	}
	bPlanWeight := tileWeight{
		groundTileWeights: []float64{
			1, 1, 1, 1,
			6, 6, 6, 6,
			1, 1, 1, 1,
		},
		skyTileWeights: []float64{
			9, 9,
			1, 1,
		},
		surfaceTileWeights: []float64{
			9, 9,
			1, 1,
		},
	}
	cPlanWeight := tileWeight{
		groundTileWeights: []float64{
			1, 1, 1, 1,
			6, 6, 6, 6,
			1, 1, 1, 1,
		},
		skyTileWeights: []float64{
			5, 5, 5, 5,
			1, 1, 1, 1,
		},
		surfaceTileWeights: []float64{
			5, 5, 5, 5,
			1, 1, 1, 1,
		},
	}
	dPlanWeight := tileWeight{
		groundTileWeights: []float64{
			5, 5, 5, 5,
			1, 1, 1, 1,
		},
		skyTileWeights: []float64{
			5, 5, 5, 5,
			1, 1, 1, 1,
		},
		surfaceTileWeights: []float64{
			5, 5, 5, 5,
			1, 1, 1, 1,
		},
	}

	tileWeights = [13]tileWeight{
		aPlanWeight,
		bPlanWeight,
		cPlanWeight,
		dPlanWeight,
		cPlanWeight,
		bPlanWeight,
		aPlanWeight,
		bPlanWeight,
		cPlanWeight,
		dPlanWeight,
		cPlanWeight,
		bPlanWeight,
		aPlanWeight,
	}

	tilePlans["A"] = tilePlan{
		groundTiles: []render.Modifiable{
			groundSheet[0][0].Copy(),
			groundSheet[0][1].Copy(),
			groundSheet[1][0].Copy(),
			groundSheet[1][1].Copy(),
			groundSheet[2][0].Copy(),
			groundSheet[2][1].Copy(),
			groundSheet[3][0].Copy(),
			groundSheet[3][1].Copy(),
		},
		skyTiles: []render.Modifiable{
			wallSheet[0][0].Copy(),
			wallSheet[1][0].Copy(),
			wallSheet[2][0].Copy(),
			wallSheet[3][0].Copy(),
		},
		surfaceTiles: []render.Modifiable{
			wallSheet[0][1].Copy(),
			wallSheet[1][1].Copy(),
			wallSheet[2][1].Copy(),
			wallSheet[3][1].Copy(),
		},
	}
	tilePlans["B"] = tilePlan{
		groundTiles: []render.Modifiable{
			groundSheet[0][0].Copy(),
			groundSheet[0][1].Copy(),
			groundSheet[1][0].Copy(),
			groundSheet[1][1].Copy(),
			groundSheet[2][0].Copy(),
			groundSheet[2][1].Copy(),
			groundSheet[3][0].Copy(),
			groundSheet[3][1].Copy(),
			groundSheet[0][2].Copy(),
			groundSheet[0][3].Copy(),
			groundSheet[1][2].Copy(),
			groundSheet[1][3].Copy(),
		},
		skyTiles: []render.Modifiable{
			wallSheet[2][0].Copy(),
			wallSheet[3][0].Copy(),
			wallSheet[0][0].Copy(),
			wallSheet[1][0].Copy(),
		},
		surfaceTiles: []render.Modifiable{
			wallSheet[2][1].Copy(),
			wallSheet[3][1].Copy(),
			wallSheet[0][1].Copy(),
			wallSheet[1][1].Copy(),
		},
	}
	tilePlans["C"] = tilePlan{
		groundTiles: []render.Modifiable{
			groundSheet[0][0].Copy(),
			groundSheet[0][1].Copy(),
			groundSheet[1][0].Copy(),
			groundSheet[1][1].Copy(),
			groundSheet[0][2].Copy(),
			groundSheet[0][3].Copy(),
			groundSheet[1][2].Copy(),
			groundSheet[1][3].Copy(),
			groundSheet[2][2].Copy(),
			groundSheet[2][3].Copy(),
			groundSheet[3][2].Copy(),
			groundSheet[3][3].Copy(),
		},
		skyTiles: []render.Modifiable{
			wallSheet[0][2].Copy(),
			wallSheet[1][2].Copy(),
			wallSheet[0][3].Copy(),
			wallSheet[1][3].Copy(),
			wallSheet[2][2].Copy(),
			wallSheet[3][2].Copy(),
			wallSheet[2][3].Copy(),
			wallSheet[3][3].Copy(),
		},
		surfaceTiles: []render.Modifiable{
			wallSheet[0][2].Copy(),
			wallSheet[1][2].Copy(),
			wallSheet[0][3].Copy(),
			wallSheet[1][3].Copy(),
			wallSheet[2][2].Copy(),
			wallSheet[3][2].Copy(),
			wallSheet[2][3].Copy(),
			wallSheet[3][3].Copy(),
		},
	}
	tilePlans["D"] = tilePlan{
		groundTiles: []render.Modifiable{
			groundSheet[0][2].Copy(),
			groundSheet[0][3].Copy(),
			groundSheet[1][2].Copy(),
			groundSheet[1][3].Copy(),
			groundSheet[2][2].Copy(),
			groundSheet[2][3].Copy(),
			groundSheet[3][2].Copy(),
			groundSheet[3][3].Copy(),
		},
		skyTiles: []render.Modifiable{
			wallSheet[2][2].Copy(),
			wallSheet[3][2].Copy(),
			wallSheet[2][3].Copy(),
			wallSheet[3][3].Copy(),
			wallSheet[0][2].Copy(),
			wallSheet[1][2].Copy(),
			wallSheet[0][3].Copy(),
			wallSheet[1][3].Copy(),
		},
		surfaceTiles: []render.Modifiable{
			wallSheet[2][2].Copy(),
			wallSheet[3][2].Copy(),
			wallSheet[2][3].Copy(),
			wallSheet[3][3].Copy(),
			wallSheet[0][2].Copy(),
			wallSheet[1][2].Copy(),
			wallSheet[0][3].Copy(),
			wallSheet[1][3].Copy(),
		},
	}

	sectionPlans = [13]sectionPlan{
		{
			tilePlan:   tilePlans["A"],
			tileWeight: tileWeights[0],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["B"],
			tileWeight: tileWeights[1],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["C"],
			tileWeight: tileWeights[2],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["D"],
			tileWeight: tileWeights[3],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["C"],
			tileWeight: tileWeights[4],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["B"],
			tileWeight: tileWeights[5],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["A"],
			tileWeight: tileWeights[6],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["B"],
			tileWeight: tileWeights[7],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["C"],
			tileWeight: tileWeights[8],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["D"],
			tileWeight: tileWeights[9],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["C"],
			tileWeight: tileWeights[10],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["B"],
			tileWeight: tileWeights[11],
			entityPlan: entityPlanA,
		},
		{
			tilePlan:   tilePlans["A"],
			tileWeight: tileWeights[12],
			entityPlan: entityPlanA,
		},
	}
}
