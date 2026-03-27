package forecasting

import (
	"math"
	"sort"
)

func tomorrowAsToday(y_t float64) float64 {
	return y_t
}

func trend(y_t, y_t_minus_1 float64) float64 {
	return y_t + (y_t - y_t_minus_1)
}

func relativeTrend(y_t, y_t_minus_1 float64) float64 {
	if y_t_minus_1 == 0 {
		return y_t
	}
	return y_t * (y_t / y_t_minus_1)
}

func simpleAvg(history []float64) float64 {
	sum := 0.0
	for _, val := range history {
		sum += val
	}
	return sum / float64(len(history))
}

func slidingAvg(history []float64, window int) float64 {
	if len(history) == 0 {
		return 0
	}
	if window > len(history) {
		window = len(history)
	}
	sum := 0.0
	for i := 0; i < window; i++ {
		sum += history[len(history)-1-i]
	}
	return sum / float64(window)
}

func exponentialAvg(y_t float64, prev_forecast float64, alpha float64) float64 {
	return alpha*y_t + (1-alpha)*prev_forecast
}

func CalculateMean(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func CalculateStdDev(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	mean := CalculateMean(data)
	sumSquares := 0.0
	for _, v := range data {
		diff := v - mean
		sumSquares += diff * diff
	}
	return math.Sqrt(sumSquares / float64(len(data)))
}

func CalculateMedian(data []float64) float64 {
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

func CalculateMode(data []float64, precision int) float64 {
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

func CalculateMinMax(data []float64) (min, max float64) {
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

func CalculateMSE(rates []float64, forecast []any) float64 {
	var sumSq float64
	count := 0
	for i := 0; i < len(rates) && i < len(forecast); i++ {
		if forecast[i] == nil {
			continue
		}
		val, ok := forecast[i].(float64)
		if !ok {
			continue
		}
		diff := rates[i] - val
		sumSq += diff * diff
		count++
	}
	if count == 0 {
		return math.MaxFloat64
	}
	return sumSq / float64(count)
}
