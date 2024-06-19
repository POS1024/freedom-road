package utils

import "math"

func moveTowards(X, Y int, Angle, Step float64) (newX int, newY int) {
	rad := Angle * math.Pi / 180
	newX = X + int(Step*math.Cos(rad))
	newY = Y + int(Step*math.Sin(rad))
	return
}

func pointInRect(aX, aY, rX, rY, rW, rH int) bool {
	return aX >= rX && aX <= rX+rW &&
		aY >= rY && aY <= rY+rH
}
