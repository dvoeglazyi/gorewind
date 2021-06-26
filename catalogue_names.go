package gorewind

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

// Astro Catalogue. Каталог небесных тел с названиями на русском языке.
// https://github.com/dvoeglazyi/astrocat

func ReadNamesCatalogue(path string) ([]*AstronomicalObject, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ','
	csvReader.FieldsPerRecord = 11

	var result []*AstronomicalObject
	for {
		fields, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		record, err := readNamesCatalogueRecord(fields)
		if err != nil {
			return nil, err
		}

		result = append(result, record)
	}
	return result, nil
}

func readNamesCatalogueRecord(fields []string) (*AstronomicalObject, error) {
	record := AstronomicalObject{
		Name:      fields[0],
		LocalName: fields[1],
		Designation: Designation{
			Constellation: fields[3],
		},
		AlternateNames: strings.Split(fields[10], ";"),
	}
	if err := record.Designation.SetCode(fields[2]); err != nil {
		return nil, err
	}
	if fields[4] != "" {
		inSystemIndex, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return nil, err
		}
		record.Designation.InSystemIndex = uint(inSystemIndex)
	}
	if fields[5] != "" {
		split := strings.Split(fields[5], " ")
		if len(split) < 2 {
			return nil, errors.New("invalid catalogue index")
		}
		record.Catalogue = split[0]
		index, err := strconv.ParseUint(split[1], 10, 64)
		if err != nil {
			return nil, err
		}
		record.Index = uint(index)
	}
	if fields[6] != "" {
		magnitude, err := strconv.ParseFloat(fields[6], 64)
		if err != nil {
			return nil, err
		}
		record.Magnitude = magnitude
	}
	if fields[7] != "" && fields[8] != "" {
		longitude, err := strconv.ParseFloat(fields[7], 64)
		if err != nil {
			return nil, err
		}
		latitude, err := strconv.ParseFloat(fields[8], 64)
		if err != nil {
			return nil, err
		}
		record.Coords = NewCoordsFromDegrees(longitude, latitude)
	}
	if fields[9] != "" {
		radius, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			return nil, err
		}
		record.Coords.Radius = float64(radius)
	}
	return &record, nil
}
