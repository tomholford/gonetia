# Contributing

## Requirements

- [Go](https://go.dev/dl/) 1.26+ (see `go.mod`)
- [golangci-lint](https://golangci-lint.run/) on `PATH`

## Quality gate

One fail-fast sequence (fmt → lint → build → test):

```sh
./ops/gauntlet.sh
```

Run it after every code change. CI runs the same sequence on pull requests and `master`. Individual steps:

| Script | Purpose |
| --- | --- |
| `./ops/fmt.sh` | format (`golangci-lint fmt`) |
| `./ops/lint.sh` | `golangci-lint run` |
| `./ops/build.sh` | build → `bin/gonetia` (gitignored) |
| `./ops/test.sh` | `go test -count=1 ./...` |
| `./ops/update-wordlist.sh` | pull the `wordlists/` submodule |

Single test:

```sh
go test -count=1 -run TestName .
```

## Layout

```
main.go           CLI + filters
main_test.go      unit tests
wordlists/        git submodule (embedded at compile time)
ops/              quality scripts (tracked)
bin/              build output (gitignored)
```

Name wordlists are `go:embed`ed from `wordlists/name/*.txt`. After updating the submodule, rebuild so the new lists are compiled in.

## Release

Releases are cut by tag. GoReleaser (`.goreleaser.yaml`, `.github/workflows/release.yml`) cross-compiles, publishes GitHub Release archives + checksums, and pushes a Homebrew cask to `tomholford/homebrew-tap`.

```sh
git tag v0.1.0
git push origin v0.1.0
```

PRs that touch `.goreleaser.yaml` or the release workflow run `goreleaser release --snapshot --skip=publish`. Locally:

```sh
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```

The tap lives in a separate repo. The default `GITHUB_TOKEN` cannot write it; set repo secret `HOMEBREW_TAP_TOKEN` to a fine-grained PAT with Contents: Read and write on `tomholford/homebrew-tap`. Prerelease tags (`v0.1.0-rc.1`) skip the tap upload (`skip_upload: auto`).
