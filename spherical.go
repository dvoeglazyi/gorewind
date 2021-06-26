package gorewind

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

import "math"

// geographic coordinates   / географические координаты   / space coordinates    / космические координаты
// longitude                / долгота                     / right ascension      / прямое восхождение
// latitude                 / широта                      / declination          / склонение

// SphericalCoords сферические координаты: долгота, широта и радиус.
type SphericalCoords struct {
	Longitude Angle   // from 0 to 360
	Latitude  Angle   // from -90 to 90
	Radius    float64 // in au
}

// NewSphericalCoords создаёт новые сферические координаты, где широта и долгота заданы градусами.
func NewCoordsFromDegrees(longitude, latitude float64) SphericalCoords {
	return NewSphericalCoords(longitude*Degree, latitude*Degree, 0)
}

// NewClockCoords создаёт новые сферические координаты, заданные через часы, минуты и секунды.
func NewClockCoords(longitudeHours uint, longitudeMinutes, longitudeSeconds float64, latitudeDegrees int, latitudeMinutes, latitudeSeconds float64) SphericalCoords {
	longitude := 15 * (float64(longitudeHours) + longitudeMinutes/60 + longitudeSeconds/3600)
	latitude := float64(latitudeDegrees) + latitudeMinutes/60 + latitudeSeconds/3600
	return NewCoordsFromDegrees(longitude, latitude)
}

// NewSphericalCoords создаёт новые сферические координаты, где широта и долгота задана радианами, а радиус астрономическими единицами.
func NewSphericalCoords(longitude, latitude, radius float64) SphericalCoords {
	return SphericalCoords{
		Longitude: newAngle(longitude),
		Latitude:  newAngle(latitude),
		Radius:    radius,
	}
}

// IsOverlap проверяет пересечение двух точек, заданных через сферические координаты, с погрешностью threshold (в радианах).
func (c *SphericalCoords) IsOverlap(coords SphericalCoords, threshold float64) bool {
	// delta is cos of angle between c and coords
	delta := c.Latitude.Sin*coords.Latitude.Sin +
		c.Latitude.Cos*coords.Latitude.Cos*
			(c.Longitude.Cos*coords.Longitude.Cos+c.Longitude.Sin*coords.Longitude.Sin)
	return math.Acos(delta) <= threshold
}

// getRotated смещеает сферические координаты на три угла Эйлера:
// rotation угол вращения относительно нулевого мередиана вокруг полюсов (изменяется долгота)
// precession угол смещения экватора относительно полюсов и смещённого нулевого мередиана (изменяются и широта и долгота).
// nutation угол вращения относительно нулевого мередиана вокруг полюсов в новой системе координат (изменяется долгота).
func (c *SphericalCoords) GetOffset(rotation, precession, nutation float64) SphericalCoords {
	coords := c.GetRotated(rotation)
	coords = coords.GetOriented(precession)
	return coords.GetRotated(nutation)
}

// GetRotated возвращает сферические координаты, вращённые относительно полюсов.
func (c SphericalCoords) GetRotated(offset float64) SphericalCoords {
	longitude := c.Longitude.float64 + offset
	if longitude >= math.Pi*2 {
		longitude -= math.Pi * 2
	}
	return SphericalCoords{
		Longitude: newAngle(longitude),
		Latitude:  c.Latitude,
		Radius:    c.Radius,
	}
}

// GetOriented возвращает координаты в новой системе координат, смещённой относительно нулевого меридиана.
func (c SphericalCoords) GetOriented(offset float64) SphericalCoords {
	const (
		degrees90  = math.Pi / 2
		degrees270 = math.Pi * 1.5
	)

	if c.Latitude.float64 == degrees90 { // is north pole
		return NewSphericalCoords(degrees90, degrees90-offset, c.Radius)
	} else if c.Latitude.float64 == -degrees90 { // is south pole
		return NewSphericalCoords(degrees270, -degrees90+offset, c.Radius)
	}

	latitudeOffsetCos := math.Cos(offset)
	latitudeOffsetSin := math.Sin(offset)
	newLatitude := math.Asin(c.Latitude.Sin*latitudeOffsetCos - c.Latitude.Cos*latitudeOffsetSin*c.Longitude.Sin)

	if c.Longitude.float64 == degrees90 || c.Longitude.float64 == degrees270 {
		return SphericalCoords{
			Longitude: c.Longitude,
			Latitude:  newAngle(newLatitude),
		}
	}

	x := c.Latitude.Cos * c.Longitude.Cos
	y := c.Latitude.Sin*latitudeOffsetSin + c.Latitude.Cos*latitudeOffsetCos*c.Longitude.Sin

	newLongitude := math.Atan(y / x)
	if x < 0 {
		newLongitude += math.Pi
	} else if y < 0 {
		newLongitude += 2 * math.Pi
	}

	return NewSphericalCoords(newLongitude, newLatitude, c.Radius)
}

// GetEcliptic возвращает эклиптические координаты.
func (c SphericalCoords) GetEcliptic() SphericalCoords {
	// TODO correction by time
	const eclipticOffset = (23 + 26/60 + 21.406/3600) * Degree
	return c.GetOriented(eclipticOffset)
}

const (
	LY     = 63241  // AU in LY
	Parsec = 206265 // AU in parsec
)

// SetRadius возвращает координаты с заданным радиусом в астрономических единицах.
func (c SphericalCoords) SetRadius(r float64) SphericalCoords {
	c.Radius = r
	return c
}
