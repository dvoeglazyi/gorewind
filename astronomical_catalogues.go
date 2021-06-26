package gorewind

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

// TODO объединение данных каталогов не сделано.

func read(namesPath, bscPath, ngcPath, ngcNamesPath string) ([]*AstronomicalObject, error) {
	namesRecords, err := ReadNamesCatalogue(namesPath)
	if err != nil {
		return nil, err
	}

	namesByIndex := make(map[uint]*AstronomicalObject)
	namesByDesignation := make(map[Designation]*AstronomicalObject)

	for _, name := range namesRecords {
		if name.Catalogue == "HR" {
			namesByIndex[name.Index] = name
		}
		if name.Designation.BayerCode != rune(0) {
			namesByDesignation[Designation{BayerCode: name.Designation.BayerCode, Constellation: name.Designation.Constellation}] = name
		} else if name.Designation.FlamsteedCode != 0 {
			namesByDesignation[Designation{FlamsteedCode: name.Designation.FlamsteedCode, Constellation: name.Designation.Constellation}] = name
		}
	}

	bscRecords, err := ReadBSCCatalogue(bscPath)
	if err != nil {
		return nil, err
	}

	var result []*AstronomicalObject
	for _, record := range bscRecords {
		nameRecord := namesByIndex[record.Index]
		if nameRecord == nil {
			if record.Designation.BayerCode != rune(0) {
				nameRecord = namesByDesignation[Designation{BayerCode: record.Designation.BayerCode, Constellation: record.Designation.Constellation}]
			} else if record.Designation.FlamsteedCode != 0 {
				nameRecord = namesByDesignation[Designation{FlamsteedCode: record.Designation.FlamsteedCode, Constellation: record.Designation.Constellation}]
			}
		}
		if nameRecord == nil {
			continue
		}
		record.Name = nameRecord.Name
		record.LocalName = nameRecord.LocalName
		record.AlternateNames = nameRecord.AlternateNames
		result = append(result, record)
	}

	ngcRecords, err := ReadNGCCatalogue(ngcPath, ngcNamesPath)
	if err != nil {
		return nil, err
	}

	for _, record := range ngcRecords {
		if record.Name != "" {
			result = append(result, record)
		}
	}

	// TODO

	return result, nil
}
