package main

import "math"

const ROUND_PRECISION = 5

func Round(v float64) float64 {
	f := math.Pow10(ROUND_PRECISION)
	return math.Round(f*v) / f
}

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func IntMin(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func IntSum(a []int) int {
	sum := 0
	for _, v := range a {
		sum += v
	}
	return sum
}

func Float64Min(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}

func Float64Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func Float64Sum(a []float64) float64 {
	var sum float64
	for _, v := range a {
		sum += v
	}
	return sum
}

func IntSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func IntSliceIndexOf(slice []int, val int) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

func IntSliceIntersect(a, b []int) []int {
	common := []int{}
	for _, v := range a {
		if IntSliceIndexOf(b, v) > -1 {
			common = append(common, v)
		}
	}
	return common
}
