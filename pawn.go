package main

import "fmt"

func NewPawn(game *Game, x, y int) *GameObject {
	pawnVariables := map[string]float64{
		"isSelected":   0,
		"maxMovement":  5,
		"leftMovement": 5,
		"leftHP":       3,
		"maxHP":        3,
		"hasShove":     1,
		"hasThrowRock": 1,
		"canAttack":    1,
	}
	return NewGameObject("Pawn", x, y, game.imagePacks["Pawn"], false, 0, "", true, game,
		pawnVariables,
		nil, nil, []string{"selectable", "damageable"})
}

func (g *Game) AddPawnToLayer(z, x, y int) {
	obj := NewPawn(g, x, y)
	g.AddObjectToMatrixLayer(obj, z, x, y)
	g.pawns = append(g.pawns, obj)
}

func (g *Game) SelectPawn(pawnObj *GameObject) {
	pawnObj.vars["isSelected"] = 1
	pawnObj.sprIdx = 1
	g.selectedPawn = pawnObj
}

func (g *Game) DeselectPawn() {
	if g.selectedPawn == nil {
		return
	}
	g.selectedPawn.vars["isSelected"] = 0
	g.selectedPawn.sprIdx = 0
	g.selectedPawn = nil
	g.DeselectAttack(false)
}

func (g *Game) CreateWalkables(vecs []vec, layerZ int) {
	for _, v := range vecs {
		obj := g.SimpleCreateObjectInMatrixLayer(underLayerZ, "walkable", v.x, v.y, "UI", false)
		obj.vars["dist"] = float64(v.dist)
		obj.sprIdx = 1
	}
}

func (g *Game) CreateWalkablesOfSelectedPawn() {
	g.CreateWalkables(
		g.FindWalkable(g.selectedPawn.x, g.selectedPawn.y, boardlayerZ, int(g.selectedPawn.vars["leftMovement"])), underLayerZ)
}

func (g *Game) InitPlayerTurn() {
	fmt.Printf("Doing player turn...\n")
	for _, pawn := range g.pawns {
		pawn.vars["leftMovement"] = pawn.vars["maxMovement"]
		pawn.vars["canAttack"] = 1.0
	}
	g.playerTurn = true
}
