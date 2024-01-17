// Setting up the main gameplay scene
package main

import "math/rand"

// Matrix layers Zs
const (
	backgroundLayerZ = iota //background ground
	underEnemyLayerZ        //enemy spawned attackables
	underLayerZ             //player spawned attackables and walkables

	boardlayerZ //Pawns, enemy objects, obstacles

	underAttacksLayerZ //UI, for signifying chosen attack
	attacksLayerZ      //Attack selection
	emptyTopLayerZ     //Testing and potentially pathfinding
)

// Free layers Zs
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
			o := g.SimpleCreateObjectInMatrixLayer(boardlayerZ, "Gravestone", randx, randy, "Obstacles", false)
			o.tags = []string{"damageable", "obstacle"}
			o.vars["leftHP"] = 1.0
		}
	}
}

func (g *Game) InitLayers() {
	g.InitBackgroundLayer()
	g.CreateNewMatrixLayerOnTop("EnemyUnder", generalGridSize, generalGridWidth, generalGridHeight, 0, 0)
	g.CreateNewMatrixLayerOnTop("Under", generalGridSize, generalGridWidth, generalGridHeight, 0, 0)
	g.CreateNewMatrixLayerOnTop("Board", generalGridSize, generalGridWidth, generalGridHeight, 0, 0)

	g.AddPawnToLayer(boardlayerZ, 1, 1)
	g.AddPawnToLayer(boardlayerZ, 2, 1)
	g.AddPawnToLayer(boardlayerZ, 3, 1)
	g.AddPawnToLayer(boardlayerZ, 4, 1)

	g.AddSkinnyToLayer(boardlayerZ, 9, 5)
	g.AddSkinnyToLayer(boardlayerZ, 10, 6)
	g.AddSkinnyToLayer(boardlayerZ, 11, 6)
	g.AddSkinnyToLayer(boardlayerZ, 8, 6)

	g.InitObstacles(10)

	g.InitUILayers()
	g.CreateNewMatrixLayerOnTop("EmptyTop", generalGridSize, generalGridWidth, generalGridHeight, 0, 0)

}
