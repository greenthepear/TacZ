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

type vec struct {
	x, y int
}

func (g *Game) DoEnemyTurn() {
	fmt.Printf("Doing enemy turn...\n")
	for _, enemy := range g.enemies {
		x, y := enemy.x, enemy.y
		possibleMoves := make([]vec, 0, 4)
		if enemy.y != 0 {
			if !g.matrixLayers[1].isOccupied(x, y-1) {
				possibleMoves = append(possibleMoves, vec{x, y - 1})
			}
		}
		if enemy.y != generalGridHeight-1 {
			if !g.matrixLayers[1].isOccupied(x, y+1) {
				possibleMoves = append(possibleMoves, vec{x, y + 1})
			}
		}
		if enemy.x != 0 {
			if !g.matrixLayers[1].isOccupied(x-1, y) {
				possibleMoves = append(possibleMoves, vec{x - 1, y})
			}
		}
		if enemy.x != generalGridWidth-1 {
			if !g.matrixLayers[1].isOccupied(x+1, y) {
				possibleMoves = append(possibleMoves, vec{x + 1, y})
			}
		}

		if len(possibleMoves) == 0 {
			continue
		}

		chosenDirVec := possibleMoves[rand.Intn(len(possibleMoves))]

		g.MoveMatrixObjects(1, x, y, chosenDirVec.x, chosenDirVec.y)
	}

	fmt.Printf("Doing player turn...\n")
	g.playerTurn = true
}
