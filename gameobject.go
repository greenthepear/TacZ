package main

import (
	"fmt"
	"log"
)

type GameObject struct {
	name  string
	x, y  int
	cellZ int

	layer *MatrixLayer

	sprites    *ImagePack
	sprMapMode bool
	sprIdx     int    //Sprite index
	sprKey     string //Sprite key, if in mapMode
	visible    bool

	gameRef *Game
	tags    []string //Can be used to determine enemies or flammable or something
	vars    map[string]float64

	children []*GameObject
	parent   *GameObject
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

func (o *GameObject) AddChild(co *GameObject) {
	o.children = append(o.children, co)
	co.parent = o
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

func (o *GameObject) MarkForDeletion() {
	o.vars["DELETED"] = 1.0
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

	gobj.cellZ = len(g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects)
	gobj.layer = g.matrixLayers[matrixLayerZ]
	g.matrixLayers[matrixLayerZ].numOfObjects++
	g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects = append(g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects, gobj)

	//Lazy cellZ refresh
	for i := range g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects {
		g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects[i].cellZ = i
	}
}

func (g *Game) MoveMatrixObject(layerZ, fromX, fromY, toX, toY, cellZ int) {
	l := g.matrixLayers[layerZ]
	o := l.ObjectAtZ(fromX, fromY, cellZ)
	if o == nil {
		log.Fatalf(
			"MOVE ERROR while trying to move from (%d,%d) to (%d,%d) in matrix layer %v\n:\tno object at (%d,%d)*%d\n%#v",
			fromX, fromY, toX, toY, l.name, fromX, fromY, cellZ, l.mat[fromY][fromX])
	}
	o.x = toX
	o.y = toY
	g.AddObjectToMatrixLayer(o, layerZ, toX, toY)
	err := l.deleteAtZ(fromX, fromY, cellZ, false)
	if err != nil {
		log.Fatalf("MOVE ERROR while moving object %#v:\n%v", o, err)
	}
}

func (o *GameObject) MoveTo(x, y int) {
	o.gameRef.MoveMatrixObject(o.layer.z, o.x, o.y, x, y, o.cellZ)
}

func (o *GameObject) checkForIntegrity() {
	g := o.gameRef
	l := o.layer

	objects := g.AllLayerObjects(l.z)

	foundself := false
	for _, oo := range objects {
		if oo == o {
			foundself = true
		}
	}
	if !foundself {
		log.Fatalf("Object %#v couldn't find itself", o)
	}
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
