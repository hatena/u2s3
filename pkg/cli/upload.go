package cli

import (
	//"bufio"
	"fmt"
	"os"

	"github.com/taku-k/log2s3-go/pkg/io"
	"github.com/urfave/cli"
)

func uploadCmd(c *cli.Context) error {
	if f := c.String("file"); f != "" {
		file, err := os.OpenFile(f, os.O_RDONLY, 0644)
		defer file.Close()
		if err != nil {
			return err
		}
		//_ := bufio.NewReader(file)
	} else {
		reader := io.NewStdinReader()
		fmt.Println(reader.Readln())
	}
	return nil
}
