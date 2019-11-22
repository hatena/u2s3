# How to Release

1. Create a release branch. (`git checkout -b v0.1.4`)
2. Bump version in `cli/cli.go` file.
3. Write changes in CHANGELOG.md. (`git log v0.1.3..`)
4. Create Pull-Request.
5. Merge Pull-Requst.
6. Update working directory. (`git checkout master && git fetch && git merge origin/master --ff`)
7. Create binaries. (`make cross-build`)
8. Create Release and attach binaries. https://github.com/hatena/u2s3/releases/new
