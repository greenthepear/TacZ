package main

import (
	"fmt"
	"log"
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

func (l *MatrixLayer) deleteAllAt(x, y int, addDeletedVar bool) {
	if addDeletedVar {
		for _, o := range l.mat[y][x].objects {
			o.MarkForDeletion()
		}
	}
	l.mat[y][x] = NewObjectCell(x, y)
	l.numOfObjects--
}

func (l *MatrixLayer) deleteFirstAt(x, y int) {
	if len(l.mat[y][x].objects) == 1 {
		l.deleteAllAt(x, y, true)
		return
	}
	l.mat[y][x].objects = l.mat[y][x].objects[1:]
}

func (l *MatrixLayer) deleteAtZ(x, y, z int, addDeletedVar bool) error {
	if len(l.mat[y][x].objects) <= z {
		return fmt.Errorf("DELETE ERROR deleteAtZ within layer %d '%v': no object at (%d,%d)*`%d`, all objects:\n%v", l.z, l.name, x, y, z, l.mat[y][x].objects)
	}
	if len(l.mat[y][x].objects) == 1 {
		l.deleteAllAt(x, y, addDeletedVar)
		return nil
	}

	if addDeletedVar {
		l.mat[y][x].objects[z].MarkForDeletion()
	}

	l.mat[y][x].objects = append(l.mat[y][x].objects[:z], l.mat[y][x].objects[z+1:]...)
	if len(l.mat[y][x].objects) == 0 {
		l.numOfObjects--
		return nil
	}
	for i := range l.mat[y][x].objects {
		l.mat[y][x].objects[i].cellZ = i
	}
	return nil
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
			for _, o := range g.matrixLayers[layerZ].mat[y][x].objects {
				o.MarkForDeletion()
			}
			g.matrixLayers[layerZ].mat[y][x] = NewObjectCell(x, y)
		}
	}
	g.matrixLayers[layerZ].numOfObjects = 0
}

func (l *MatrixLayer) AllObjects() []*GameObject {
	r := make([]*GameObject, 0)
	for y := range l.mat {
		for x := range l.mat[y] {
			o := l.AllObjectsAt(x, y)
			if o == nil {
				continue
			}
			r = append(r, o...)
		}
	}
	return r
}

func (l *MatrixLayer) checkForIntegrity() {
	for y := range l.mat {
		for x := range l.mat[y] {
			obs := l.mat[y][x]
			if obs.x != x || obs.y != y {
				log.Printf("Layer '%s' has mismatched cell x y: real (%d,%d) != (%d, %d):\n%#v",
					l.name, x, y, obs.x, obs.y, obs)
			}

			for z, o := range obs.objects {
				if o.x != x || o.y != y {
					log.Printf("Layer '%s' has mismatched object x y: real (%d,%d) != (%d, %d):\n%#v",
						l.name, x, y, o.x, o.y, o)
				}
				if o.IsMarkedForDeletion() {
					log.Printf("Layer '%s' has undeleted object market for deletion: %#v",
						l.name, o)
				}
				if l.ObjectAtZ(x, y, z) != o {
					if o.cellZ != z {
						log.Printf("Layer '%s' has mismatched object z: real %d != %d:\n%#v\nwithin\n%#v",
							l.name, len(l.AllObjectsAt(x, y)), z, o, obs)
					}
				}
			}
		}
	}
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
	return g.MatrixLayerAtZ(layerZ).AllObjects()
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
