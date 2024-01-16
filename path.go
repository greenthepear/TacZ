// Pathfinding
package main

import (
	"log"
)

type vec struct {
	x, y int

	//for pathfinding
	dist int
	prev *vec
}

func NewVec(x, y int) vec {
	return vec{x, y, 0, nil}
}

func NewVecDist(x, y, d int) vec {
	return vec{x, y, d, nil}
}

func (v vec) NewChildVec(x, y int) vec {
	return vec{x, y, v.dist, &v}
}

func getNeighboringCells(v vec) []vec {
	possibleMoves := make([]vec, 0, 4)
	if v.y != 0 {
		possibleMoves = append(possibleMoves, v.NewChildVec(v.x, v.y-1))
	}
	if v.y != generalGridHeight-1 {
		possibleMoves = append(possibleMoves, v.NewChildVec(v.x, v.y+1))
	}
	if v.x != 0 {
		possibleMoves = append(possibleMoves, v.NewChildVec(v.x-1, v.y))
	}
	if v.x != generalGridWidth-1 {
		possibleMoves = append(possibleMoves, v.NewChildVec(v.x+1, v.y))
	}
	return possibleMoves
}

func getNeighboringCellsOfObject(gobj GameObject) []vec {
	return getNeighboringCells(vec{gobj.x, gobj.y, 0, nil})
}

func initBFSdata(fromX, fromY, width, height int) (Queue, [][]bool, map[vec]int) {
	q := Queue{}
	q.push(vec{fromX, fromY, 0, nil})

	wasVisited := make([][]bool, height)
	for i := range wasVisited {
		wasVisited[i] = make([]bool, width)
	}

	distances := make(map[vec]int)

	return q, wasVisited, distances
}

// Find points (generate slice of vectors)
// that are walkable within distance (steps) in a layer from (x,y).
// Checks only for empty object cells, unless they're occupied by
// object with tag collIgnoreTag.
func (g *Game) FindWalkable(fromX, fromY, layerZ, distance int) []vec {
	vecs := make([]vec, 0)

	q, wasVisited, distances := initBFSdata(
		fromX, fromY, g.MatrixLayerAtZ(layerZ).width, g.MatrixLayerAtZ(layerZ).height)

	for !q.isEmpty() {
		cell, err := q.pop()
		if err != nil {
			log.Fatal(err)
		}

		if (cell.x == fromX && cell.y == fromY) ||
			!g.MatrixLayerAtZ(layerZ).isOccupied(cell.x, cell.y) {
			vecs = append(vecs, cell)
		}

		if distances[cell] > distance-1 {
			continue
		}
		neighbors := getNeighboringCells(cell)
		for _, nc := range neighbors {
			if !wasVisited[nc.y][nc.x] && !g.MatrixLayerAtZ(layerZ).isOccupied(nc.x, nc.y) {
				wasVisited[nc.y][nc.x] = true
				nc.dist = distances[cell] + 1
				distances[nc] = nc.dist
				q.push(nc)
			}
		}
	}
	return vecs
}

type objectWithVec struct {
	v vec
	o *GameObject
}

func (g *Game) FindObjectsWithTagWithinDistance(fromX, fromY, layerZ, distance int, tag string) []objectWithVec {
	objects := make([]objectWithVec, 0)

	q, wasVisited, distances := initBFSdata(
		fromX, fromY, g.MatrixLayerAtZ(layerZ).width, g.MatrixLayerAtZ(layerZ).height)

	for !q.isEmpty() {
		cell, err := q.pop()

		//fmt.Printf("Checking %v\n", cell)
		if err != nil {
			log.Fatal(err)
		}

		if o := g.MatrixLayerAtZ(layerZ).FirstObjectAt(cell.x, cell.y); o != nil && o.HasTag(tag) {
			objects = append(objects, objectWithVec{cell, o})
		}

		if distances[cell] > distance-1 {
			continue
		}
		neighbors := getNeighboringCells(cell)
		for _, nc := range neighbors {
			if o := g.MatrixLayerAtZ(layerZ).FirstObjectAt(nc.x, nc.y); o != nil && o.HasTag(tag) {
				objects = append(objects, objectWithVec{nc, o})
				continue
			}

			if !wasVisited[nc.y][nc.x] && !g.MatrixLayerAtZ(layerZ).isOccupied(nc.x, nc.y) {
				wasVisited[nc.y][nc.x] = true
				nc.dist = distances[cell] + 1
				distances[nc] = nc.dist
				q.push(nc)
			}
		}
	}
	return objects
}

func (g *Game) ObjectsWithinWalkable(fromX, fromY, layerZ, distance int, tag string) []objectWithVec {
	r := make([]objectWithVec, 0)
	vecs := g.FindWalkable(fromX, fromY, boardlayerZ, distance)
	l := g.MatrixLayerAtZ(layerZ)
	for _, v := range vecs {
		if o := l.FirstObjectAt(v.x, v.y); o != nil {
			r = append(r, objectWithVec{v, o})
		}
	}
	return r
}

func (g *Game) FindObjectsWithTagWithinWalkable(fromX, fromY, layerZ, distance int,
	tag string, oIgnore *GameObject) []objectWithVec {

	objects := g.ObjectsWithinWalkable(fromX, fromY, layerZ, distance, tag)
	r := make([]objectWithVec, 0)
	for _, o := range objects {
		if oIgnore != nil && o.o == oIgnore {
			continue
		}
		if o.o.HasTag(tag) {
			r = append(r, o)
		}
	}

	return r
}
