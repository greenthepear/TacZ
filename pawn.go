package main

func NewPawn(game *Game, x, y int) *GameObject {
	return NewGameObject("Pawn", x, y, game.imagePacks["Pawn"], 0, true, game,
		map[string]float64{"isSelected": 0}, nil, nil, []string{"selectable"})
}

func (g *Game) selectPawn(pawnObj *GameObject) {
	pawnObj.vars["isSelected"] = 1
	pawnObj.sprIdx = 1
	g.selectedPawn = pawnObj
}

func (g *Game) deselectPawn() {
	g.selectedPawn.vars["isSelected"] = 0
	g.selectedPawn.sprIdx = 0
	g.selectedPawn = nil
}
