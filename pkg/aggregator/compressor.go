package aggregator

type Compressor struct {
}

func NewCompressor() *Compressor {
	return &Compressor{}
}

func (c *Compressor) HasEpochFile(epoch string) bool {
	return false
}

func (c *Compressor) AddEpoch(e *Epoch) {

}

func (c *Compressor) Compress(key, l string) {

}
