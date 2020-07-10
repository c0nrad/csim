package csim

import "math"

func Factorial(i int) float64 {
	out := 1.0
	for i > 1 {
		out *= float64(i)
		i--
	}

	return out
}

func LegendrePoly(l int) Func {
	prefix := 1 / (math.Pow(2, float64(l)) * Factorial(l))

	inner := func(x float64) float64 {
		out := 1.0
		for i := 0; i < l; i++ {
			out *= x*x - 1
		}
		return out
	}

	out := inner
	for i := 0; i < l; i++ {
		out = D(out)
	}
	return func(x float64) float64 {
		return prefix * out(x)
	}
}
