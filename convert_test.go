package tradfri

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKelvinToMired(t *testing.T) {
	assert.Equal(t, 454, KelvinToMired(2100))
	assert.Equal(t, 454, KelvinToMired(2200))
	assert.Equal(t, 345, KelvinToMired(2900))
	assert.Equal(t, 250, KelvinToMired(4000))
	assert.Equal(t, 250, KelvinToMired(5000))
}

func TestMiredToKelvin(t *testing.T) {
	assert.Equal(t, 2203, MiredToKelvin(500))
	assert.Equal(t, 2203, MiredToKelvin(454))
	assert.Equal(t, 2899, MiredToKelvin(345))
	assert.Equal(t, 4000, MiredToKelvin(250))
	assert.Equal(t, 4000, MiredToKelvin(220))
}

func TestPercentageToDim(t *testing.T) {
	assert.Equal(t, 254, PercentageToDim(110))
	assert.Equal(t, 254, PercentageToDim(100))
	assert.Equal(t, 251, PercentageToDim(99))
	assert.Equal(t, 3, PercentageToDim(1))
	assert.Equal(t, 0, PercentageToDim(0))
	assert.Equal(t, 0, PercentageToDim(-10))
}

func TestDimToPercentage(t *testing.T) {
	assert.Equal(t, 100, DimToPercentage(265))
	assert.Equal(t, 100, DimToPercentage(254))
	assert.Equal(t, 99, DimToPercentage(251))
	assert.Equal(t, 1, DimToPercentage(3))
	assert.Equal(t, 0, DimToPercentage(0))
	assert.Equal(t, 0, DimToPercentage(-10))
}

func TestMsToDuration(t *testing.T) {
	assert.Equal(t, 10, MsToDuration(1000))
}

var hexRGBToColorXYDimTable = []struct {
	rgb       string
	x, y, dim int
}{
	{"ff0000", 44506, 21022, 80},
	{"00ff00", 11299, 48941, 170},
	{"0000ff", 8880, 2613, 12},
	{"ffffff", 20943, 21992, 255},
}

func TestHexRGBToColorXYDim(t *testing.T) {
	var x, y, dim int
	var err error

	for _, row := range hexRGBToColorXYDimTable {
		x, y, dim, err = HexRGBToColorXYDim(row.rgb)
		assert.Equal(t, row.x, x)
		assert.Equal(t, row.y, y)
		assert.Equal(t, row.dim, dim)
		assert.NoError(t, err)
	}
}

func TestKelvinToRGB(t *testing.T) {
	// pure white
	r, g, b := KelvinToRGB(6600)
	assert.Equal(t, 1., r)
	assert.Equal(t, 1., g)
	assert.Equal(t, 1., b)
}
