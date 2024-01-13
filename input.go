package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) IsXYWithingMatrixLayerBounds(layerZ, x, y int) bool {
	l := g.MatrixLayerAtZ(layerZ)
	return x >= int(l.xOffset) &&
		x < l.width*int(l.squareLength)+int(l.xOffset) &&
		y >= int(l.yOffset) &&
		y < l.height*int(l.squareLength)+int(l.yOffset)
}

func (g *Game) CursorXYtoMatrixGrid(layerZ, sx, sy int) (int, int) {
	l := g.MatrixLayerAtZ(layerZ)
	return (sx - int(l.xOffset)) / int(l.squareLength),
		(sy - int(l.yOffset)) / int(l.squareLength)
}

func (g *Game) HandleBoardSelection(sx, sy int) {
	sx, sy = g.CursorXYtoMatrixGrid(boardlayerZ, sx, sy)
	objWalkable := g.MatrixLayerAtZ(underLayerZ).findObjectWithNameAt(sx, sy, "walkable")
	if g.selectedPawn != nil && objWalkable != nil && !g.MatrixLayerAtZ(boardlayerZ).isOccupied(sx, sy) {
		g.MoveMatrixObjects(boardlayerZ, g.selectedPawn.x, g.selectedPawn.y, sx, sy)
		g.selectedPawn.vars["leftMovement"] -= objWalkable.vars["dist"]
		g.ClearMatrixLayer(underLayerZ)
		g.CreateWalkables(g.FindWalkable(g.selectedPawn.x, g.selectedPawn.y, boardlayerZ, int(g.selectedPawn.vars["leftMovement"])), underLayerZ)
		return
	}

	if obj := g.matrixLayers[boardlayerZ].findObjectWithTagAt(sx, sy, "selectable"); obj != nil {
		if g.selectedPawn != obj {
			if g.selectedPawn != nil {
				g.DeselectPawn()
				g.ClearMatrixLayer(underLayerZ)
			}
			g.SelectPawn(obj)
			//fmt.Println(g.FindWalkable(obj.x, obj.y, boardlayerZ, int(obj.vars["leftMovement"])))
			g.CreateWalkables(g.FindWalkable(obj.x, obj.y, boardlayerZ, int(obj.vars["leftMovement"])), underLayerZ)
		} else {
			g.DeselectPawn()
			g.ClearMatrixLayer(underLayerZ)
		}
	}
}

func (g *Game) HandleAttackSelection(sx, sy int) {
	sx, sy = g.CursorXYtoMatrixGrid(attacksLayerZ, sx, sy)
	o := g.matrixLayers[attacksLayerZ].FirstObjectAt(sx, sy)
	if g.selectedAttack != o {
		if g.selectedAttack != nil {
			g.DeselectAttack()
		}
		g.SelectAttack(o)
	} else {
		g.DeselectAttack()
	}
}

func (g *Game) HandleClickControls() {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
		return
	}
	cx, cy := ebiten.CursorPosition()
	sx, sy := snapXYtoGrid(generalGridSize, cx, cy)
	if g.IsXYWithingMatrixLayerBounds(boardlayerZ, sx, sy) {
		g.HandleBoardSelection(sx, sy)
		return
	}

	if g.selectedPawn != nil &&
		g.IsXYWithingMatrixLayerBounds(attacksLayerZ, sx, cy) {
		g.HandleAttackSelection(sx, sy)
	}
}

func (g *Game) CheckForTurnEndButton() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.playerTurn = false
	}
}
