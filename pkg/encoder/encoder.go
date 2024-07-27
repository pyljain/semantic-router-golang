package encoder

type Encoder interface {
	Encode(data []string) ([][]float64, error)
}
