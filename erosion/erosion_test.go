package erosion

import (
	"erosion/heightmap"
	"testing"
)

func TestErodeChangesTerrain(t *testing.T) {
	m := heightmap.New(24, 24)
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			m.Set(x, y, uint16(50000-x*1000))
		}
	}

	before := append([]uint16(nil), m.Cells...)
	Erode(m, 100, 20, 42)

	changed := false
	for index, value := range m.Cells {
		if value != before[index] {
			changed = true
			break
		}
	}
	if !changed {
		t.Fatal("erosion should change at least one height")
	}
}

func TestErodeIsDeterministic(t *testing.T) {
	first := heightmap.New(24, 24)
	second := heightmap.New(24, 24)
	for y := 0; y < first.Height; y++ {
		for x := 0; x < first.Width; x++ {
			height := uint16(50000 - x*1000)
			first.Set(x, y, height)
			second.Set(x, y, height)
		}
	}

	Erode(first, 100, 20, 42)
	Erode(second, 100, 20, 42)

	for index := range first.Cells {
		if first.Cells[index] != second.Cells[index] {
			t.Fatalf("cell %d differs for the same seed", index)
		}
	}
}
