package main

import (
	"fmt"
	"math/rand"
)

func NewEnemy(game *Game, x, y int) *GameObject {
	return NewGameObject("Zombie", x, y, game.imagePacks["Zombie"], 0, true, game,
		map[string]float64{}, nil, nil, []string{"enemy"})
}

func (g *Game) AddEnemyToLayer(z, x, y int) {
	obj := NewEnemy(g, x, y)
	g.AddObjectToMatrixLayer(obj, z, x, y)
	g.enemies = append(g.enemies, obj)
}

func (g *Game) DoEnemyTurn() {
	fmt.Printf("Doing enemy turn...\n")
	playLayerZ := 1
	for _, enemy := range g.enemies {
		//x, y :=
		possibleMoves := make([]vec, 0, 4)
		neighboringCells := getNeighboringCellsOfObject(*enemy)
		for _, v := range neighboringCells {
			if !g.matrixLayers[playLayerZ].isOccupied(v.x, v.y) {
				possibleMoves = append(possibleMoves, v)
			}
		}

		if len(possibleMoves) == 0 {
			continue
		}

		chosenDirVec := possibleMoves[rand.Intn(len(possibleMoves))]

		g.MoveMatrixObjects(playLayerZ, enemy.x, enemy.y, chosenDirVec.x, chosenDirVec.y)
	}

	fmt.Printf("Doing player turn...\n")
	g.playerTurn = true
}
