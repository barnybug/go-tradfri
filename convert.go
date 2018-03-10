package tradfri

func KelvinToMired(k int) int {
	mired := round(1000000 / float64(k))
	if mired < MiradsMin {
		mired = MiradsMin
	} else if mired > MiradsMax {
		mired = MiradsMax
	}
	return mired
}

func MiredToKelvin(mired int) int {
	if mired < MiradsMin {
		mired = MiradsMin
	} else if mired > MiradsMax {
		mired = MiradsMax
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
