package util

import (
	"net/http"
	"os"
	"regexp"
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

/**
 * Parses url with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func GetParams(regEx, url string) (paramsMap map[string]string) {

	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}
