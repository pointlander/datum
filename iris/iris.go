// Copyright 2019 The Datum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package iris

import (
	"bytes"
	"compress/bzip2"
	"encoding/gob"
	"strconv"
)

var Labels = map[string]int{
	"Iris-setosa":     0,
	"Iris-versicolor": 1,
	"Iris-virginica":  2,
}

// Iris is a iris
type Iris struct {
	Measures []float64
	Label    string
}

// Datum is the Iris data
type Datum struct {
	Fisher []Iris
	Bezdek []Iris
}

// Load loads the Iris data set
func Load() (datum Datum, err error) {
	load := func(asset []byte, out *[]Iris) error {
		var data [][]string
		decoder := gob.NewDecoder(bzip2.NewReader(bytes.NewReader(asset)))
		err = decoder.Decode(&data)
		if err != nil {
			return err
		}
		for _, row := range data {
			iris := Iris{
				Measures: make([]float64, 4),
				Label:    row[4],
			}
			for i, measure := range row[:4] {
				value, err := strconv.ParseFloat(measure, 64)
				if err != nil {
					return err
				}
				iris.Measures[i] = value
			}
			*out = append(*out, iris)
		}
		return nil
	}
	err = load(AssetFisher, &datum.Fisher)
	if err != nil {
		return datum, err
	}
	err = load(AssetBezdek, &datum.Bezdek)
	if err != nil {
		return datum, err
	}

	return datum, nil
}
