package main

type Attack struct {
	name         string
	objectVarKey string
	imagePackKey string
}

func NewAttacks() []Attack {
	return []Attack{
		{"shove", "hasShove", "shove"},
		{"throwRock", "hasThrowRock", "throwRock"},
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