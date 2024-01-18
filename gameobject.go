package main

import (
	"fmt"
	"log"
)

type GameObject struct {
	name  string
	x, y  int
	cellZ int

	sprites    *ImagePack
	sprMapMode bool
	sprIdx     int    //Sprite index
	sprKey     string //Sprite key, if in mapMode
	visible    bool

	gameRef *Game
	tags    []string //Can be used to determine enemies or flammable or something
	vars    map[string]float64

	children []*GameObject
}

func NewGameObject(
	name string, x, y int, sprites *ImagePack, mapMode bool, sprIdx int, sprMap string, visible bool, gameRef *Game,
	vars map[string]float64, tags []string, children []*GameObject) *GameObject {

	if vars == nil {
		vars = make(map[string]float64)
	}
	if children == nil {
		children = make([]*GameObject, 0)
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
		vars:       vars,
		children:   children,
	}
	return gobj
}

func (o *GameObject) HasTag(tag string) bool {
	for _, t := range o.tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (o *GameObject) HasChildren() bool {
	return len(o.children) > 0
}

func (o *GameObject) XY() (int, int) {
	return o.x, o.y
}

func (o *GameObject) Vec() vec {
	return NewVec(o.x, o.y)
}

func (o *GameObject) ScreenPosition(l MatrixLayer) (int, int) {
	return o.x*int(l.squareLength) + int(l.xOffset),
		o.y*int(l.squareLength) + int(l.yOffset)

}

func (o *GameObject) IsMarkedForDeletion() bool {
	return o.vars["DELETED"] == 1.0
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

	gobj := NewGameObject(objName, gridx, gridy, imgPack, sprMapMode, 0, sprKey, true, g, nil, []string{}, nil)
	g.AddObjectToMatrixLayer(gobj, matrixLayerZ, gridx, gridy)
	return gobj
}

func (g *Game) AddObjectToMatrixLayer(gobj *GameObject, matrixLayerZ, gridx, gridy int) {
	if len(g.matrixLayers) < matrixLayerZ {
		log.Fatalf("No layer %d", matrixLayerZ)
	}
	gobj.x, gobj.y = gridx, gridy
	objectcell := &g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects
	gobj.cellZ = len(*objectcell)
	g.matrixLayers[matrixLayerZ].numOfObjects++
	*objectcell = append(*objectcell, gobj)
}

func (g *Game) MoveMatrixObject(layerZ, fromX, fromY, toX, toY, cellZ int) {
	l := g.matrixLayers[layerZ]
	o := l.ObjectAtZ(fromX, fromY, cellZ)
	o.x = toX
	o.y = toY
	g.AddObjectToMatrixLayer(o, layerZ, toX, toY)
	l.deleteAtZ(fromX, fromY, cellZ, false)
}

func (g *Game) MoveMatrixObjects(layerZ, fromX, fromY, toX, toY int) error {
	l := g.matrixLayers[layerZ]
	cell := l.ObjectCellAt(fromX, fromY)
	if len(cell.objects) == 0 {
		return fmt.Errorf("[%d] Cell at (%d, %d) empty:\n%#v", layerZ, fromX, fromY, cell)
	}
	for _, obj := range cell.objects {
		obj.x = toX
		obj.y = toY
	}
	g.matrixLayers[layerZ].mat[toY][toX].objects = cell.objects
	g.matrixLayers[layerZ].mat[fromY][fromX] = NewObjectCell(fromX, fromY)
	return nil
}

func (g *Game) AddObjectToFreeLayer(z int, o *GameObject) {
	if z > len(g.freeLayers) {
		log.Fatalf("Error while adding object:\n\n%v\n\nNo layer `%d`", o, z)
	}
	g.freeLayers[z].objects = append(g.freeLayers[z].objects, o)
}
