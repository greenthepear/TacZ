// Setting up the main gameplay scene
package main

import "math/rand"

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
	boardlayer := g.CreateNewMatrixLayerOnTop("Board", generalGridSize, generalGridWidth, generalGridHeight)
	//g.SimpleCreateObjectInMatrixLayer(boardlayer.z, "Pawn", 10, 5, "Pawn")
	/*
		g.AddObjectToMatrixLayer(NewPawn(g), boardlayer.z, 1, 1)
		g.AddObjectToMatrixLayer(NewPawn(g), boardlayer.z, 2, 2)
		g.AddObjectToMatrixLayer(NewPawn(g), boardlayer.z, 3, 3)
		g.AddObjectToMatrixLayer(NewPawn(g), boardlayer.z, 4, 4)
	*/
	g.AddPawnToLayer(boardlayer.z, 1, 1)
	g.AddPawnToLayer(boardlayer.z, 2, 1)
	g.AddPawnToLayer(boardlayer.z, 3, 1)
	g.AddPawnToLayer(boardlayer.z, 4, 1)

	g.AddEnemyToLayer(boardlayer.z, 9, 5)
	g.AddEnemyToLayer(boardlayer.z, 10, 6)
}
