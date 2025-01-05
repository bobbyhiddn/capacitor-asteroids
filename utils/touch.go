package utils

// IsPointInCircle checks if a point (x, y) is within a circle at (circleX, circleY) with given radius
func IsPointInCircle(x, y, circleX, circleY, radius float64) bool {
	dx := x - circleX
	dy := y - circleY
	distanceSquared := dx*dx + dy*dy
	return distanceSquared <= radius*radius
}
