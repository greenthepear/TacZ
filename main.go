package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	screenWidth  = 320
	screenHeight = 160
)

var (
	generalGridSize   px = 16
	generalGridWidth     = screenWidth / int(generalGridSize)
	generalGridHeight    = screenHeight / int(generalGridSize)
)

const maxNumberOfLayers = 100

//lint:ignore U1000 Added later
type FreeObjectLayer struct {
	name    string
	z       int
	objects []*GameObject
}

type Game struct {
	freeLayers     []*GameObject
	matrixLayers   []*MatrixLayer
	matrixLayerNum int
	imagePacks     map[string]*ImagePack
	selectedPawn   *GameObject
	playerTurn     bool
	enemies        []*GameObject
	pawns          []*GameObject
}

func (g *Game) Update() error {
	if g.playerTurn {
		g.HandleClickControls()
		g.checkForTurnEndButton()
	} else {
		g.DoEnemyTurn()
	}
	return nil
}

func main() {
	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("TacZ")
	g := &Game{
		imagePacks:     initImagePacks(),
		freeLayers:     make([]*GameObject, maxNumberOfLayers),
		matrixLayers:   make([]*MatrixLayer, maxNumberOfLayers),
		matrixLayerNum: 0,
		selectedPawn:   nil,
		playerTurn:     true,
		enemies:        make([]*GameObject, 0),
		pawns:          make([]*GameObject, 0),
	}

	g.Init()
	g.matrixLayers[0].printMatrix()
	g.matrixLayers[1].printMatrix()

	//g.MoveMatrixObjects(1, 10, 5, 1, 1)
	g.matrixLayers[1].printMatrix()
	/*
		g.CreateNewMatrixLayerOnTop("Ground", generalGridSize, generalGridWidth, generalGridHeight)
		g.SimpleCreateObjectInMatrixLayer(0, "terr", 1, 2, "Terrain")
		g.matrixLayers[0].printMatrix()
	*/
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
