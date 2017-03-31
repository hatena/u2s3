package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var globalFlags = []cli.Flag{}

var commands = []cli.Command{
	{
		Name:   "upload",
		Usage:  "Upload a log file to S3",
		Action: uploadCmd,
		Flags: []cli.Flag{
			cli.StringFlag{Name: "f, file", Usage: "input file"},
			cli.IntFlag{Name: "s, step", Value: 30, Usage: "step for splitting logs. (min)"},
			cli.StringFlag{Name: "p, prefix", Usage: "prefix of splitted logs"},
			cli.StringFlag{Name: "o, output", Usage: "directory for output"},
			cli.StringFlag{Name: "l, log-format", Value: "tsv", Usage: "tsv or ssv (default: tsv)"},
			cli.StringFlag{Name: "k, key", Value: "{{.Output}}/{{.Year}}/{{.Month}}/{{.Day}}/{{.Hostname}}-{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.Minute}}_{{.Seq}}.log.gz", Usage: "object key format"},
			cli.StringFlag{Name: "b, bucket", Usage: "bucket name"},
			cli.IntFlag{Name: "m, max-retry", Value: 5, Usage: "the maximum number of retry times"},
			cli.IntFlag{Name: "cpu", Usage: "cpu usage limitation (%)"},
			cli.IntFlag{Name: "memory", Usage: "memory usage limitation (MB)"},
			cli.IntFlag{Name: "rate", Usage: "bandwidth rate limit (MB)"},
			cli.StringFlag{Name: "dev", Value: "eth0", Usage: "rate limit device (default: eth0)"},
			cli.BoolFlag{Name: "ca, content-aware", Usage: "Content-aware uploading is enabled"},
		},
	},
}

func commandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

func Main() {
	app := cli.NewApp()
	app.Name = "log2s3"
	app.Version = "0.1.0"
	app.Author = "taku-k"
	app.Email = "taakuu19@gmail.com"

	app.Flags = globalFlags
	app.Commands = commands
	app.CommandNotFound = commandNotFound

	app.Run(os.Args)
}
