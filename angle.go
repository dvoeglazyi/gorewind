package gorewind

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

import "math"

const (
	Degree = math.Pi / 180 // from radians
	Radian = 180 / math.Pi // from degrees
)

// Angle угол с заранее рассчитанными синусом и косинусом.
type Angle struct {
	float64 // radians
	Sin     float64
	Cos     float64
}

// Radians возвращает угол в радианах.
func (a Angle) Radians() float64 {
	return a.float64
}

// Degrees возвращает угол в градусах.
func (a Angle) Degrees() float64 {
	return a.float64 * Radian
}

// NewAngleFromDegrees создаёт новый угол, заданный градусами.
func NewAngleFromDegrees(degree int, minutes, seconds float64) Angle {
	return newAngle(float64(degree) + minutes/60 + seconds/3600)
}

// NewClockAngle создаёт новый угол, заданный часами, минутами и секундами.
func NewClockAngle(hours uint, minutes, seconds float64) Angle {
	return newAngle(15 * (float64(hours) + minutes/60 + seconds/3600))
}

func newAngle(angle float64) Angle {
	return Angle{
		float64: angle,
		Sin:     math.Sin(angle),
		Cos:     math.Cos(angle),
	}
}
