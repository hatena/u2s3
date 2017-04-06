package core

type Aggregator interface {
	Run() error
	GetUploadableFiles() []UploadableFile
	GenFetchJobs() chan *fetchJob
	Close()
}
