package main

import (
	"fmt"

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

	objWalkable := g.MatrixLayerAtZ(underLayerZ).findObjectWithNameAt(sx, sy, "walkable")
	if g.selectedPawn != nil && objWalkable != nil && !g.MatrixLayerAtZ(boardlayerZ).isOccupied(sx, sy) {
		g.MoveMatrixObjects(boardlayerZ, g.selectedPawn.x, g.selectedPawn.y, sx, sy)
		g.selectedPawn.vars["leftMovement"] -= objWalkable.vars["dist"]
		g.clearMatrixLayer(underLayerZ)
		g.createWalkables(g.findWalkable(g.selectedPawn.x, g.selectedPawn.y, boardlayerZ, int(g.selectedPawn.vars["leftMovement"])), underLayerZ)
		return
	}

	if obj := g.matrixLayers[boardlayerZ].findObjectWithTagAt(sx, sy, "selectable"); obj != nil {
		if g.selectedPawn != obj {
			if g.selectedPawn != nil {
				g.deselectPawn()
				g.clearMatrixLayer(underLayerZ)
			}
			g.selectPawn(obj)
			fmt.Println(g.findWalkable(obj.x, obj.y, boardlayerZ, int(obj.vars["leftMovement"])))
			g.createWalkables(g.findWalkable(obj.x, obj.y, boardlayerZ, int(obj.vars["leftMovement"])), underLayerZ)
		} else {
			g.deselectPawn()
			g.clearMatrixLayer(underLayerZ)
		}
	}
}

func (g *Game) checkForTurnEndButton() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.playerTurn = false
	}
}
