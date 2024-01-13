// Setting up the main gameplay scene
package main

import "math/rand"

const (
	backgroundLayerZ = 0
	underLayerZ      = 1
	boardlayerZ      = 2
)

const (
	pawnInfoLayerZ = 0
)

func (g *Game) InitBackgroundLayer() {
	layer := g.CreateNewMatrixLayerOnTop("Background", generalGridSize, generalGridWidth, generalGridHeight, 0, 0)
	for y := 0; y < layer.height; y++ {
		for x := 0; x < layer.width; x++ {
			gobj := g.SimpleCreateObjectInMatrixLayer(layer.z, "sGround", x, y, "Terrain", false)
			gobj.sprIdx = rand.Intn(len(gobj.sprites.imagesQ))
		}
	}
}

func (g *Game) InitObstacles(randomObstacleNum int) {
	type fence struct {
		x, y   int
		suffix string
	}

	makeFence := func(f fence) {
		gobj := g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", f.x, f.y, "Obstacles/Fences", true)
		gobj.sprKey = "Fences/fence" + f.suffix
	}

	fences := [...]fence{
		{13, 0, "UpLeft"}, {13, 2, "UpLeft"}, {13, 3, "EndLeft"},
		{14, 3, "EndMid"}, {15, 3, "2"}, {16, 3, "3"}, {17, 3, "EndRight"}, {19, 3, "EndLeft"},
	}

	for _, f := range fences {
		makeFence(f)
	}

	//Random grave stones
	for i := 0; i < randomObstacleNum; i++ {
		randx := rand.Intn(generalGridWidth)
		randy := rand.Intn(generalGridHeight)
		if !g.MatrixLayerAtZ(boardlayerZ).isOccupied(randx, randy) {
			g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "gravestone", randx, randy, "Obstacles", false)
		}
	}
}

func (g *Game) Init() {
	g.InitBackgroundLayer()
	g.CreateNewMatrixLayerOnTop("Under", generalGridSize, generalGridWidth, generalGridHeight, 0, 0)
	g.CreateNewMatrixLayerOnTop("Board", generalGridSize, generalGridWidth, generalGridHeight, 0, 0)

	g.AddPawnToLayer(boardlayerZ, 1, 1)
	g.AddPawnToLayer(boardlayerZ, 2, 1)
	g.AddPawnToLayer(boardlayerZ, 3, 1)
	g.AddPawnToLayer(boardlayerZ, 4, 1)

	g.AddEnemyToLayer(boardlayerZ, 9, 5)
	g.AddEnemyToLayer(boardlayerZ, 10, 6)

	g.InitObstacles(10)

	g.CreateNewFreeLayerOnTop("pawnInfoLayer")
}
