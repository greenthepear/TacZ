package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
)

func skinnyScript(g *Game, o *GameObject) {
	players := g.FindObjectsWithTagWithinDistance(o.x, o.y, boardlayerZ, 4, "player")

	if len(players) > 0 {
		p := players[0]
		fmt.Printf("Found player at (%d, %d), from (%d,%d)\n",
			p.v.x, p.v.y, p.v.prev.x, p.v.prev.y)
		toX, toY := p.v.prev.x, p.v.prev.y
		if !(toX == o.x && toY == o.y) {
			g.MoveMatrixObjects(boardlayerZ, o.x, o.y, toX, toY)
		}
		g.attacks["punch"].script(g, o, p.v.x, p.v.y)
	}
}

func (g *Game) InitEnemyScripts() {
	g.enemyScripts = map[string]func(*Game, *GameObject){
		"Skinny": skinnyScript,
	}
}

func NewSkinny(game *Game, x, y int) *GameObject {
	enemyVars := map[string]float64{
		"leftHP":   3,
		"maxHP":    3,
		"hasPunch": 1,
	}
	return NewGameObject("Skinny", x, y, game.imagePacks["Zombie"], false, 0, "", true, game,
		enemyVars, []string{"enemy", "damageable"}, nil)
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

	g.pawns = g.AllLayerObjectsWithTag(boardlayerZ, "player")

	if len(g.pawns) == 0 {
		fmt.Printf("All pawns defeated. You lose!\n")
		os.Exit(0)
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

func (g *Game) ApplyEnemyAttack(oAttackable *GameObject) {
	receiverLayer := g.MatrixLayerAtZ(boardlayerZ)
	receiver := receiverLayer.FirstObjectAt(oAttackable.x, oAttackable.y)
	if receiver != nil && receiver.HasTag("damageable") {
		if dmg := oAttackable.vars["damage"]; dmg != 0.0 {
			g.ApplyDamage(dmg, receiver, receiverLayer)
		}
	}
}

func (g *Game) ApplyEnemyAttackables() {
	attackables := g.AllLayerObjects(underEnemyLayerZ)

	for _, a := range attackables {
		g.ApplyEnemyAttack(a)
		g.MatrixLayerAtZ(underEnemyLayerZ).deleteFirstAt(a.x, a.y)
	}
}
