package heightmap

import (
	"image"
	"image/png"
	"os"
	"slices"
	"testing"
)

func TestGenerateProducesValidHeightValues(t *testing.T) {
	m := Generate(32, 24, 42)

	if len(m.Cells) != 32*24 {
		t.Fatalf("got %d cells, want %d", len(m.Cells), 32*24)
	}

	if m.Cells[0] == m.Cells[len(m.Cells)-1] {
		t.Fatal("generated terrain should contain elevation variation")
	}

}

func TestGenerateDepthChangesTerrainDetail(t *testing.T) {
	shallow := GenerateWithDepth(32, 24, 42, 1)
	deep := GenerateWithDepth(32, 24, 42, 6)

	if slices.Equal(shallow.Cells, deep.Cells) {
		t.Fatal("different depths should produce different terrain")
	}
}

func TestSavePNG(t *testing.T) {
	path := t.TempDir() + "/heightmap.png"
	if err := Generate(8, 6, 42).SavePNG(path); err != nil {
		t.Fatalf("save PNG: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open PNG: %v", err)
	}
	defer file.Close()

	decoded, err := png.Decode(file)
	if err != nil {
		t.Fatalf("decode PNG: %v", err)
	}
	if got := decoded.Bounds().Size(); got.X != 8 || got.Y != 6 {
		t.Fatalf("got image size %v, want 8x6", got)
	}

	if _, ok := decoded.(*image.Gray16); !ok {
		t.Fatalf("got %T, want *image.Gray16", decoded)
	}
}

func TestSaveDifferencePNG(t *testing.T) {
	before := New(2, 1)
	before.Set(0, 0, 1000)
	after := before.Clone()
	after.Set(1, 0, 4000)
	after.Set(0, 0, 0)

	path := t.TempDir() + "/difference.png"
	if err := after.SaveDifferencePNG(before, path); err != nil {
		t.Fatalf("save difference PNG: %v", err)
	}

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open difference PNG: %v", err)
	}
	defer file.Close()

	decoded, err := png.Decode(file)
	if err != nil {
		t.Fatalf("decode difference PNG: %v", err)
	}
	if red, _, _, _ := decoded.At(0, 0).RGBA(); red == 0 {
		t.Fatal("height decrease should be visible in the difference image")
	}
	if _, _, blue, _ := decoded.At(1, 0).RGBA(); blue == 0 {
		t.Fatal("height increase should be visible in the difference image")
	}
}
