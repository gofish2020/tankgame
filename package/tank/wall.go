package tank

import "math"

type Line struct {
	X1, Y1, X2, Y2 float64
}

func (l *Line) angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

type Object struct {
	Walls []Line
}

func (o Object) points() [][2]float64 {
	// Get one of the endpoints for all segments,
	// + the startpoint of the first one, for non-closed paths
	var points [][2]float64
	for _, wall := range o.Walls {
		points = append(points, [2]float64{wall.X2, wall.Y2})
	}
	p := [2]float64{o.Walls[0].X1, o.Walls[0].Y1}
	if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
		points = append(points, [2]float64{o.Walls[0].X1, o.Walls[0].Y1})
	}
	return points
}

// 矩形的四个边
func Rect(x, y, w, h float64) []Line {
	return []Line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}



