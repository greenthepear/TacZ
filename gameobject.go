package main

import (
	"log"
)

type GameObject struct {
	name       string
	x, y       int
	sprites    *ImagePack
	sprIdx     int //Sprite index
	visible    bool
	gameRef    *Game
	tags       []string //Can be used to determine enemies or flammable or something
	vars       map[string]float64
	updateFunc func()
	createFunc func()
}

func NewGameObject(
	name string, x, y int, sprites *ImagePack, sprIdx int, visible bool, gameRef *Game,
	vars map[string]float64, updateFunc func(), createFunc func(), tags []string) *GameObject {

	varsmap := vars
	if vars == nil {
		varsmap = make(map[string]float64)
	}

	gobj := &GameObject{
		name:       name,
		x:          x,
		y:          y,
		sprites:    sprites,
		sprIdx:     sprIdx,
		visible:    visible,
		gameRef:    gameRef,
		tags:       tags,
		vars:       varsmap,
		updateFunc: updateFunc,
		createFunc: createFunc,
	}
	return gobj
}

func (g *Game) SimpleCreateObjectInMatrixLayer(matrixLayerZ int, objName string, gridx, gridy int, imagePackName string) *GameObject {
	if g.matrixLayerNum < matrixLayerZ {
		log.Fatalf("No layer %d", matrixLayerZ)
	}

	objectcell := &g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects
	gobj := NewGameObject(objName, 0, 0, g.imagePacks[imagePackName], 0, true, g, nil, nil, nil, []string{})
	*objectcell = append(*objectcell, gobj)
	return gobj
}

func (g *Game) AddObjectToMatrixLayer(gobj *GameObject, matrixLayerZ, gridx, gridy int) {
	if g.matrixLayerNum < matrixLayerZ {
		log.Fatalf("No layer %d", matrixLayerZ)
	}
	objectcell := &g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects
	for _, o := range *objectcell {
		o.x = gridx
		o.y = gridy
	}
	*objectcell = append(*objectcell, gobj)
}

func (g *Game) MoveMatrixObjects(layerZ, fromX, fromY, toX, toY int) {
	l := g.matrixLayers[layerZ]
	o := l.ObjectCellAt(fromX, fromY).objects
	if len(o) == 0 {
		log.Fatal("Cell empty.")
	}
	for _, obj := range o {
		obj.x = toX
		obj.y = toY
	}
	g.matrixLayers[layerZ].mat[toY][toX].objects = o
	g.matrixLayers[layerZ].mat[fromY][fromX].objects = []*GameObject{}
}
