package analysis

import (
	"math"
	"sort"
)

func Mean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func StdDev(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	mean := Mean(data)
	sumSquares := 0.0
	for _, v := range data {
		diff := v - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(data)))
}

func Variance(data []float64, mean float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, v := range data {
		diff := v - mean
		sumSquares += diff * diff
	}
	return sumSquares / float64(len(data))
}

func Median(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

func Mode(data []float64, precision int) float64 {
	if len(data) == 0 {
		return 0
	}
	counts := make(map[float64]int)
	factor := math.Pow10(precision)
	for _, v := range data {
		rounded := math.Round(v*factor) / factor
		counts[rounded]++
	}
	maxCount := 0
	var mode float64
	for v, count := range counts {
		if count > maxCount {
			maxCount = count
			mode = v
		}
	}
	return mode
}

func MinMax(data []float64) (min, max float64) {
	if len(data) == 0 {
		return 0, 0
	}
	min = data[0]
	max = data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return
}

func MSE(actual []float64, calculated []float64) float64 {
	var sumSq float64
	count := 0
	for i := 0; i < len(actual) && i < len(calculated); i++ {
		diff := actual[i] - calculated[i]
		sumSq += diff * diff
		count++
	}
	if count == 0 {
		return math.MaxFloat64
	}
	return sumSq / float64(count)
}
