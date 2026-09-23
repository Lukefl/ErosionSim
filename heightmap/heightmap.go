package heightmap

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
)

// Map stores one 16-bit elevation value for every cell. Zero is low and 65535
// is high.
type Map struct {
	Width  int
	Height int
	Cells  []uint16
}

func New(width, height int) *Map {
	if width <= 0 || height <= 0 {
		panic("heightmap dimensions must be positive")
	}

	return &Map{
		Width:  width,
		Height: height,
		Cells:  make([]uint16, width*height),
	}
}

func (m *Map) At(x, y int) uint16 {
	return m.Cells[y*m.Width+x]
}

func (m *Map) Set(x, y int, value uint16) {
	m.Cells[y*m.Width+x] = value
}

func (m *Map) Clone() *Map {
	clone := New(m.Width, m.Height)
	copy(clone.Cells, m.Cells)
	return clone
}

// Generate creates a deterministic starting landscape using the default depth.
func Generate(width, height int, seed int64) *Map {
	return GenerateWithDepth(width, height, seed, 4)
}

// GenerateWithDepth creates a deterministic mountain landscape by summing
// ridged sine-wave layers. Each layer doubles the frequency and reduces the
// amplitude, adding smaller foothills and peaks.
func GenerateWithDepth(width, height int, seed int64, depth int) *Map {
	if depth <= 0 {
		panic("heightmap depth must be positive")
	}

	m := New(width, height)
	random := rand.New(rand.NewSource(seed))
	phases := make([]float64, depth*4)
	for index := range phases {
		phases[index] = random.Float64() * 2 * math.Pi
	}

	totalAmplitude := 0.0
	amplitude := 1.0
	for layer := 0; layer < depth; layer++ {
		totalAmplitude += amplitude
		amplitude *= 0.5
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			normalizedX := float64(x) / float64(width-1)
			normalizedY := float64(y) / float64(height-1)
			mountains := 0.0
			detail := 0.0
			frequency := 1.0
			amplitude := 1.0
			for layer := 0; layer < depth; layer++ {
				phase := phases[layer*4:]
				signal := (math.Sin(normalizedX*math.Pi*2.2*frequency+phase[0]) +
					math.Sin(normalizedY*math.Pi*2.7*frequency+phase[1]) +
					math.Sin((normalizedX+normalizedY)*math.Pi*3.4*frequency+phase[2]) +
					math.Sin((normalizedX-normalizedY)*math.Pi*2.9*frequency+phase[3])) / 4
				ridge := 1 - math.Abs(signal)
				mountains += amplitude * ridge * ridge
				detail += amplitude * signal
				frequency *= 2
				amplitude *= 0.5
			}

			mountainValue := mountains / totalAmplitude
			detailValue := 0.5 + detail/(2*totalAmplitude)
			value := 0.15*detailValue + 0.85*math.Pow(mountainValue, 1.8)

			value = math.Max(0, math.Min(1, value))
			m.Set(x, y, uint16(value*math.MaxUint16))
		}
	}

	return m
}

func (m *Map) SavePNG(path string) error {
	imageMap := image.NewGray16(image.Rect(0, 0, m.Width, m.Height))
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			imageMap.SetGray16(x, y, color.Gray16{Y: m.At(x, y)})
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, imageMap)
}

// SaveDifferencePNG writes a diagnostic image: red is material removed,
// blue is material deposited, and black is unchanged terrain.
func (m *Map) SaveDifferencePNG(before *Map, path string) error {
	if m.Width != before.Width || m.Height != before.Height {
		return fmt.Errorf("heightmap dimensions do not match")
	}

	maxDifference := uint16(0)
	for index, afterHeight := range m.Cells {
		beforeHeight := before.Cells[index]
		difference := afterHeight - beforeHeight
		if beforeHeight > afterHeight {
			difference = beforeHeight - afterHeight
		}
		if difference > maxDifference {
			maxDifference = difference
		}
	}

	imageMap := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			index := y*m.Width + x
			beforeHeight := before.Cells[index]
			afterHeight := m.Cells[index]
			if maxDifference == 0 || beforeHeight == afterHeight {
				imageMap.SetRGBA(x, y, color.RGBA{A: 255})
				continue
			}

			var difference uint16
			pixel := color.RGBA{A: 255}
			if afterHeight > beforeHeight {
				difference = afterHeight - beforeHeight
				pixel.B = uint8(uint32(difference) * 255 / uint32(maxDifference))
			} else {
				difference = beforeHeight - afterHeight
				pixel.R = uint8(uint32(difference) * 255 / uint32(maxDifference))
			}
			imageMap.SetRGBA(x, y, pixel)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, imageMap)
}
