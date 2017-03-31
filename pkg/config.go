package pkg

type UploadConfig struct {
	FileName        string
	LogFormat       string
	KeyFormat       string
	OutputPrefixKey string
	Step            int
	Bucket          string
	Gzipped         bool
	MaxRetry        int
	CPULimit        int
	MemoryLimit     int
	RateLimit       int
	Device          string
	ContentAware    bool
}
