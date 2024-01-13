package main

import (
	"log"
)

type GameObject struct {
	name string
	x, y int

	sprites    *ImagePack
	sprMapMode bool
	sprIdx     int    //Sprite index
	sprKey     string //Sprite key, if in mapMode
	visible    bool

	gameRef    *Game
	tags       []string //Can be used to determine enemies or flammable or something
	vars       map[string]float64
	updateFunc func()
	createFunc func()
}

func NewGameObject(
	name string, x, y int, sprites *ImagePack, mapMode bool, sprIdx int, sprMap string, visible bool, gameRef *Game,
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
		sprMapMode: mapMode,
		sprIdx:     sprIdx,
		sprKey:     sprMap,
		visible:    visible,
		gameRef:    gameRef,
		tags:       tags,
		vars:       varsmap,
		updateFunc: updateFunc,
		createFunc: createFunc,
	}
	return gobj
}

func (g *Game) SimpleCreateObjectInMatrixLayer(matrixLayerZ int, objName string, gridx, gridy int, imagePackName string, sprMapMode bool) *GameObject {
	if len(g.matrixLayers) < matrixLayerZ {
		log.Fatalf("No layer %d", matrixLayerZ)
	}

	imgPack := g.imagePacks[imagePackName]
	imgPackImages := imgPack.images
	sprKey := ""
	if sprMapMode {
		for k := range imgPackImages { //random key
			sprKey = k
			break
		}
	}

	objectcell := &g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects
	gobj := NewGameObject(objName, gridx, gridy, imgPack, sprMapMode, 0, sprKey, true, g, nil, nil, nil, []string{})

	*objectcell = append(*objectcell, gobj)
	return gobj
}

func (g *Game) AddObjectToMatrixLayer(gobj *GameObject, matrixLayerZ, gridx, gridy int) {
	if len(g.matrixLayers) < matrixLayerZ {
		log.Fatalf("No layer %d", matrixLayerZ)
	}
	gobj.x, gobj.y = gridx, gridy
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

func (g *Game) AddObjectToFreeLayer(z int, o *GameObject) {
	if z > len(g.freeLayers) {
		log.Fatalf("Error while adding object:\n\n%v\n\nNo layer `%d`", o, z)
	}
	g.freeLayers[z].objects = append(g.freeLayers[z].objects, o)
}
