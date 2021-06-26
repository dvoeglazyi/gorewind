package gorewind

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
)

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

// GeoNames Gazetteer.
// Для использования с каталогом городов, например, cities15000.zip.
// http://download.geonames.org/export/dump/

func ReadCitiesCatalogue(path string) ([]*Location, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []*Location
	csvReader := csv.NewReader(file)
	csvReader.Comma = '\t'

	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		latitude, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}
		longitude, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, err
		}
		population, err := strconv.ParseUint(record[14], 10, 64)
		if err != nil {
			return nil, err
		}
		location := Location{
			Name:        record[1],
			LocalName:   getLocalCityName(record[3]),
			Population:  uint(population),
			CountryCode: record[8],
			Coords:      NewCoordsFromDegrees(longitude, latitude),
		}
		result = append(result, &location)
	}
	return result, nil
}

func getLocalCityName(s string) string {
	fields := strings.Split(s, ",")
	for _, field := range fields {
		if isRussianText(field) {
			return field
		}
	}
	return ""
}

func isRussianText(text string) bool {
	for _, r := range []rune(text) {
		if (r < '\u0410' || r > '\u042F') && (r < '\u0430' || r > '\u044F') {
			return false
		}
	}
	return true
}
