package section

import (
	"path/filepath"

	"github.com/200sc/go-dist/floatrange"

	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak/render"

	"github.com/oakmound/weekly87/internal/characters/enemies"
)

type entityPlan struct {
	chestCount        int
	chestRange        floatrange.Range
	enemyCount        int
	enemyDistribution [enemies.TypeLimit]float64
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

var tilePlans = map[string]tilePlan{}
var tileWeights = [13]tileWeight{}
var sectionPlans [13]sectionPlan

func Init() {
	dir := filepath.Join("assets", "images")
	wallSheet, err := render.LoadSprites(dir, filepath.Join("16x16", "wallTiles.png"), 16, 16, 0)
	dlog.ErrorCheck(err)
	groundSheet, err := render.LoadSprites(dir, filepath.Join("16x16", "floorTiles.png"), 16, 16, 0)
	dlog.ErrorCheck(err)

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
		},
		{
			tilePlan:   tilePlans["B"],
			tileWeight: tileWeights[1],
		},
		{
			tilePlan:   tilePlans["C"],
			tileWeight: tileWeights[2],
		},
		{
			tilePlan:   tilePlans["D"],
			tileWeight: tileWeights[3],
		},
		{
			tilePlan:   tilePlans["C"],
			tileWeight: tileWeights[4],
		},
		{
			tilePlan:   tilePlans["B"],
			tileWeight: tileWeights[5],
		},
		{
			tilePlan:   tilePlans["A"],
			tileWeight: tileWeights[6],
		},
		{
			tilePlan:   tilePlans["B"],
			tileWeight: tileWeights[7],
		},
		{
			tilePlan:   tilePlans["C"],
			tileWeight: tileWeights[8],
		},
		{
			tilePlan:   tilePlans["D"],
			tileWeight: tileWeights[9],
		},
		{
			tilePlan:   tilePlans["C"],
			tileWeight: tileWeights[10],
		},
		{
			tilePlan:   tilePlans["B"],
			tileWeight: tileWeights[11],
		},
		{
			tilePlan:   tilePlans["A"],
			tileWeight: tileWeights[12],
		},
	}
}
