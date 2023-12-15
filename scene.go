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
			gobj := g.SimpleCreateObjectInMatrixLayer(layer.z, "sGround", x, y, "Terrain")
			gobj.sprIdx = rand.Intn(len(gobj.sprites.imagesQ))
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

}
