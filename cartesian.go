package gorewind

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

import (
	"math"
)

// CartesianCoords прямоугольные координаты.
type CartesianCoords struct {
	X, Y, Z float64
}

// GetSpherical переводит прямоугольные координаты в сферические.
func (c *CartesianCoords) GetSpherical() SphericalCoords {
	longitude := math.Atan(c.Y / c.X)
	if c.X < 0 {
		longitude += math.Pi
	} else if c.Y < 0 {
		longitude += 2 * math.Pi
	}

	diagonal := math.Sqrt(c.X*c.X + c.Y*c.Y)
	latitude := math.Atan(c.Z / diagonal)
	radius := math.Sqrt(diagonal*diagonal + c.Z*c.Z)

	return NewSphericalCoords(longitude, latitude, radius)
}
