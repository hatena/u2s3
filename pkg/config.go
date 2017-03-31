package pkg

type UploadConfig struct {
	FileName        string
	LogFormat       string
	KeyFormat       string
	OutputPrefixKey string
	Step            int
	Bucket          string
	MaxRetry        int
	CPULimit        int
	MemoryLimit     int
	RateLimit       int
	Device          string
	ContentAware    bool
	FilenameFormat  string
}

type UploadKeyTemplate struct {
	Year   string
	Month  string
	Day    string
	Hour   string
	Minute string
	Second string
	Output string
	Seq    int
}
