package gorewind

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

// Yale Catalogue of Bright Stars.
// http://cdsarc.u-strasbg.fr/viz-bin/Cat?V/50

func ReadBSCCatalogue(path string) ([]*AstronomicalObject, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []*AstronomicalObject
	reader := bufio.NewReader(file)
	for {
		fields, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		record, err := getRecord(fields)
		if err != nil {
			return nil, err
		} else if record == nil {
			continue
		}
		records = append(records, record)
	}
	return records, nil
}

func getDesignation(s string) (Designation, error) {
	result := Designation{
		Constellation: strings.TrimSpace(s[7:]),
		BayerCode:     GetBayerRune(strings.TrimSpace(s[3:6])),
	}
	if s[6] != ' ' {
		index, err := strconv.ParseUint(string(s[6]), 10, 64)
		if err != nil {
			return result, err
		}
		result.InSystemIndex = uint(index)
	}
	if code := strings.TrimSpace(s[:3]); code != "" {
		index, err := strconv.ParseUint(code, 10, 64)
		if err != nil {
			return result, err
		}
		result.FlamsteedCode = uint(index)
	}
	return result, nil
}

func getRecord(s string) (*AstronomicalObject, error) {
	index, err := strconv.ParseUint(strings.TrimSpace(s[:4]), 10, 64)
	if err != nil {
		return nil, err
	}
	longitudeHoursString := strings.TrimSpace(s[75:77])
	if longitudeHoursString == "" {
		// catalogue has some records without coordinates
		// that records will be skipped
		return nil, nil
	}

	longitudeHours, err := strconv.ParseUint(longitudeHoursString, 10, 64)
	if err != nil {
		return nil, err
	}
	longitudeMinutes, err := strconv.ParseUint(strings.TrimSpace(s[77:79]), 10, 64)
	if err != nil {
		return nil, err
	}
	longitudeSeconds, err := strconv.ParseFloat(strings.TrimSpace(s[79:83]), 64)
	if err != nil {
		return nil, err
	}

	latitudeDegrees, err := strconv.ParseInt(strings.TrimSpace(s[83:86]), 10, 64)
	if err != nil {
		return nil, err
	}
	latitudeMinutes, err := strconv.ParseUint(strings.TrimSpace(s[86:88]), 10, 64)
	if err != nil {
		return nil, err
	}
	latitudeSeconds, err := strconv.ParseUint(strings.TrimSpace(s[88:90]), 10, 64)
	if err != nil {
		return nil, err
	}

	var magnitude float64
	if magnitudeString := strings.TrimSpace(s[102:107]); magnitudeString != "" {
		if magnitude, err = strconv.ParseFloat(magnitudeString, 64); err != nil {
			return nil, err
		}
	}

	result := AstronomicalObject{
		Catalogue: "HR",
		Index:     uint(index),
		Coords:    NewClockCoords(uint(longitudeHours), float64(longitudeMinutes), longitudeSeconds, int(latitudeDegrees), float64(latitudeMinutes), float64(latitudeSeconds)),
		Magnitude: magnitude,
	}
	if result.Designation, err = getDesignation(s[4:14]); err != nil {
		return nil, err
	}
	return &result, nil
}

func readBSCNotes(path string) (map[uint]BSCNote, error) {
	namesFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer namesFile.Close()

	notesMap := make(map[uint]BSCNote)
	reader := bufio.NewReader(namesFile)
	for {
		s, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if len(s) < 41 {
			continue
		}
		index, err := strconv.ParseUint(strings.TrimSpace(s[1:5]), 10, 64)
		if err != nil {
			return nil, err
		}
		noteIndex, err := strconv.ParseUint(strings.TrimSpace(s[5:7]), 10, 64)
		if err != nil {
			return nil, err
		}
		key := strings.TrimSpace(s[7:11])
		value := s[12:]
		notesMap[uint(index)].add(key, uint(noteIndex), value)
	}
	return notesMap, nil
}

type BSCNote map[string]string

func (n BSCNote) add(key string, index uint, value string) {
	if index > 1 {
		n[key] += value
		return
	}
	n[key] = value
}
