package csim

import (
	"gonum.org/v1/plot/plotter"
)

func Sample(f Func, start, stop, h float64) plotter.XYs {
	out := []plotter.XY{}
	for x := start; x < stop; x += h {
		out = append(out, plotter.XY{X: x, Y: f(x)})
	}
	return out
}
