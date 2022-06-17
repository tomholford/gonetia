# gonetia

Gonetia is a simple command-line utility for generating a list of planet names
issuable from an Urbit star. It is inspired by
[Venetia](https://github.com/tylershuster/venetia), but written in Go instead
of JS.

## Usage

This project requires go 1.18+. One option is to use [gvm](https://github.com/moovweb/gvm).

1. Clone this repo
2. `git submodule update --init --recursive`
3. `go run main.go`
4. Enter a star in patp format (e.g., `~marzod`), and select a strategy

Output is written to `./planets.txt`.

## Special Thanks

- Thanks to @tylershuster for creating [Venetia](https://github.com/tylershuster/venetia)
- Thanks to @deelawn for creating [urbit-gob](https://github.com/deelawn/urbit-gob)
- Thanks to @ashelkovnykov for creating [urbit-wordlists](https://github.com/ashelkovnykov/urbit-wordlists)
