# gonetia

Gonetia is a simple command-line utility for generating a list of planet names
issuable from an Urbit star. It is inspired by
[Venetia](https://github.com/tylershuster/venetia), but written in Go instead
of JS.

## Preview

https://user-images.githubusercontent.com/16504501/190836787-212afe87-b2f9-4352-a552-0cadff6b1338.mp4

## Setup

### Prerequisites

This project requires Go 1.18+. One option is to use [gvm](https://github.com/moovweb/gvm). Alternatively, on MacOS it is simple to install Go with [homebrew](https://brew.sh/).

Once `brew` is installed, install Go like so:
```
brew install go
```
### Instructions

1. Clone this repo
2. `git submodule update --init --recursive`
3. `go build`

## Usage

Enter a star in patp format (e.g., `~marzod`) either as a command line argument, or in response to the interactive prompt.

### Command Line Argument

```
./gonetia "~marzod"
```

### Interactive Mode

```sh
./gonetia

# Enter a star in patp format
>Which star? (e.g., ~marzod):

~marzod
```

### Output

A [planet](https://developers.urbit.org/reference/glossary/planet) identity is a four-syllable name composed of two six-character segments, such as `~sampel-palnet`. 

The script uses [urbit-wordlists](https://github.com/ashelkovnykov/urbit-wordlists) and  various strategies to generate output:

- *AnyEnglish*: at least one segment matches `wordlists/name/english-single` or `wordlists/name/english-double`.
- *OnlyEnglish*: both segments match `wordlists/name/english-single` or `wordlists/name/english-double`.
- *AnyApprox*:  at least one segment matches `wordlists/name/approx-single` or `wordlists/name/approx-double`.
- *OnlyApprox*:  both segments match any of the wordlists.
- *Doubles*: both segments are identical (e.g., ~datnut-datnut)


Output for each is written to `./output/[star]/[strategy]_planets.txt`, for 5 total files per run.

## Special Thanks

- Thanks to @tylershuster for creating [Venetia](https://github.com/tylershuster/venetia)
- Thanks to @deelawn for creating [urbit-gob](https://github.com/deelawn/urbit-gob)
- Thanks to @ashelkovnykov for creating [urbit-wordlists](https://github.com/ashelkovnykov/urbit-wordlists)
- Thanks to @michelleylai for feature ideas, feedback, testing, and proofreading docs
