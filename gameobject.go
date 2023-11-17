package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameObject struct {
	name    string
	x, y    int
	sprites *ImagePack
	sprIdx  int //Sprite index
	visible bool
	tags    []string //Can be used to determine enemies or flammable or something
}

func NewGameObject(name string, x, y int, sprites *ImagePack, sprIdx int, visible bool, tags []string) *GameObject {
	gobj := &GameObject{
		name:    name,
		x:       x,
		y:       y,
		sprites: sprites,
		sprIdx:  sprIdx,
		visible: visible,
		tags:    tags,
	}
	return gobj
}

func (o *GameObject) CurrSprite() *ebiten.Image {
	return o.sprites.imagesQ[o.sprIdx]
}

func (g *Game) SimpleCreateObjectInMatrixLayer(matrixLayerZ int, objName string, gridx, gridy int, imagePackName string) *GameObject {
	if g.matrixLayerNum < matrixLayerZ {
		log.Fatalf("No layer %d", matrixLayerZ)
	}

	objectcell := &g.matrixLayers[matrixLayerZ].mat[gridy][gridx].objects
	gobj := NewGameObject(objName, 0, 0, g.imagePacks[imagePackName], 0, true, []string{})
	*objectcell = append(*objectcell, gobj)
	return gobj
}

func (g *Game) MoveMatrixObjects(layerZ, fromX, fromY, toX, toY int) {
	l := g.matrixLayers[layerZ]
	o := l.ObjectCellAt(fromX, fromY).objects
	if len(o) == 0 {
		log.Fatal("Cell empty.")
	}
	g.matrixLayers[layerZ].mat[toY][toX].objects = o
	g.matrixLayers[layerZ].mat[fromY][fromX].objects = []*GameObject{}
}
