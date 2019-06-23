// Copyright 2019 The Datum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mnist

import (
	"bytes"
	"compress/bzip2"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// Width is the width of the image
	Width = 28
	// Height is the height of the image
	Height = 28
)

// Image is a mnist image
type Image []byte

// ColorModel returns the color model of the image
func (i Image) ColorModel() color.Model {
	return color.GrayModel
}

// Bounds returns the bounds of the image
func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, Width, Height)
}

// At get the pixel at xy
func (i Image) At(x, y int) color.Color {
	return color.Gray{Y: i[y*Width+x]}
}

// Set is a set of training data
type Set struct {
	Width  int
	Height int
	Images []Image
	Labels []uint8
}

// WriteImage writes the image i to dir
func (s *Set) WriteImage(i int, dir string) error {
	name := fmt.Sprintf("image%d_%d.png", i, s.Labels[i])
	out, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, s.Images[i])
}

// Datum is the mnist training data set
type Datum struct {
	Train Set
	Test  Set
}

// Load loads the MNIST data set
func Load() (datum Datum, err error) {
	decoder := gob.NewDecoder(bzip2.NewReader(bytes.NewReader(Asset)))
	err = decoder.Decode(&datum)
	if err != nil {
		return datum, err
	}
	return datum, nil
}

// Server starts an image server
func (d *Datum) Server(address string) error {
	handler := func(set Set) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.URL.EscapedPath(), "/")
			length := len(parts)
			if length != 3 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			name := strings.Split(parts[length-1], ".")
			if len(name) != 2 {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			switch name[1] {
			case "png":
				i, err := strconv.Atoi(name[0])
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.Header().Add("Content-Type", "image/png")
				if i < 0 {
					i = -1
				}
				if i >= len(set.Images) {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				png.Encode(w, set.Images[i])
			case "html":
				i, err := strconv.Atoi(name[0])
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				w.Header().Add("Content-Type", "text/html")
				if i < 0 {
					i = -1
				}
				if i >= len(set.Images) {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				fmt.Fprintf(w, "<html>\n")
				fmt.Fprintf(w, " <head><title>image %d</title></head>\n", i)
				fmt.Fprintf(w, " <body>\n")
				fmt.Fprintf(w, "  <b>image %d: %d</b><br/>\n", i, set.Labels[i])
				fmt.Fprintf(w, "  <img src=\"%d.png\"/>", i)
				fmt.Fprintf(w, " </body>\n")
				fmt.Fprintf(w, "</html>\n")
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}
	}
	http.HandleFunc("/train/", handler(d.Train))
	http.HandleFunc("/test/", handler(d.Test))

	return http.ListenAndServe(address, nil)
}
