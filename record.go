package gorewind

import (
	"strconv"
	"strings"
)

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

// AstronomicalObject небесный объект.
type AstronomicalObject struct {
	Catalogue      string
	Index          uint
	Designation    Designation
	Name           string
	LocalName      string
	AlternateNames []string
	Magnitude      float64
	Coords         SphericalCoords
}

func (ao *AstronomicalObject) GetCoords() SphericalCoords {
	return ao.Coords
}

func (ao *AstronomicalObject) GetRecord() []string {
	s := []string{
		ao.Name,
		ao.LocalName,
		"",
		ao.Designation.Constellation,
		"",
		"",
		"",
		"",
		"",
		"",
		strings.Join(ao.AlternateNames, ","),
	}
	if code := ao.Designation.GetCode(); code != "" {
		s[2] = code
	}
	if ao.Designation.InSystemIndex != 0 {
		s[4] = strconv.FormatUint(uint64(ao.Designation.InSystemIndex), 10)
	}
	if ao.Catalogue != "" && ao.Index != 0 {
		s[5] = ao.Catalogue + " " + strconv.FormatUint(uint64(ao.Index), 10)
	}
	if ao.Magnitude != 0 {
		s[6] = strconv.FormatFloat(ao.Magnitude, 'f', 2, 64)
	}
	if lat, long := ao.Coords.Latitude.Degrees(), ao.Coords.Longitude.Degrees(); lat != 0 || long != 0 {
		s[7] = strconv.FormatFloat(long, 'f', 6, 64)
		s[8] = strconv.FormatFloat(lat, 'f', 6, 64)
	}
	if ao.Coords.Radius != 0 {
		s[9] = strconv.FormatUint(uint64(ao.Coords.Radius), 10)
	}
	return s
}

type Location struct {
	Name        string
	LocalName   string
	Population  uint
	CountryCode string
	Coords      SphericalCoords
}

func (l *Location) GetCoords() SphericalCoords {
	return l.Coords
}
