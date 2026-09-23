package main

import (
	"flag"
	"fmt"
	"log"

	"erosion/erosion"
	"erosion/heightmap"
)

func main() {
	width := flag.Int("width", 512, "heightmap width in cells")
	height := flag.Int("height", 512, "heightmap height in cells")
	depth := flag.Int("depth", 4, "number of sine-wave detail layers")
	seed := flag.Int64("seed", 42, "random seed")
	droplets := flag.Int("droplets", 10000, "number of droplets")
	steps := flag.Int("steps", 30, "number of steps")
	beforeOutput := flag.String("beforeOutput", "heightmap.png", "before output PNG path")
	afterOutput := flag.String("afterOutput", "eroded.png", "after output PNG path")
	differenceOutput := flag.String("differenceOutput", "difference.png", "erosion difference PNG path")
	flag.Parse()

	mapData := heightmap.GenerateWithDepth(*width, *height, *seed, *depth)
	beforeMap := mapData.Clone()
	if err := mapData.SavePNG(*beforeOutput); err != nil {
		log.Fatal(err)
	}
	erosion.Erode(mapData, *droplets, *steps, *seed)
	if err := mapData.SavePNG(*afterOutput); err != nil {
		log.Fatal(err)
	}
	if err := mapData.SaveDifferencePNG(beforeMap, *differenceOutput); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %dx%d before map to %s\n", *width, *height, *beforeOutput)
	fmt.Printf("\nwrote %dx%d eroded map to %s\n", *width, *height, *afterOutput)
	fmt.Printf("wrote erosion difference to %s\n", *differenceOutput)
}
