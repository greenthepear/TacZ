package main

type Attack struct {
	name         string
	objectVarKey string
	imagePackKey string
	desc         string

	script func(*Game, *GameObject)
}

func shoveScript(g *Game, o *GameObject) {
	attackableObj := *NewGameObject("attackable", o.x, o.y, g.imagePacks["UI"], true, 0, "attackable", true, g,
		map[string]float64{
			"damage":   1,
			"shoveDir": 0, //0 - no shove, 1 - north, 2 - east, 3 - south, 4 - west
		}, nil, nil, []string{},
	)
	l := g.MatrixLayerAtZ(boardlayerZ)
	if l.isWithinBounds(o.x, o.y-1) {
		attackableObj.vars["shoveDir"] = 1
		g.AddObjectToMatrixLayer(&attackableObj, underLayerZ, o.x, o.y-1)
	}
	if l.isWithinBounds(o.x+1, o.y) {
		attackableObj.vars["shoveDir"] = 2
		g.AddObjectToMatrixLayer(&attackableObj, underLayerZ, o.x+1, o.y)
	}
	if l.isWithinBounds(o.x, o.y+1) {
		attackableObj.vars["shoveDir"] = 3
		g.AddObjectToMatrixLayer(&attackableObj, underLayerZ, o.x, o.y+1)
	}
	if l.isWithinBounds(o.x-1, o.y) {
		attackableObj.vars["shoveDir"] = 4
		g.AddObjectToMatrixLayer(&attackableObj, underLayerZ, o.x-1, o.y)
	}
}

func throwScript(g *Game, o *GameObject) {
	x, y := o.x, o.y
	vecsToThrow := [...]vec{
		NewVec(x+2, y), NewVec(x+3, y),
		NewVec(x-2, y), NewVec(x-3, y),
		NewVec(x, y+2), NewVec(x, y+3),
		NewVec(x, y-2), NewVec(x, y-3),
	}

	attackableObj := *NewGameObject("attackable", o.x, o.y, g.imagePacks["UI"], true, 0, "attackable", true, g,
		map[string]float64{
			"damage": 1,
		}, nil, nil, []string{},
	)

	l := g.MatrixLayerAtZ(underLayerZ)
	for _, v := range vecsToThrow {
		if l.isWithinBounds(v.x, v.y) {
			g.AddObjectToMatrixLayer(&attackableObj, underLayerZ, v.x, v.y)
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
