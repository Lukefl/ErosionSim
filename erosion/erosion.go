package erosion

import (
	"erosion/heightmap"
	"math"
	"math/rand"
)

const (
	gravity         = 4.0
	erosionRate     = 0.3
	depositionRate  = 0.3
	capacityFactor  = 0.01
	evaporationRate = 0.02
	maxErosion      = 512.0
)

func Erode(m *heightmap.Map, droplets, steps int, seed int64) {
	for i := 0; i < droplets; i++ {
		erodeDroplet(m, steps, seed+int64(i))
	}
}

func erodeDroplet(m *heightmap.Map, steps int, seed int64) {
	if m.Width == 0 || m.Height == 0 || steps <= 0 {
		return
	}

	water := 1.0
	sediment := 0.0
	speed := 1.0

	random := rand.New(rand.NewSource(seed))
	x := random.Intn(m.Width)
	y := random.Intn(m.Height)

	for i := 0; i < steps; i++ {
		if !inBounds(m, x, y) || water <= 0 {
			return
		}

		dx, dy, found, steepness := getDownHill(m, x, y)
		if !found {
			deposit(m, x, y, sediment)
			return
		}

		distance := math.Sqrt(float64(dx*dx + dy*dy))
		deltaHeight := steepness * distance
		capacity := math.Max(deltaHeight*speed*water*capacityFactor, 1)

		if sediment > capacity {
			depositionAmount := (sediment - capacity) * depositionRate
			deposit(m, x, y, depositionAmount)
			sediment -= depositionAmount
		} else {
			erosionAmount := math.Min((capacity-sediment)*erosionRate, maxErosion)
			erosionAmount = math.Min(erosionAmount, float64(m.At(x, y)))
			remove(m, x, y, erosionAmount)
			sediment += erosionAmount
		}

		x += dx
		y += dy
		speed = math.Sqrt(math.Max(0, speed*speed+deltaHeight*gravity))
		water *= 1 - evaporationRate
	}

	deposit(m, x, y, sediment)
}

func remove(m *heightmap.Map, x, y int, amount float64) {
	if !inBounds(m, x, y) || amount <= 0 {
		return
	}

	newHeight := math.Max(float64(m.At(x, y))-amount, 0)
	m.Set(x, y, uint16(newHeight))
}

func deposit(m *heightmap.Map, x, y int, amount float64) {
	if !inBounds(m, x, y) || amount <= 0 {
		return
	}

	newHeight := math.Min(float64(m.At(x, y))+amount, math.MaxUint16)
	m.Set(x, y, uint16(newHeight))
}

func getDownHill(m *heightmap.Map, x, y int) (dx, dy int, found bool, steepness float64) {
	center := float64(m.At(x, y))

	for oy := -1; oy <= 1; oy++ {
		for ox := -1; ox <= 1; ox++ {
			if ox == 0 && oy == 0 {
				continue
			}

			nx := x + ox
			ny := y + oy

			if !inBounds(m, nx, ny) {
				continue
			}

			height := float64(m.At(nx, ny))

			if height >= center {
				continue
			}

			// Distance to the neighboring cell.
			distance := math.Sqrt(float64(ox*ox + oy*oy))

			// Height drop per unit of horizontal distance.
			slope := (center - height) / distance

			if !found || slope > steepness {
				dx = ox
				dy = oy
				steepness = slope
				found = true
			}
		}
	}

	return dx, dy, found, steepness
}

func inBounds(m *heightmap.Map, x, y int) bool {
	return x >= 0 && x < m.Width &&
		y >= 0 && y < m.Height
}
