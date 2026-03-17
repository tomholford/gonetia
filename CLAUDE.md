# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Gonetia is a Go CLI tool that generates lists of valid Urbit planet names issuable from a given star. It filters planets through wordlist-matching strategies (English, approximate/slang, doubles) and writes results to `./output/[star]/`.

## Build & Run

```bash
# Build
go build

# Run with argument
./gonetia "~marzod"

# Run interactive (prompts for star)
./gonetia

# Update wordlists submodule
bin/update-wordlist
```

## Setup

```bash
git submodule update --init --recursive
go build
```

No tests, linter, or CI exist yet.

## Architecture

Single-file app (`main.go`, ~300 lines). All logic lives here:

- **Strategies** (enum-like consts): `All`, `AnyApprox`, `AnyEnglish`, `OnlyApprox`, `OnlyEnglish`, `Doubles` — define how planets are filtered
- **Wordlist loading** (`generateWords`): reads files from `wordlists/` git submodule into four global maps (single/double × english/approx)
- **Planet generation** (`makePlanets`): uses `urbit-gob` to convert star hex → all 65,535 child planet patps
- **Filtering** (`filterPlanets`): applies strategy-specific match functions to planet lists
- **Output** (`writeResults`): writes filtered lists to `./output/[star]/[strategy]_planets.txt`
- **Validation** (`validate`): ensures input is a valid Urbit star using `urbit-gob`

## Key Dependencies

- `github.com/deelawn/urbit-gob` — Urbit patp/hex conversion and clan detection
- `github.com/manifoldco/promptui` — interactive CLI prompts
- `wordlists/` — git submodule pointing to `ashelkovnykov/urbit-wordlists`
