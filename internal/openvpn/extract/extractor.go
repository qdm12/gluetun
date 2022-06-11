package extract

type Extractor struct{}

func New() *Extractor {
	return new(Extractor)
}
