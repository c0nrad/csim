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

type FuncIntercept struct {
	In map[float64]int
}

func NewIntercept() FuncIntercept {
	return FuncIntercept{In: make(map[float64]int)}
}

func (c FuncIntercept) F(f csim.Func) csim.Func {
	return func(x float64) float64 {
		c.In[x]++
		return f(x)
	}
}

func (c FuncIntercept) Unique() int {
	return len(c.In)
}

func main() {

	rows := 2
	cols := 2

	plots := make([][]*plot.Plot, rows)
	for j := 0; j < rows; j++ {
		plots[j] = make([]*plot.Plot, cols)
	}

	f1 := func(x float64) float64 {
		return 1.0 / (x + 2.0)
	}
	soln := math.Log(3) - math.Log(1)
	p1 := PlotErrors(f1, soln, -1, 1)
	p1.Title.Text = "Integral[1/(1+x), -1, 1]"
	// panic("done")

	// 2
	f2 := func(x float64) float64 {
		return math.Sin(x * x)
	}
	soln2 := 0.310268301723381101808152423165396507574509388832446717732
	p2 := PlotErrors(f2, soln2, 0, 1)
	p2.Title.Text = "Integral[sin(x^2), 0, 1]"

	// 3
	f3 := func(x float64) float64 {
		return 2*x*x*x*x + 3*x*x*x + 4*x*x + 5*x + 6
	}
	soln3 := 82787.0
	p3 := PlotErrors(f3, soln3, -10, 10)
	p3.Title.Text = "Integral[x^5 + 2x^4 + 3x^3 + 4x^2 + 5x + 6, -10, 10]"
	// panic("here")

	// 4
	// f4 := func(x float64) float64 {
	// 	return math.Exp(math.Cos(x))
	// }
	// soln4 := 5.1834282608
	// p4 := PlotErrors(f4, soln4, 0, 5)
	// p4.Title.Text = "Integral[e^cos(x), 0, 5]"
	f4 := func(x float64) float64 {
		return math.Exp(-2 * x)
	}
	soln4 := .49084
	p4 := PlotErrors(f4, soln4, 0, 2)
	p4.Title.Text = "Integral[e^-2x, 0, 2]"

	plots[0][0] = p1
	plots[1][0] = p2
	plots[0][1] = p3
	plots[1][1] = p4

	img := vgimg.New(vg.Points(800), vg.Points(800))
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

	w, err := os.Create("integrate_example.png")
	if err != nil {
		panic(err)
	}
	defer w.Close()
	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}

}

func PlotErrors(f1 csim.Func, soln float64, startX, endX float64) *plot.Plot {
	// Ns := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}
	Ns := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 15, 20, 25, 30, 40, 50, 60, 70, 80, 90, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	p.Title.Text = "Change me"
	p.X.Label.Text = "N"
	p.Y.Label.Text = "Error"
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = plot.LogTicks{}
	p.X.Min = float64(Ns[0] - 1)
	p.X.Max = float64(Ns[len(Ns)-1] + 1)

	p.Y.Scale = plot.LogScale{}
	p.Y.Tick.Marker = plot.LogTicks{}
	p.Y.Min = .000002
	p.Y.Max = 1

	rectScatter, _ := plotter.NewLine(plotter.XYs{})
	rectScatter.Color = plotutil.Color(0)
	// rectScatter.Shape = draw.CircleGlyph{}

	trapScatter, _ := plotter.NewLine(plotter.XYs{})
	trapScatter.Color = plotutil.Color(1)

	simpScatter, _ := plotter.NewLine(plotter.XYs{})
	simpScatter.Color = plotutil.Color(2)

	gaussScatter, _ := plotter.NewLine(plotter.XYs{})
	gaussScatter.Color = plotutil.Color(3)

	// Add the functions and their legend entries.
	p.Add(rectScatter, trapScatter, simpScatter, gaussScatter)
	p.Legend.Add("Rectangular", rectScatter)
	p.Legend.Add("Trapezoidal", trapScatter)
	p.Legend.Add("Simpson", simpScatter)
	p.Legend.Add("Gauss-Legendre", gaussScatter)

	p.Legend.ThumbnailWidth = 0.5 * vg.Inch
	for _, N := range Ns {
		// H := (endX - startX) / float64(N)

		fmt.Printf("Soln=%.06f, N=%d\n", soln, N)
		fmt.Println("------")
		fmt.Println("Method, Guess, Error")

		rIntercept := NewIntercept()
		rIntegrate := csim.Integrator{F: rIntercept.F(f1), N: N}
		rGuess := rIntegrate.Rectangular(startX, endX)
		fmt.Printf("Rect  %.06f %1.2e %d\n", rGuess, Error(rGuess, soln), rIntercept.Unique())
		if rGuess != 0 {
			rectScatter.XYs = append(rectScatter.XYs, plotter.XY{X: float64(N), Y: Error(rGuess, soln)})
		}

		tIntercept := NewIntercept()
		tIntegrate := csim.Integrator{F: tIntercept.F(f1), N: N}
		tGuess := tIntegrate.Trapezoidal(startX, endX)
		fmt.Printf("Trap  %.06f %1.2e %d\n", tGuess, Error(tGuess, soln), tIntercept.Unique())
		if tGuess != 0 {
			trapScatter.XYs = append(trapScatter.XYs, plotter.XY{X: float64(N), Y: Error(tGuess, soln)})
		}

		sIntercept := NewIntercept()
		sIntegrate := csim.Integrator{F: sIntercept.F(f1), N: N}
		sGuess := sIntegrate.Simpson(startX, endX)
		fmt.Printf("Simp  %.06f %1.2e %d\n", sGuess, Error(sGuess, soln), sIntercept.Unique())
		if sGuess != 0 {
			if Error(sGuess, soln) != 0 {
				simpScatter.XYs = append(simpScatter.XYs, plotter.XY{X: float64(N), Y: Error(sGuess, soln)})
			}
		}

		gIntercept := NewIntercept()
		gIntegrate := csim.Integrator{F: gIntercept.F(f1), N: N}
		gGuess := gIntegrate.GaussLegendre(startX, endX, N)
		fmt.Printf("Gauss %.06f %1.2e %d\n", gGuess, Error(gGuess, soln), gIntercept.Unique())
		if gGuess != 0 {
			gaussScatter.XYs = append(gaussScatter.XYs, plotter.XY{X: float64(N), Y: Error(gGuess, soln)})
		}
	}

	return p

}

func Error(guess, correct float64) float64 {
	return math.Abs((guess - correct) / correct)
}
