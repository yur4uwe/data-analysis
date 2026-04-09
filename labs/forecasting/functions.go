package forecasting

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
