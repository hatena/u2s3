# u2s3

CPUやメモリ，ネットワーク帯域のリソース制限を行いつつAmazon S3にログファイルなどをアップロードできます．


## Usage

* Upload log files which have content-awareness

```
$ u2s3 upload-log \
       -f access_log.tsv \
       -b test_bucket \
       -s 30 \
       -kf "{{.Year}}/{{.Month}}/{{.Day}}/{{.Hostname}}-{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.Minute}}_{{.Seq}}.log.gz"
&config.UploadConfig{
  FileName:        "access_log.tsv",
  LogFormat:       "tsv",
  KeyFormat:       "{{.Year}}/{{.Month}}/{{.Day}}/{{.Hostname}}-{{.Year}}{{.Month}}{{.Day}}{{.Hour}}{{.Minute}}_{{.Seq}}.log.gz",
  OutputPrefixKey: "",
  Step:            30,
  Bucket:          "test_bucket",
  MaxRetry:        5,
  CPULimit:        0,
  MemoryLimit:     0,
  RateLimit:       0,
  Device:          "eth0",
  FilenameFormat:  "",
}
2017/04/06 06:25:01 [info] No limit resources
2017/04/06 06:25:01 [info] Uploaded 2017/02/24/ubuntu-xenial-201702241000_1.log.gz
2017/04/06 06:25:01 [info] Uploaded 2017/02/24/ubuntu-xenial-201702241030_1.log.gz
```

* Upload binaries

```
$ ls 2017*.tsv
20170331.tsv  20170406.tsv
$ u2s3 upload-file \
       -f "2017*.tsv" \
       -b test_bucket \
       -ff "(?P<Year>\d{4})(?P<Month>\d{2})(?P<Day>\d{2}).tsv" \
       -kf "{{.Year}}/{{.Month}}/{{.Day}}/{{.Hostname}}-{{.Year}}{{.Month}}{{.Day}}_{{.Seq}}.log.gz"
&config.UploadConfig{
  FileName:        "2017*.tsv",
  LogFormat:       "",
  KeyFormat:       "{{.Year}}/{{.Month}}/{{.Day}}/{{.Hostname}}-{{.Year}}{{.Month}}{{.Day}}_{{.Seq}}.log.gz",
  OutputPrefixKey: "",
  Step:            0,
  Bucket:          "test_bucket",
  MaxRetry:        5,
  CPULimit:        0,
  MemoryLimit:     0,
  RateLimit:       0,
  Device:          "eth0",
  FilenameFormat:  "(?P<Year>\\d{4})(?P<Month>\\d{2})(?P<Day>\\d{2}).tsv",
}
2017/04/06 06:25:01 [info] No limit resources
2017/04/06 06:25:16 [info] Uploaded 2017/03/31/ubuntu-xenial-20170331_1.log.gz
2017/04/06 06:25:21 [info] Uploaded 2017/04/06/ubuntu-xenial-20170406_1.log.gz
```

## Install

To install, use `go get`:

```bash
$ go get github.com/hatena/u2s3
```

## Contribution

1. Fork ([https://github.com/hatena/u2s3/fork](https://github.com/hatena/u2s3/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[taku-k](https://github.com/taku-k)
[itchyny](https://github.com/itchyny)
