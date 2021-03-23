package models

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (p *Point) IsEqualToPoint(point Point) bool {
	return p.X == point.X && p.Y == point.Y
}

func (p *Point) IsEqualToCords(x int, y int) bool {
	return p.X == x && p.Y == y
}
