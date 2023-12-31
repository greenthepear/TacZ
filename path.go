// Pathfinding
package main

import "log"

type vec struct {
	x, y int
	dist int //for pathfinding
}

func getNeighboringCells(v vec) []vec {
	possibleMoves := make([]vec, 0, 4)
	if v.y != 0 {
		possibleMoves = append(possibleMoves, vec{v.x, v.y - 1, v.dist})
	}
	if v.y != generalGridHeight-1 {
		possibleMoves = append(possibleMoves, vec{v.x, v.y + 1, v.dist})
	}
	if v.x != 0 {
		possibleMoves = append(possibleMoves, vec{v.x - 1, v.y, v.dist})
	}
	if v.x != generalGridWidth-1 {
		possibleMoves = append(possibleMoves, vec{v.x + 1, v.y, v.dist})
	}
	return possibleMoves
}

func getNeighboringCellsOfObject(gobj GameObject) []vec {
	return getNeighboringCells(vec{gobj.x, gobj.y, 0})
}

// Find points (generate slice of vectors)
// that are walkable within distance (steps) in a layer from (x,y).
// Checks only for empty object cells.
func (g *Game) findWalkable(fromX, fromY, layerZ, distance int) []vec {
	vecs := make([]vec, 0)

	q := Queue{}
	q.push(vec{fromX, fromY, 0})

	wasVisited := make([][]bool, generalGridHeight)
	for i := range wasVisited {
		wasVisited[i] = make([]bool, generalGridWidth)
	}

	distances := make(map[vec]int)
	for !q.isEmpty() {
		cell, err := q.pop()
		if err != nil {
			log.Fatal(err)
		}

		if (cell.x == fromX && cell.y == fromY) ||
			!g.MatrixLayerAtZ(layerZ).isOccupied(cell.x, cell.y) {
			vecs = append(vecs, cell)
		}

		neighbors := getNeighboringCells(cell)
		if distances[cell] > distance-1 {
			continue
		}
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
