package main

import (
	"fmt"
	"math"
	"os"

	"github.com/c0nrad/csim"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func main() {
	rows := 1
	cols := 4

	plots := make([][]*plot.Plot, rows)
	for j := 0; j < rows; j++ {
		plots[j] = make([]*plot.Plot, cols)
	}

	plots[0][0] = PlotDifferentiators(.4)
	plots[0][1] = PlotDifferentiators(.2)
	plots[0][2] = PlotDifferentiators(.1)
	plots[0][3] = PlotDifferentiators(.05)

	img := vgimg.New(vg.Points(1200), vg.Points(300))
	dc := draw.New(img)

	t := draw.Tiles{
		Rows:      rows,
		Cols:      cols,
		PadX:      vg.Millimeter,
		PadY:      vg.Millimeter,
		PadTop:    vg.Points(2),
		PadBottom: vg.Points(2),
		PadLeft:   vg.Points(5),
		PadRight:  vg.Points(5),
	}

	canvases := plot.Align(plots, t, dc)
	for j := 0; j < rows; j++ {
		for i := 0; i < cols; i++ {
			if plots[j][i] != nil {
				plots[j][i].Draw(canvases[j][i])
			}
		}
	}

	w, err := os.Create("multiple_h.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}

}

func Sample(f csim.Func, start, stop, h float64) plotter.XYs {
	out := []plotter.XY{}
	for x := start; x < stop; x += h {
		out = append(out, plotter.XY{X: x, Y: f(x)})
	}
	return out
}

func PlotDifferentiators(H float64) *plot.Plot {
	f := math.Sin
	nDiff := csim.NumericalDifferentiator{H: H, F: f}
	xStart := 0.0 + math.Pi/2
	xStop := 2*math.Pi + math.Pi/2

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = fmt.Sprintf("\t\tNumerical Differentiation, H=%.03f, f=Sin(x)", H)
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	dfdt := plotter.NewFunction(math.Cos)
	dfdt.Color = plotutil.Color(0)
	dfdtF, _ := plotter.NewScatter(Sample(nDiff.ForwardDifferenceDerivative(), xStart, xStop, H))
	dfdtF.Color = plotutil.Color(1)
	dfdtF.Shape = draw.PyramidGlyph{}

	dfdtB, _ := plotter.NewScatter(Sample(nDiff.BackwardDifferenceDerivative(), xStart, xStop, H))
	dfdtB.Color = plotutil.Color(2)
	dfdtB.Shape = draw.BoxGlyph{}

	dfdtC, _ := plotter.NewScatter(Sample(nDiff.CentralDifferenceDerivative(), xStart, xStop, H))
	dfdtC.Color = plotutil.Color(3)
	dfdtC.Shape = draw.CircleGlyph{}

	// Add the functions and their legend entries.
	p.Add(dfdt, dfdtF, dfdtB, dfdtC)
	p.Legend.Add("d/dx Sin(x)", dfdt)
	p.Legend.Add("Forward df/dx", dfdtF)
	p.Legend.Add("Backward df/dx", dfdtB)
	p.Legend.Add("Central df/dx", dfdtC)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	p.X.Min = xStart
	p.X.Max = xStop
	p.Y.Min = -1.2
	p.Y.Max = 1.2

	return p
}
