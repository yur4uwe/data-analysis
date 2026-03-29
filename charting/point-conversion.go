package charting

func ToAnySlice(data []float64) []any {
	res := make([]any, len(data))
	for i, v := range data {
		val := v
		res[i] = &val
	}
	return res
}

func ToFloat64PtrSlice(data []float64) []*float64 {
	res := make([]*float64, len(data))
	for i, v := range data {
		val := v
		res[i] = &val
	}
	return res
}

func ToDataPointPtrSlice(data []DataPoint) []DataPoint {
	return data
}

func PointsToAnySlice(data []DataPoint) []any {
	res := make([]any, len(data))
	for i, v := range data {
		res[i] = v
	}
	return res
}

func AnyToPointsSlice(data []any) []DataPoint {
	res := make([]DataPoint, len(data))
	for i, v := range data {
		if v == nil {
			res[i] = DataPoint{X: float64(i), Y: nil}
			continue
		}
		if p, ok := v.(float64); ok {
			val := p
			res[i] = DataPoint{X: float64(i), Y: &val}
		} else if p, ok := v.(*float64); ok {
			res[i] = DataPoint{X: float64(i), Y: p}
		} else if p, ok := v.(*DataPoint); ok {
			res[i] = *p
		} else if p, ok := v.(DataPoint); ok {
			res[i] = p
		} else if hp, ok := v.(*HeatmapPoint); ok {
			res[i] = hp.DataPoint
		} else if hp, ok := v.(HeatmapPoint); ok {
			res[i] = hp.DataPoint
		}
	}
	return res
}

func IndexedDataPoints(data []float64) []DataPoint {
	res := make([]DataPoint, len(data))
	for i, v := range data {
		val := v
		res[i] = DataPoint{X: float64(i), Y: &val}
	}
	return res
}
