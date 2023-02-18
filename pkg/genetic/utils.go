package genetic

import "math"

func Max(s []int) int {
	m := math.MinInt
	for _, v := range s {
		if v > m {
			m = v
		}
	}
	return m
}

func Min(s []int) int {
	m := math.MaxInt
	for _, v := range s {
		if v < m {
			m = v
		}
	}
	return m
}

func Avg(s []int) int {
	var sum int
	for _, v := range s {
		sum += v
	}
	return sum / len(s)
}

func Std(s []int) int {
	s2 := make([]int, len(s))
	for i := range s2 {
		s2[i] = s[i]*s[i]
	}
	return int(math.Sqrt(float64(Avg(s2) - Avg(s))))
}