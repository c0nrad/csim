package csim

type Func func(float64) float64

type NumericalDifferentiator struct {
	H float64

	F Func
}

func (d NumericalDifferentiator) ForwardDifferenceDerivative() Func {
	return func(x float64) float64 {
		return (d.F(x+d.H) - d.F(x)) / d.H
	}
}

func (d NumericalDifferentiator) BackwardDifferenceDerivative() Func {
	return func(x float64) float64 {
		return (d.F(x) - d.F(x-d.H)) / d.H
	}
}

func (d NumericalDifferentiator) CentralDifferenceDerivative() Func {
	return func(x float64) float64 {
		return (d.F(x+d.H) - d.F(x-d.H)) / (2 * d.H)
	}
}

func D(f Func) Func {
	return NumericalDifferentiator{H: .01, F: f}.CentralDifferenceDerivative()
}
