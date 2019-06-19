package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/dsnet/compress/bzip2"
	"github.com/petar/GoMNIST"
)

var MNIST = flag.Bool("mnist", false, "generate mnist gob file")

// Set is a set of training data
type Set struct {
	Width  int
	Height int
	Images [][]byte
	Labels []uint8
}

// Datum is the mnist training data set
type Datum struct {
	Train Set
	Test  Set
}

func main() {
	flag.Parse()

	if *MNIST {
		fmt.Println("loading data...")
		train, test, err := GoMNIST.Load("./")
		if err != nil {
			panic(err)
		}
		convertImages := func(images []GoMNIST.RawImage) [][]byte {
			output := make([][]byte, len(images))
			for i, image := range images {
				output[i] = []byte(image)
			}
			return output
		}
		convertLabels := func(labels []GoMNIST.Label) []uint8 {
			output := make([]uint8, len(labels))
			for i, label := range labels {
				output[i] = uint8(label)
			}
			return output
		}
		mnist := Datum{
			Train: Set{
				Width:  train.NCol,
				Height: train.NRow,
				Images: convertImages(train.Images),
				Labels: convertLabels(train.Labels),
			},
			Test: Set{
				Width:  train.NCol,
				Height: train.NRow,
				Images: convertImages(test.Images),
				Labels: convertLabels(test.Labels),
			},
		}

		fmt.Println("encoding data...")
		buffer := bytes.Buffer{}
		compress, err := bzip2.NewWriter(&buffer, &bzip2.WriterConfig{Level: bzip2.BestCompression})
		if err != nil {
			panic(err)
		}
		encoder := gob.NewEncoder(compress)
		encoder.Encode(mnist)
		err = compress.Close()
		if err != nil {
			panic(err)
		}

		fmt.Println("writing data...")
		out, err := os.Create("assets.go")
		if err != nil {
			panic(err)
		}
		defer out.Close()
		fmt.Fprintf(out, "package mnist\n\n")
		fmt.Fprintf(out, "var Asset = []byte(%+q)\n", buffer.Bytes())
	}
}
