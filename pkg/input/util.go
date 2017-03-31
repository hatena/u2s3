package input

import (
	"net/http"
	"os"
)

func IsGzipped(fp *os.File) bool {
	defer fp.Seek(0, 0)

	buff := make([]byte, 512)
	_, err := fp.Read(buff)
	if err != nil {
		return false
	}
	filetype := http.DetectContentType(buff)
	return filetype == "application/x-gzip"
}
