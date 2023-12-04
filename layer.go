package main

import (
	"fmt"
)

type px int

type ObjectCell struct {
	x, y    int
	objects []*GameObject
}

type MatrixLayer struct {
	name          string
	z             int
	squareLength  px
	width, height int
	mat           [][]*ObjectCell
}

func NewObjectCell(x, y int) *ObjectCell {
	oc := &ObjectCell{
		x:       x,
		y:       y,
		objects: make([]*GameObject, 0),
	}
	return oc
}

func NewObjectMatrix(width, height int) [][]*ObjectCell {
	m := make([][]*ObjectCell, height)
	for y := 0; y < height; y++ {
		m[y] = make([]*ObjectCell, width)

		for x := 0; x < width; x++ {
			m[y][x] = NewObjectCell(x, y)
		}
	}
	return m
}

func NewMatrixLayer(name string, z int, squareLength px, width, height int) *MatrixLayer {
	return &MatrixLayer{
		name:         name,
		z:            z,
		squareLength: squareLength,
		width:        width,
		height:       height,
		mat:          NewObjectMatrix(width, height),
	}
}

/*
func (g *Game) CreateNewMatrixLayerAtZ(z int, name string, squareLength px, width, height int) {
	if z < 0 || z > maxNumberOfLayers {
		log.Fatal("z not withing layer range")
	}
	if g.matrixLayers[z] != nil {
		log.Fatalf("z already occupied by layer \"%s\"", g.matrixLayers[z].name)
	}
	g.matrixLayers[z] = NewMatrixLayer(name, z, squareLength, width, height)
}
*/

func (g *Game) CreateNewMatrixLayerOnTop(name string, squareLength px, width, height int) *MatrixLayer {
	ln := g.matrixLayerNum
	g.matrixLayers[ln] = NewMatrixLayer(name, ln, squareLength, width, height)
	g.matrixLayerNum++
	return g.matrixLayers[ln]
}

func (g *Game) MatrixLayerAtZ(z int) *MatrixLayer {
	return g.matrixLayers[z]
}

func (l *MatrixLayer) ObjectCellAt(x, y int) *ObjectCell {
	return l.mat[y][x]
}

func (l *MatrixLayer) AllObjectsAt(x, y int) []*GameObject {
	return l.mat[y][x].objects
}

func (l *MatrixLayer) ObjectAtZ(x, y, z int) *GameObject {
	le := len(l.AllObjectsAt(x, y))
	if le <= z {
		return nil
	}
	return l.AllObjectsAt(x, y)[z]
}

func (l *MatrixLayer) FirstObjectAt(x, y int) *GameObject {
	return l.ObjectAtZ(x, y, 0)
}

func (l *MatrixLayer) isOccupied(x, y int) bool {
	return l.FirstObjectAt(x, y) != nil
}

func (l *MatrixLayer) findObjectWithTagAt(x, y int, tag string) *GameObject {
	if !l.isOccupied(x, y) {
		return nil
	}
	for _, obj := range l.AllObjectsAt(x, y) {
		for _, t := range obj.tags {
			if t == tag {
				return obj
			}
		}
	}
	return nil
}

func (l *MatrixLayer) hasObjectWithTagAt(x, y int, tag string) bool {
	return l.findObjectWithTagAt(x, y, tag) == nil
}

func (l MatrixLayer) printMatrix() {
	fmt.Printf("Layer '%s' (%d) %d x %d \n", l.name, l.z, l.width, l.height)
	for y := 0; y < l.height; y++ {
		for x := 0; x < l.width; x++ {
			if l.isOccupied(x, y) {
				fmt.Print(l.FirstObjectAt(x, y).name)
			} else {
				fmt.Print(" nil ")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}
