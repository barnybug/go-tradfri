package tradfri

func KelvinToMired(k int) int {
	mired := 1000000 / k
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
	return 1000000 / mired
}

func PercentageToDim(p int) int {
	dim := p * DimMax / 100
	if dim < DimMin {
		dim = DimMin
	} else if dim > DimMax {
		dim = DimMax
	}
	return dim
}

func MsToDuration(ms int) int {
	return ms / 100
}
