package main

import (
	"fmt"
	"log"
	"math/rand"
)

func skinnyScript(g *Game, o *GameObject) {
	players := g.FindObjectsWithTagWithinWalkable(o.x, o.y, boardlayerZ, 4, "player", o)
	//dist := len(players) > 0
	fmt.Println(players)
}

func (g *Game) InitEnemyScripts() {
	g.enemyScripts = map[string]func(*Game, *GameObject){
		"Skinny": skinnyScript,
	}
}

func NewSkinny(game *Game, x, y int) *GameObject {
	enemyVars := map[string]float64{
		"leftHP": 3,
		"maxHP":  3,
	}
	return NewGameObject("Skinny", x, y, game.imagePacks["Zombie"], false, 0, "", true, game,
		enemyVars, nil, nil, []string{"enemy", "damageable"})
}

func (g *Game) AddSkinnyToLayer(z, x, y int) {
	obj := NewSkinny(g, x, y)
	g.AddObjectToMatrixLayer(obj, z, x, y)
	g.enemies = append(g.enemies, obj)
}

func (g *Game) DoEnemyTurn() {
	fmt.Printf("Doing enemy turn...\n")
	for _, enemy := range g.enemies {
		if g.enemyScripts[enemy.name] == nil {
			log.Printf("Enemy '%v' has no script.\n", enemy.name)
			g.MoveObjectInRandomDirection(enemy)
		} else {
			g.enemyScripts[enemy.name](g, enemy)
		}
	}

	g.InitPlayerTurn()
}

func (g *Game) MoveObjectInRandomDirection(o *GameObject) {
	possibleMoves := make([]vec, 0, 4)
	neighboringCells := getNeighboringCellsOfObject(*o)
	for _, v := range neighboringCells {
		if !g.matrixLayers[boardlayerZ].isOccupied(v.x, v.y) {
			possibleMoves = append(possibleMoves, v)
		}
	}

	if len(possibleMoves) == 0 {
		return
	}

	chosenDirVec := possibleMoves[rand.Intn(len(possibleMoves))]

	g.MoveMatrixObjects(boardlayerZ, o.x, o.y, chosenDirVec.x, chosenDirVec.y)
}
