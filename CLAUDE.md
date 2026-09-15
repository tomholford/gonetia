# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Gonetia is a Go CLI tool that generates lists of valid Urbit planet names issuable from a given star. It filters planets through wordlist-matching strategies (English, approximate/slang, doubles, alliteration) and writes results to `./output/[star]/`.

## Commands

Quality scripts live under `ops/` (the built binary is `bin/gonetia`, gitignored):

| Script | Purpose |
| --- | --- |
| `./ops/fmt.sh` | `golangci-lint fmt` (gofmt + gci) |
| `./ops/lint.sh` | `golangci-lint run` |
| `./ops/build.sh` | `go build -o bin/gonetia .` |
| `./ops/test.sh` | `go test -count=1 ./...` |
| `./ops/gauntlet.sh` | **fmt → lint → build → test** (fail-fast) |
| `./ops/update-wordlist.sh` | pull the `wordlists/` submodule |

Requires [golangci-lint](https://golangci-lint.run/) on `PATH`.

### Agent loop (required)

After **every** code change, run:

```sh
./ops/gauntlet.sh
```

Do not stop at “it seems fine.” Iterate until the gauntlet exits 0. Fix format/lint failures in the same change set.

```sh
# Run with argument
./bin/gonetia "~marzod"

# Run interactive (prompts for star)
./bin/gonetia
```

## Setup

```bash
git submodule update --init --recursive
./ops/build.sh
```

Submodules are required at compile time: the four name wordlists are `go:embed`ed into the binary.

## Architecture

Single-file app (`main.go`). All logic lives here:

- **Strategies** (enum-like consts): `All`, `AnyApprox`, `AnyEnglish`, `OnlyApprox`, `OnlyEnglish`, `Doubles`, `Alliteration` — define how planets are filtered
- **Wordlist loading** (`loadWordlists` / `generateWords`): reads embedded `wordlists/name/*.txt` into four global maps (single/double × english/approx)
- **Planet generation** (`makePlanets`): uses `urbit-gob` to convert star hex → all 65,535 child planet patps
- **Filtering** (`filterPlanets`): applies strategy-specific match functions to planet lists
- **Output** (`writeResults`): writes filtered lists to `./output/[star]/[strategy]_planets.txt`
- **Validation** (`validate`): ensures input is a valid Urbit star using `urbit-gob`

## Key Dependencies

- `github.com/deelawn/urbit-gob` — Urbit patp/hex conversion and clan detection
- `github.com/manifoldco/promptui` — interactive CLI prompts
- `wordlists/` — git submodule pointing to `ashelkovnykov/urbit-wordlists` (compile-time embed source)
