package main

import (
	"fmt"
	"math"
)

type Point2D struct {
	x, y float64
}

func NewPoint2D(x, y float64) Point2D {
	return Point2D{x, y}
}
func (p Point2D) Distance(other Point2D) float64 {
	return math.Sqrt(math.Pow(other.x-p.x, 2) + math.Pow(other.y-p.y, 2))
}
func main() {
	p := NewPoint2D(1, 2)
	fmt.Println(p.Distance(NewPoint2D(2, 1)))
}
