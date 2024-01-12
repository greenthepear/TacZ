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

type Game struct {
	freeLayers   []*FreeObjectLayer
	matrixLayers []*MatrixLayer

	imagePacks   map[string]*ImagePack
	selectedPawn *GameObject
	playerTurn   bool
	enemies      []*GameObject
	pawns        []*GameObject
}

func (g *Game) Update() error {
	if g.playerTurn {
		g.HandleClickControls()
		g.checkForTurnEndButton()
	} else {
		g.deselectPawn()
		g.clearMatrixLayer(underLayerZ)
		g.DoEnemyTurn()
	}
	return nil
}

func main() {
	g := &Game{
		imagePacks:   initImagePacks(),
		freeLayers:   make([]*FreeObjectLayer, 0),
		matrixLayers: make([]*MatrixLayer, 0),
		selectedPawn: nil,
		playerTurn:   true,
		enemies:      make([]*GameObject, 0),
		pawns:        make([]*GameObject, 0),
	}

	g.Init()

	ebiten.SetWindowSize(screenWidth*4, screenHeight*4)
	ebiten.SetWindowTitle("TacZ")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
