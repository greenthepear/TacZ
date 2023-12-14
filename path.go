// Pathfinding
package main

import "log"

type vec struct {
	x, y int
}

func getNeighboringCells(x, y int) []vec {
	possibleMoves := make([]vec, 0, 4)
	if y != 0 {
		possibleMoves = append(possibleMoves, vec{x, y - 1})
	}
	if y != generalGridHeight-1 {
		possibleMoves = append(possibleMoves, vec{x, y + 1})
	}
	if x != 0 {
		possibleMoves = append(possibleMoves, vec{x - 1, y})
	}
	if x != generalGridWidth-1 {
		possibleMoves = append(possibleMoves, vec{x + 1, y})
	}
	return possibleMoves
}

func getNeighboringCellsOfObject(gobj GameObject) []vec {
	return getNeighboringCells(gobj.x, gobj.y)
}

// Find points (generate slice of vectors)
// that are walkable within distance (steps) in a layer from (x,y).
// Checks only for empty object cells.
func (g *Game) findWalkable(fromX, fromY, layerZ, distance int) []vec {
	vecs := make([]vec, 0)

	q := Queue{}
	q.push(vec{fromX, fromY})

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

		neighbors := getNeighboringCells(cell.x, cell.y)
		for _, nc := range neighbors {
			if !wasVisited[nc.y][nc.x] { //&& !g.MatrixLayerAtZ(layerZ).isOccupied(cell.x, cell.y) {
				wasVisited[nc.y][nc.x] = true
				distances[nc] = distances[cell] + 1
				q.push(nc)
			}
		}
	}
	return vecs
}
