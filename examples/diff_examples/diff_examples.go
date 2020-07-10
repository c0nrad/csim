package main

import (
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

type Example struct {
	Title string
	F     csim.Func

	YMin, YMax float64
}

var H = .01

func main() {

	rows := 2
	cols := 2

	examples := []Example{
		{"f(x) = Sin(x)",
			func(x float64) float64 { return math.Sin(x) }, 0, 0},
		{"f(x) = cos(2x) + exp(-x**2/2)*sin(10x)",
			func(x float64) float64 { return math.Cos(2*x) + math.Exp(math.Pow(x, 2)/-2)*math.Sin(10*x) }, 0, 0},
		{"f(x) = 5x - 3",
			func(x float64) float64 { return 2*x + -3 }, -5.0, 10.0},
		{"f(x) = exp(x)",
			func(x float64) float64 { return math.Exp(x) }, 0, 0},
	}
	plots := make([][]*plot.Plot, rows)
	for j := 0; j < rows; j++ {
		plots[j] = make([]*plot.Plot, cols)
	}

	plots[0][0] = GeneratePlot(examples[0])
	plots[1][0] = GeneratePlot(examples[1])
	plots[0][1] = GeneratePlot(examples[2])
	plots[1][1] = GeneratePlot(examples[3])

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

	w, err := os.Create("diff_examples.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}
}

func GeneratePlot(e Example) *plot.Plot {

	nDiff := csim.NumericalDifferentiator{H: H, F: e.F}
	xStart := 0.0
	xStop := 5.0

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = e.Title
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	f := plotter.NewFunction(e.F)
	f.Color = plotutil.Color(0)

	dfdtC, _ := plotter.NewScatter(csim.Sample(nDiff.CentralDifferenceDerivative(), xStart, xStop, H))
	dfdtC.Color = plotutil.Color(2)
	dfdtC.Shape = draw.CircleGlyph{}

	// Add the functions and their legend entries.
	p.Add(dfdtC, f)
	p.Legend.Add("f(x)", f)
	p.Legend.Add("df/dx", dfdtC)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	p.X.Min = xStart
	p.X.Max = xStop

	if e.YMax != 0 {
		p.Y.Max = e.YMax
		p.Y.Min = e.YMin
	}

	return p
}
