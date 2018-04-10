package tradfri

import (
	"encoding/hex"
	"errors"
	"math"
)

func KelvinToMired(k int) int {
	mired := round(1000000 / float64(k))
	if mired < MiredMin {
		mired = MiredMin
	} else if mired > MiredMax {
		mired = MiredMax
	}
	return mired
}

func MiredToKelvin(mired int) int {
	if mired < MiredMin {
		mired = MiredMin
	} else if mired > MiredMax {
		mired = MiredMax
	}
	return round(1000000 / float64(mired))
}

func round(f float64) int {
	if f < -0.5 {
		return int(f - 0.5)
	}
	if f > 0.5 {
		return int(f + 0.5)
	}
	return 0
}

func PercentageToDim(p int) int {
	dim := round(float64(p) * DimMax / 100)
	if dim < DimMin {
		dim = DimMin
	} else if dim > DimMax {
		dim = DimMax
	}
	return dim
}

func DimToPercentage(dim int) int {
	p := round(float64(dim) * 100 / DimMax)
	if p > 100 {
		p = 100
	}
	return p
}

func MsToDuration(ms int) int {
	return ms / 100
}

// Gamma correction of rgb component
func norm(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	} else {
		return math.Pow((v+0.055)/1.055, 2.4)
	}
}

// Convert sRGB D65 -> xy colour space
func RGBToColorXYDim(rgb string) (x int, y int, dim int, err error) {
	if len(rgb) != 6 {
		err = errors.New("Incorrect length color hex string")
		return
	}
	var s []byte
	s, err = hex.DecodeString(rgb)
	if err != nil {
		return
	}
	// Gamma correct sRGB -> sRGB'
	r := norm(float64(s[0]) / 255)
	g := norm(float64(s[1]) / 255)
	b := norm(float64(s[2]) / 255)
	// Wide RGB D65 conversion formula
	X := r*0.664511 + g*0.154324 + b*0.162028
	Y := r*0.313881 + g*0.668433 + b*0.047685
	Z := r*0.000088 + g*0.072310 + b*0.986039
	// Convert XYZ -> xy
	x = int(X / (X + Y + Z) * 65535)
	y = int(Y / (X + Y + Z) * 65535)
	if Y > 1 {
		Y = 1
	}
	dim = int(Y * 255)
	return
}
