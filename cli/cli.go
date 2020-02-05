package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var globalFlags = []cli.Flag{}

var commonFlags = []cli.Flag{
	cli.StringFlag{Name: "f, file", Usage: "input file"},
	cli.StringFlag{Name: "o, output", Usage: "directory for output"},
	cli.StringFlag{Name: "kf, key-format", Value: "{{.Output}}/{{.Year}}/{{.Month}}/{{.Day}}/{{.Hostname}}-{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.Minute}}_{{.Seq}}.log.gz", Usage: "object key format"},
	cli.StringFlag{Name: "b, bucket", Usage: "bucket name"},
	cli.IntFlag{Name: "m, max-retry", Value: 5, Usage: "limit retry times"},
	cli.IntFlag{Name: "cpu", Usage: "cpu usage limitation (%)"},
	cli.IntFlag{Name: "memory", Usage: "memory usage limitation (MB)"},
	cli.IntFlag{Name: "rate", Usage: "bandwidth rate limit (MB)"},
	cli.StringFlag{Name: "dev", Value: "eth0", Usage: "rate limit device (default: eth0)"},
}

var commands = []cli.Command{
	{
		Name:   "upload-log",
		Usage:  "Upload log files to S3",
		Action: uploadLogCmd,
		Flags: append(commonFlags,
			cli.IntFlag{Name: "s, step", Value: 30, Usage: "step for splitting logs. (min)"},
			cli.StringFlag{Name: "l, log-format", Value: "tsv", Usage: "tsv or ssv (default: tsv)"},
		),
	},
	{
		Name:   "upload-file",
		Usage:  "Upload files to S3",
		Action: uploadFileCmd,
		Flags: append(commonFlags,
			cli.StringFlag{Name: "ff, filename-format", Usage: "file name format e.g. " + `(?<Year>\d{4})-(?<Month>\d{2})-(?<Day>\d{2})`},
		),
	},
	{
		Name:   "sync-file",
		Usage:  "Sync files with S3",
		Action: syncFileCmd,
		Flags: append(commonFlags,
			cli.StringFlag{Name: "ff, filename-format", Usage: "file name format e.g. " + `(?<Year>\d{4})-(?<Month>\d{2})-(?<Day>\d{2})`},
		),
	},
}

func commandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

var version = "unset"

func Main() {
	app := cli.NewApp()
	app.Name = "u2s3"
	app.Version = version

	app.Flags = globalFlags
	app.Commands = commands
	app.CommandNotFound = commandNotFound

	app.Run(os.Args)
}
