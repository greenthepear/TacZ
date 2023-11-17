package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) HandleClickControls() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		return
	}
	cx, cy := ebiten.CursorPosition()
	if cx < 0 || cx >= screenWidth || cy < 0 || cy >= screenHeight {
		return
	}
	sx, sy := snapXYtoGrid(generalGridSize, cx, cy)
	sx /= int(generalGridSize)
	sy /= int(generalGridSize)
	if sx == g.TESTpawnX && sy == g.TESTpawnY {
		return
	}
	g.MoveMatrixObjects(1, g.TESTpawnX, g.TESTpawnY, int(sx), int(sy))
	g.TESTpawnX = sx
	g.TESTpawnY = sy

	g.matrixLayers[1].printMatrix()
}
