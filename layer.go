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
	name             string
	z                int
	squareLength     px
	width, height    int
	mat              [][]*ObjectCell
	xOffset, yOffset float64
	numOfObjects     int
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

func NewMatrixLayer(name string, z int, squareLength px, width, height int, xOffset, yOffset float64) *MatrixLayer {
	return &MatrixLayer{
		name:         name,
		z:            z,
		squareLength: squareLength,
		width:        width,
		height:       height,
		mat:          NewObjectMatrix(width, height),
		xOffset:      xOffset,
		yOffset:      yOffset,
		numOfObjects: 0,
	}
}

func (g *Game) CreateNewMatrixLayerOnTop(name string, squareLength px, width, height int, xOffset, yOffset float64) *MatrixLayer {
	ln := len(g.matrixLayers)
	g.matrixLayers = append(g.matrixLayers, NewMatrixLayer(name, ln, squareLength, width, height, xOffset, yOffset))
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

func (l *MatrixLayer) isWithinBounds(x, y int) bool {
	return x >= 0 && x < l.width && y >= 0 && y < l.height
}

func (l *MatrixLayer) findObjectWithNameAt(x, y int, name string) *GameObject {
	if !l.isOccupied(x, y) {
		return nil
	}
	for _, obj := range l.AllObjectsAt(x, y) {
		if obj.name == name {
			return obj
		}
	}
	return nil
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

func (l *MatrixLayer) deleteFirstAt(x, y int) {
	l.mat[y][x].objects = l.mat[y][x].objects[1:]
	if len(l.mat[y][x].objects) == 0 {
		l.numOfObjects--
	}
}

func (l *MatrixLayer) deleteAllAt(x, y int) {
	l.mat[y][x] = NewObjectCell(x, y)
	l.numOfObjects--
}

//lint:ignore U1000 shut up lint
func (l *MatrixLayer) hasObjectWithTagAt(x, y int, tag string) bool {
	return l.findObjectWithTagAt(x, y, tag) == nil
}

//lint:ignore U1000 its for debugging
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

func (g *Game) ClearMatrixLayer(layerZ int) {
	for y := range g.matrixLayers[layerZ].mat {
		for x := range g.matrixLayers[layerZ].mat[y] {
			g.matrixLayers[layerZ].mat[y][x] = NewObjectCell(x, y)
		}
	}
	g.matrixLayers[layerZ].numOfObjects = 0
}

type FreeObjectLayer struct {
	name    string
	z       int
	objects []*GameObject
}

func NewFreeObjectLayer(name string, z int) *FreeObjectLayer {
	return &FreeObjectLayer{
		name:    name,
		z:       z,
		objects: make([]*GameObject, 0),
	}
}

func (g *Game) CreateNewFreeLayerOnTop(name string) {
	l := NewFreeObjectLayer(name, len(g.freeLayers))
	g.freeLayers = append(g.freeLayers, l)
}

func (g *Game) ClearFreeLayer(layerZ int) {
	g.freeLayers[layerZ].objects = make([]*GameObject, 0)
}

func (g *Game) AllLayerObjects(layerZ int) []*GameObject {
	r := make([]*GameObject, 0)
	l := g.MatrixLayerAtZ(layerZ)
	for y := range g.matrixLayers[layerZ].mat {
		for x := range g.matrixLayers[layerZ].mat[y] {
			o := l.AllObjectsAt(x, y)
			if o == nil {
				continue
			}
			r = append(r, o...)
		}
	}
	return r
}

func (g *Game) AllLayerObjectsWithTag(layerZ int, tag string) []*GameObject {
	r := make([]*GameObject, 0)
	for _, o := range g.AllLayerObjects(layerZ) {
		if o.HasTag(tag) {
			r = append(r, o)
		}
	}
	return r
}
