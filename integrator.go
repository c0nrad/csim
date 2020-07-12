package csim

import (
	"math"
)

type Integrator struct {
	F Func
	N int
}

func (i Integrator) H(a, b float64) float64 {
	return (b - a) / float64(i.N)
}

func (i Integrator) Rectangular(a, b float64) float64 {
	out := 0.0
	x := a
	i.N++
	h := i.H(a, b)
	for step := 0; step < i.N; step++ {
		out += i.F(x) * h
		x += h
	}
	i.N--
	return out
}

func (i Integrator) Trapezoidal(a, b float64) float64 {
	out := 0.0
	x := a

	h := i.H(a, b)

	for step := 0; step < i.N; step++ {
		out += i.F(x) + i.F(x+h)
		x += h
	}
	out *= h / 2
	return out
}

func (i Integrator) Simpson(a, b float64) float64 {
	if i.N%2 == 1 || i.N <= 2 {
		return 0
	}

	out := 0.0

	// So, we'll be including the very last piece, meaning we need to make h slightly wider

	h := i.H(a, b)

	x := a + h

	for step := 1; step < i.N; step++ {
		if step%2 == 1 {
			out += (4 * i.F(x))
		} else {
			out += (2 * i.F(x))
		}
		x += h
	}
	out += i.F(b) + i.F(a)
	out *= h / 3
	return out
}

type LegendrePoint struct {
	N  int
	Xi []float64
	Wi []float64
}

var LegendrePoints = []LegendrePoint{
	{N: 0},
	{N: 1, Xi: []float64{0}, Wi: []float64{0}},
	{N: 2,
		Xi: []float64{1.0 / math.Sqrt(3), -1.0 / math.Sqrt(3)},
		Wi: []float64{1.0, 1.0},
	},
	{N: 2,
		Xi: []float64{0, math.Sqrt(3.0 / 5.0), -math.Sqrt(3.0 / 5.0)},
		Wi: []float64{8.0 / 9.0, 5.0 / 9.0, 5.0 / 9.0},
	},
	{N: 4,
		Xi: []float64{
			math.Sqrt(3.0/7.0 - ((2.0 / 7.0) * math.Sqrt(6.0/5.0))),
			-math.Sqrt(3.0/7.0 - ((2.0 / 7.0) * math.Sqrt(6.0/5.0))),
			math.Sqrt(3.0/7.0 + ((2.0 / 7.0) * math.Sqrt(6.0/5.0))),
			-math.Sqrt(3.0/7.0 + ((2.0 / 7.0) * math.Sqrt(6.0/5.0)))},
		Wi: []float64{
			(18.0 + math.Sqrt(30)) / 36.0,
			(18.0 + math.Sqrt(30)) / 36.0,
			(18.0 - math.Sqrt(30)) / 36.0,
			(18.0 - math.Sqrt(30)) / 36.0}},
	{N: 5,
		Xi: []float64{0,
			(1.0 / 3.0) * math.Sqrt(5.0-(2*math.Sqrt(10.0/7.0))),
			-(1.0 / 3.0) * math.Sqrt(5.0-(2*math.Sqrt(10.0/7.0))),
			(1.0 / 3.0) * math.Sqrt(5.0+(2*math.Sqrt(10.0/7.0))),
			-(1.0 / 3.0) * math.Sqrt(5.0+(2*math.Sqrt(10.0/7.0))),
		},
		Wi: []float64{
			128.0 / 225.0,
			(322 + 13.0*math.Sqrt(70)) / 900.0,
			(322 + 13.0*math.Sqrt(70)) / 900.0,
			(322 - 13.0*math.Sqrt(70)) / 900.0,
			(322 - 13.0*math.Sqrt(70)) / 900.0,
		},
	},
}

func (i Integrator) GaussLegendre(a, b float64, n int) float64 {
	if n >= len(LegendrePoints) {
		return 0
	}

	bigF := func(x float64) float64 {
		prefix := (b - a) / 2
		postfix := (b + a) / 2
		return prefix * i.F(prefix*x+postfix)
	}

	points := LegendrePoints[n]
	out := 0.0
	for j, xi := range points.Xi {
		out += points.Wi[j] * bigF(xi)
	}

	return out
}
