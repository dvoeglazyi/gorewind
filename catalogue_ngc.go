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

// New General Catalogue of Nebulae and Clusters of Stars.
// https://cdsarc.unistra.fr/viz-bin/cat/VII/118

func ReadNGCCatalogue(path, namesPath string) ([]*AstronomicalObject, error) {
	namesMap, err := readNGCNames(namesPath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []*AstronomicalObject
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		record, err := getNGCRecord(line)
		if err != nil {
			return nil, err
		}
		if names := namesMap[ngcKey{catalogue: record.Catalogue, index: record.Index}]; len(names) > 0 {
			record.Name = names[0]
			record.AlternateNames = names[1:]
		}
		records = append(records, record)
	}
	return records, nil
}

func getNGCRecord(s string) (*AstronomicalObject, error) {
	longitudeHours, err := strconv.ParseUint(strings.TrimSpace(s[10:12]), 10, 64)
	if err != nil {
		return nil, err
	}
	longitudeMinutes, err := strconv.ParseFloat(strings.TrimSpace(s[13:17]), 64)
	if err != nil {
		return nil, err
	}

	latitudeDegrees, err := strconv.ParseInt(strings.TrimSpace(s[19:22]), 10, 64)
	if err != nil {
		return nil, err
	}
	latitudeMinutes, err := strconv.ParseUint(strings.TrimSpace(s[23:25]), 10, 64)
	if err != nil {
		return nil, err
	}

	var magnitude float64
	if magnitudeString := strings.TrimSpace(s[40:44]); magnitudeString != "" {
		if magnitude, err = strconv.ParseFloat(magnitudeString, 64); err != nil {
			return nil, err
		}
	}

	key, err := getNGCKey(s[:5])
	if err != nil {
		return nil, err
	}

	return &AstronomicalObject{
		Catalogue: key.catalogue,
		Index:     key.index,
		Designation: Designation{
			Constellation: strings.TrimSpace(s[29:32]),
		},
		Magnitude: magnitude,
		Coords:    NewClockCoords(uint(longitudeHours), longitudeMinutes, 0, int(latitudeDegrees), float64(latitudeMinutes), 0),
	}, nil
}

func readNGCNames(path string) (map[ngcKey][]string, error) {
	namesFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer namesFile.Close()

	namesMap := make(map[ngcKey][]string)
	reader := bufio.NewReader(namesFile)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		} else if len(s) < 41 {
			continue
		}
		key, err := getNGCKey(s[36:41])
		if err != nil {
			return nil, err
		}
		name := strings.TrimSpace(s[:35])
		namesMap[key] = append(namesMap[key], name)
	}
	return namesMap, nil
}

type ngcKey struct {
	catalogue string
	index     uint
}

func getNGCKey(s string) (ngcKey, error) {
	var key ngcKey
	if s[0] == 'I' {
		key.catalogue = "IC"
		s = s[1:]
	} else {
		key.catalogue = "NGC"
	}
	index, err := strconv.ParseUint(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return key, err
	}
	key.index = uint(index)
	return key, nil
}
