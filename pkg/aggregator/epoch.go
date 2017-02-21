package aggregator

import (
	"compress/gzip"
	"os"
)

type Epoch struct {
	fp     *os.File
	writer *gzip.Writer
}

type EpochManager struct {
}

func NewEpoch(epochKey, keyFmt string) *Epoch {

	return &Epoch{}
}

func NewEpochManager() *EpochManager {
	return &EpochManager{}
}

func (m *EpochManager) HasEpoch(key string) bool {
	return false
}

func (m *EpochManager) GetEpoch(key string) *Epoch {
	return nil
}
