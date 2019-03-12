package run

import (
	"path/filepath"

	"github.com/oakmound/oak/dlog"

	"github.com/oakmound/oak/render"
)

type sectionPlanTile struct {
	groundTiles  []render.Modifiable
	skyTiles     []render.Modifiable
	surfaceTiles []render.Modifiable
}

type sectionPlanWeight struct {
	groundTileWeights  []float64
	skyTileWeights     []float64
	surfaceTileWeights []float64
}

type sectionPlan struct {
	sectionPlanTile
	sectionPlanWeight
	effects []render.Modifiable
}

var sectionPlanTiles = map[string]sectionPlanTile{}
var sectionPlanWeights = [13]sectionPlanWeight{}
var sectionPlans [13]sectionPlan

func Init() {
	dir := filepath.Join("assets", "images")
	wallSheet, err := render.LoadSprites(dir, filepath.Join("16x16", "wallTiles.png"), 16, 16, 0)
	dlog.ErrorCheck(err)
	groundSheet, err := render.LoadSprites(dir, filepath.Join("16x16", "floorTiles.png"), 16, 16, 0)
	dlog.ErrorCheck(err)

	aPlanWeight := sectionPlanWeight{
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
	bPlanWeight := sectionPlanWeight{
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
	cPlanWeight := sectionPlanWeight{
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
	dPlanWeight := sectionPlanWeight{
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

	sectionPlanWeights = [13]sectionPlanWeight{
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

	sectionPlanTiles["A"] = sectionPlanTile{
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
	sectionPlanTiles["B"] = sectionPlanTile{
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
	sectionPlanTiles["C"] = sectionPlanTile{
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
	sectionPlanTiles["D"] = sectionPlanTile{
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
			sectionPlanTile:   sectionPlanTiles["A"],
			sectionPlanWeight: sectionPlanWeights[0],
		},
		{
			sectionPlanTile:   sectionPlanTiles["B"],
			sectionPlanWeight: sectionPlanWeights[1],
		},
		{
			sectionPlanTile:   sectionPlanTiles["C"],
			sectionPlanWeight: sectionPlanWeights[2],
		},
		{
			sectionPlanTile:   sectionPlanTiles["D"],
			sectionPlanWeight: sectionPlanWeights[3],
		},
		{
			sectionPlanTile:   sectionPlanTiles["C"],
			sectionPlanWeight: sectionPlanWeights[4],
		},
		{
			sectionPlanTile:   sectionPlanTiles["B"],
			sectionPlanWeight: sectionPlanWeights[5],
		},
		{
			sectionPlanTile:   sectionPlanTiles["A"],
			sectionPlanWeight: sectionPlanWeights[6],
		},
		{
			sectionPlanTile:   sectionPlanTiles["B"],
			sectionPlanWeight: sectionPlanWeights[7],
		},
		{
			sectionPlanTile:   sectionPlanTiles["C"],
			sectionPlanWeight: sectionPlanWeights[8],
		},
		{
			sectionPlanTile:   sectionPlanTiles["D"],
			sectionPlanWeight: sectionPlanWeights[9],
		},
		{
			sectionPlanTile:   sectionPlanTiles["C"],
			sectionPlanWeight: sectionPlanWeights[10],
		},
		{
			sectionPlanTile:   sectionPlanTiles["B"],
			sectionPlanWeight: sectionPlanWeights[11],
		},
		{
			sectionPlanTile:   sectionPlanTiles["A"],
			sectionPlanWeight: sectionPlanWeights[12],
		},
	}
}
