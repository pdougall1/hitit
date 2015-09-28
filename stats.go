package main

import (
	"math"
	"sort"
)

func GetStats(numbers []float64) (stats statistics) {
	stats.numbers = numbers
	sort.Float64s(stats.numbers)
	stats.mean = sum(numbers) / float64(len(numbers))
	stats.median = median(numbers)
	stats.modes = mode(numbers)
	stats.stdDev = stdDev(numbers, stats.mean)
	stats.quickest = min(numbers)
	stats.slowest = max(numbers)
	return stats
}

func min(numbers []float64) float64 {
	sort.Float64s(numbers)
	return numbers[0]
}

func max(numbers []float64) float64 {
	sort.Sort(sort.Reverse(sort.Float64Slice(numbers)))
	return numbers[0]
}

func sum(numbers []float64) (total float64) {
	for _, x := range numbers {
		total += x
	}
	return total
}

func median(numbers []float64) float64 {
	middle := len(numbers) / 2
	result := numbers[middle]
	if len(numbers)%2 == 0 {
		result = (result + numbers[middle-1]) / 2
	}
	return result
}

func mode(numbers []float64) (modes []float64) {
	frequencies := make(map[float64]int, len(numbers))
	highestFrequency := 0
	for _, x := range numbers {
		frequencies[x]++
		if frequencies[x] > highestFrequency {
			highestFrequency = frequencies[x]
		}
	}
	for x, frequency := range frequencies {
		if frequency == highestFrequency {
			modes = append(modes, x)
		}
	}
	if highestFrequency == 1 || len(modes) == len(numbers) {
		modes = modes[:0] // Or: modes = []float64{}
	}
	sort.Float64s(modes)
	return modes
}

func stdDev(numbers []float64, mean float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += math.Pow(number-mean, 2)
	}
	variance := total / float64(len(numbers)-1)
	return math.Sqrt(variance)
}

type statistics struct {
	numbers  []float64
	mean     float64
	median   float64
	modes    []float64
	stdDev   float64
	quickest float64
	slowest  float64
}
