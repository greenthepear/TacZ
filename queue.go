package main

import "fmt"

type Queue struct {
	elements []vec
}

func (q *Queue) getLen() int {
	return len(q.elements)
}

func (q *Queue) isEmpty() bool {
	return q.getLen() == 0
}

func (q *Queue) push(v vec) {
	q.elements = append(q.elements, v)
}

func (q *Queue) pop() (vec, error) {
	if q.isEmpty() {
		return vec{0, 0, 0}, fmt.Errorf("popping empty queue")
	}
	element := q.elements[0]
	if q.getLen() == 1 {
		q.elements = nil
		return element, nil
	}

	q.elements = q.elements[1:]
	return element, nil
}
