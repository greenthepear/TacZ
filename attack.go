package main

import "fmt"

const shoveBumpDamage = 1

type Attack struct {
	name         string
	objectVarKey string
	imagePackKey string
	desc         string

	script func(*Game, *GameObject)
}

func NewAttackable(game *Game, o *GameObject, vars map[string]float64) *GameObject {
	return NewGameObject("attackable", o.x, o.y, game.imagePacks["UI"], true, 0, "attackable", true, game,
		vars, nil, nil, []string{},
	)
}

func shoveScript(g *Game, o *GameObject) {
	l := g.MatrixLayerAtZ(boardlayerZ)
	x, y := o.XY()
	vecsToShove := [...]vec{
		NewVec(x, y-1), NewVec(x+1, y),
		NewVec(x, y+1), NewVec(x-1, y),
	}
	for i, v := range vecsToShove {
		if l.isWithinBounds(v.x, v.y) {
			g.AddObjectToMatrixLayer(
				NewAttackable(g, o, map[string]float64{
					"damage":   1,
					"shoveDir": 1 + float64(i), //0 - none, 1 - north, 2 - east, 3 - south, 4 - west
				}),
				underLayerZ, v.x, v.y)
		}
	}
}

func throwScript(g *Game, o *GameObject) {
	x, y := o.XY()
	vecsToThrow := [...]vec{
		NewVec(x+2, y), NewVec(x+3, y),
		NewVec(x-2, y), NewVec(x-3, y),
		NewVec(x, y+2), NewVec(x, y+3),
		NewVec(x, y-2), NewVec(x, y-3),
	}

	l := g.MatrixLayerAtZ(underLayerZ)
	for _, v := range vecsToThrow {
		if l.isWithinBounds(v.x, v.y) {
			g.AddObjectToMatrixLayer(NewGameObject("attackable", o.x, o.y, g.imagePacks["UI"], true, 0, "attackable", true, g,
				map[string]float64{
					"damage": 1,
				}, nil, nil, []string{},
			), underLayerZ, v.x, v.y)
		}
	}
}

func (g *Game) InitAttacks() {
	g.attacks = map[string]Attack{
		"shove": {"shove", "hasShove", "shove",
			"1 damage and pushes in direction",
			shoveScript},
		"throwRock": {"throwRock", "hasThrowRock", "throwRock",
			"1 damage from distance",
			throwScript},
	}
}

func (g *Game) CreateAttackObjectFromReference(a Attack) *GameObject {
	return NewGameObject(
		a.name, 0, 0, g.imagePacks["Attacks"],
		true, 0, a.imagePackKey, true, g, map[string]float64{}, nil, nil, []string{"attack"})
}

func (g *Game) CreateAttackObjectsOf(o *GameObject) []*GameObject {
	r := make([]*GameObject, 0)
	for _, a := range g.attacks {
		if o.vars[a.objectVarKey] != 0.0 {
			r = append(r, g.CreateAttackObjectFromReference(a))
		}
	}
	return r
}

func (g *Game) SelectAttack(o *GameObject, attacker *GameObject) {
	g.selectedAttack = o
	g.SimpleCreateObjectInMatrixLayer(underAttacksLayerZ, "selectedAttackIndicator", o.x, o.y, "UI", false)

	g.ClearMatrixLayer(underLayerZ)
	g.attacks[o.name].script(g, attacker)
}

func (g *Game) DeselectAttack(recreateWalkables bool) {
	g.selectedAttack = nil
	g.ClearMatrixLayer(underAttacksLayerZ)
	g.ClearMatrixLayer(underLayerZ)

	if recreateWalkables {
		g.CreateWalkablesOfSelectedPawn()
	}
}

func (g *Game) ClearAttackLayer() {
	g.ClearMatrixLayer(attacksLayerZ)
}

// Applies and returns if object has been destroyed
func (g *Game) ApplyDamage(dmg float64, receiver *GameObject, receiverLayer *MatrixLayer) bool {
	receiver.vars["leftHP"] -= dmg
	fmt.Printf("%s damaged for %.0f.\n", receiver.name, dmg)
	if receiver.vars["leftHP"] <= 0 {
		fmt.Printf("%s destroyed!\n", receiver.name)
		receiverLayer.deleteAllAt(receiver.x, receiver.y)
		return true
	}
	return false
}

func (g *Game) ApplyShove(dir float64, receiver *GameObject) bool {
	x, y := receiver.XY()
	shoveVecs := map[float64]vec{
		1.0: NewVec(x, y-1),
		2.0: NewVec(x+1, y),
		3.0: NewVec(x, y+1),
		4.0: NewVec(x-1, y),
	}
	shoveVec := shoveVecs[dir]
	l := g.MatrixLayerAtZ(boardlayerZ)
	if !l.isWithinBounds(shoveVec.x, shoveVec.y) {
		return false
	}

	oAtShoveDir := l.FirstObjectAt(shoveVec.x, shoveVec.y)
	if oAtShoveDir == nil {
		g.MoveMatrixObjects(boardlayerZ, x, y, shoveVec.x, shoveVec.y)
		return false
	}

	if oAtShoveDir.HasTag("damageable") {
		g.ApplyDamage(shoveBumpDamage, oAtShoveDir, l)
	}
	return g.ApplyDamage(shoveBumpDamage, receiver, l)
}

func (g *Game) ApplyPawnAttack(oAttackable *GameObject, receiver *GameObject, receiverLayer *MatrixLayer) {

	died := false
	if dir := oAttackable.vars["shoveDir"]; dir != 0.0 {
		died = g.ApplyShove(dir, receiver)
	}
	if dmg := oAttackable.vars["damage"]; dmg != 0.0 && !died {
		g.ApplyDamage(dmg, receiver, receiverLayer)
	}

	g.selectedPawn.vars["canAttack"] = 0.5
	g.DeselectAttack(true)
	g.ClearAttackLayer()
}

func NewEnemyAttackable(game *Game, o *GameObject, vars map[string]float64) *GameObject {
	return NewGameObject("enemyAttackable", o.x, o.y, game.imagePacks["UI"], true, 0, "attackable", true, game,
		vars, nil, nil, []string{},
	)
}
