package core

type Aggregator interface {
	Run() error
	GetUploadableFiles() []UploadableFile
	Close()
}
