package main

import (
	"fmt"

	"github.com/c0nrad/csim"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var H = .01
var count = 5

func main() {

	xStart := -1.0
	xStop := 1.0

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Legendre Polynomials via Rodrigues' Formula"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	for i := 0; i < count+1; i++ {
		l := csim.LegendrePoly(i)
		lPlot, err := plotter.NewScatter(csim.Sample(l, xStart, xStop, H))
		if err != nil {
			panic(err)
		}
		lPlot.Color = plotutil.Color(i)
		lPlot.Radius = 1.5
		lPlot.Shape = draw.CircleGlyph{}
		p.Add(lPlot)
		p.Legend.Add(fmt.Sprintf("L%d(x)", i), lPlot)
	}

	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	p.X.Min = xStart
	p.X.Max = xStop

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "legendre.png"); err != nil {
		panic(err)
	}
}
