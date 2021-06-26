package gorewind

import (
	"strconv"
	"strings"
)

// gorewind. Golang library for working with astronomical and geographical objects in spherical coordinates system.
// Библиотека для работы с астрономическими и географическими объектами в сферической системе координат.
// Copyright © 2021 Dvoeglazyi
// License: http://github.com/dvoeglazyi/gorewind/LICENSE

type Designation struct {
	BayerCode        rune
	FlamsteedCode    uint // from 1 if exist
	VariableStarCode string
	InSystemIndex    uint // from 1 if exist
	Constellation    string
}

func (d *Designation) GetCode() string {
	if d.BayerCode != 0 {
		return string(d.BayerCode)
	} else if d.FlamsteedCode != 0 {
		return strconv.FormatUint(uint64(d.FlamsteedCode), 10)
	}
	return d.VariableStarCode
}

func (d *Designation) SetCode(code string) error {
	if code == "" {
		return nil
	}
	// Bayer Designation:
	// Greek letter as in-constellation index
	runes := []rune(code)
	if len(runes) == 1 {
		for _, letter := range greekLetters {
			if runes[0] == letter.rune {
				d.BayerCode = letter.rune
				return nil
			}
		}
	}
	// Variable Star Designation:
	// Letter-digit code as in-constellation index
	if runes[0] < '0' || runes[0] > '9' {
		d.VariableStarCode = code
		return nil
	}
	// Flamsteed Designation:
	// Number as in-constellation index
	num, err := strconv.ParseUint(code, 10, 64)
	if err != nil {
		return nil
	}
	d.FlamsteedCode = uint(num)
	return nil
}

type letter struct {
	rune      rune
	shortName string
}

var greekLetters = []letter{
	{rune: 'α', shortName: "alp"},
	{rune: 'β', shortName: "bet"},
	{rune: 'γ', shortName: "gam"},
	{rune: 'δ', shortName: "del"},
	{rune: 'ε', shortName: "eps"},
	{rune: 'ζ', shortName: "zet"},
	{rune: 'η', shortName: "eta"},
	{rune: 'θ', shortName: "the"},
	{rune: 'ι', shortName: "iot"},
	{rune: 'κ', shortName: "kap"},
	{rune: 'λ', shortName: "lam"},
	{rune: 'μ', shortName: "mu"},
	{rune: 'ν', shortName: "nu"},
	{rune: 'ξ', shortName: "xi"},
	{rune: 'π', shortName: "pi"},
	{rune: 'ρ', shortName: "rho"},
	{rune: 'σ', shortName: "sig"},
	{rune: 'τ', shortName: "tau"},
	{rune: 'υ', shortName: "ups"},
	{rune: 'φ', shortName: "phi"},
	{rune: 'χ', shortName: "chi"},
	{rune: 'ψ', shortName: "psi"},
	{rune: 'ω', shortName: "omi"},
}

func GetBayerRune(code string) rune {
	code = strings.ToLower(code)
	for _, letter := range greekLetters {
		if letter.shortName == code {
			return letter.rune
		}
	}
	return rune(0)
}
