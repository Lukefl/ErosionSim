package main

import (
	"flag"
	"fmt"
	"log"

	"erosion/heightmap"
)

func main() {
	width := flag.Int("width", 512, "heightmap width in cells")
	height := flag.Int("height", 512, "heightmap height in cells")
	depth := flag.Int("depth", 4, "number of sine-wave detail layers")
	seed := flag.Int64("seed", 42, "random seed")
	output := flag.String("output", "heightmap.png", "output PNG path")
	flag.Parse()

	mapData := heightmap.GenerateWithDepth(*width, *height, *seed, *depth)
	if err := mapData.SavePNG(*output); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("wrote %dx%d heightmap to %s\n", *width, *height, *output)
}
