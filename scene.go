// Setting up the main gameplay scene
package main

import "math/rand"

const (
	backgroundLayerZ = 0
	underLayerZ      = 1
	boardlayerZ      = 2
)

func (g *Game) InitBackgroundLayer() {
	layer := g.CreateNewMatrixLayerOnTop("Background", generalGridSize, generalGridWidth, generalGridHeight)
	for y := 0; y < layer.height; y++ {
		for x := 0; x < layer.width; x++ {
			gobj := g.SimpleCreateObjectInMatrixLayer(layer.z, "sGround", x, y, "Terrain", false)
			gobj.sprIdx = rand.Intn(len(gobj.sprites.imagesQ))
		}
	}
}

func (g *Game) InitObstacles(randomObstacleNum int) {
	gobj := g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", 13, 0, "Obstacles/Fences", true)
	gobj.sprKey = "Fences/fenceUpLeft"
	gobj = g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", 13, 2, "Obstacles/Fences", true)
	gobj.sprKey = "Fences/fenceUpLeft"
	gobj = g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", 13, 3, "Obstacles/Fences", true)
	gobj.sprKey = "Fences/fenceEndLeft"

	gobj = g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", 14, 3, "Obstacles/Fences", true)
	gobj.sprKey = "Fences/fenceEndMid"
	gobj = g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", 15, 3, "Obstacles/Fences", true)
	gobj.sprKey = "Fences/fence2"
	gobj = g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", 16, 3, "Obstacles/Fences", true)
	gobj.sprKey = "Fences/fence3"
	gobj = g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", 17, 3, "Obstacles/Fences", true)
	gobj.sprKey = "Fences/fenceEndRight"

	gobj = g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "fence", 19, 3, "Obstacles/Fences", true)
	gobj.sprKey = "Fences/fenceEndLeft"

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
	g.CreateNewMatrixLayerOnTop("Under", generalGridSize, generalGridWidth, generalGridHeight)
	g.CreateNewMatrixLayerOnTop("Board", generalGridSize, generalGridWidth, generalGridHeight)

	g.AddPawnToLayer(boardlayerZ, 1, 1)
	g.AddPawnToLayer(boardlayerZ, 2, 1)
	g.AddPawnToLayer(boardlayerZ, 3, 1)
	g.AddPawnToLayer(boardlayerZ, 4, 1)

	g.AddEnemyToLayer(boardlayerZ, 9, 5)
	g.AddEnemyToLayer(boardlayerZ, 10, 6)

	g.InitObstacles(10)

	g.CreateNewFreeLayerOnTop("freeLayerTest")
	g.AddObjectToFreeLayer(0,
		NewGameObject(
			"test", 60, 80, g.imagePacks["UI"], false, 1, "UI/walkable", true, g, nil, nil, nil, nil),
	)
}
