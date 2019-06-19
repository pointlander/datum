package main

import (
	"flag"
	"fmt"

	"github.com/pointlander/datum/mnist"
)

var (
	image  = flag.Int("image", 0, "the image to output")
	set    = flag.String("set", "train", "the image set")
	server = flag.Bool("server", false, "start the server")
)

func main() {
	flag.Parse()

	datum, err := mnist.Load()
	if err != nil {
		panic(err)
	}

	if *server {
		address := ":8080"
		fmt.Printf("starting server at %s\n", address)
		err := datum.Server(address)
		if err != nil {
			panic(err)
		}
		return
	}

	s := datum.Train
	if *set == "test" {
		s = datum.Test
	}
	err = s.WriteImage(*image, "./")
	if err != nil {
		panic(err)
	}
}
