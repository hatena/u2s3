package util

import (
	"bytes"
	"net/http"
	"os"
	"regexp"
	"text/template"

	"github.com/hatena/u2s3/pkg/config"
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

func GenerateUploadKey(keyTemp *config.UploadKeyTemplate, keyFmt string) (string, error) {
	host, err := os.Hostname()
	if err != nil {
		return "", err
	}
	temp := template.New("key")
	template.Must(temp.Parse(keyFmt))
	var res bytes.Buffer
	err = temp.Execute(&res, map[string]interface{}{
		"Output":   keyTemp.Output,
		"Year":     keyTemp.Year,
		"Month":    keyTemp.Month,
		"Day":      keyTemp.Day,
		"Hour":     keyTemp.Hour,
		"Minute":   keyTemp.Minute,
		"Second":   keyTemp.Second,
		"Hostname": host,
		"Seq":      keyTemp.Seq,
	})
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
