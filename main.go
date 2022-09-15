package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/deelawn/urbit-gob/co"
	"github.com/manifoldco/promptui"
)

// available strategies
type Strategy int

const (
	All Strategy = iota
	AnyApprox
	AnyEnglish
	OnlyApprox
	OnlyEnglish
	Doubles
)

type StrategyPrompt struct {
	name     string
	strategy Strategy
}

var strategyPrompts = []StrategyPrompt{
	{
		"All", All,
	},
	{
		"Any English or Slang", AnyApprox,
	},
	{
		"Any English", AnyEnglish,
	},
	{
		"Only English or Slang", OnlyApprox,
	},
	{
		"Only English", OnlyEnglish,
	},
	{
		"Doubles", Doubles,
	},
}

// generates an English word map
func generateWords(fileName string) map[string]bool {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	var words = make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return words
}

// loaded words globals
var singleEnglishWords = generateWords("./wordlists/name/english-single.txt")
var doubleEnglishWords = generateWords("./wordlists/name/english-double.txt")
var singleApproxWords = generateWords("./wordlists/name/approx-single.txt")
var doubleApproxWords = generateWords("./wordlists/name/approx-double.txt")

// strategy filters
func matchApprox(phoneme string) bool {
	return singleApproxWords[phoneme] || doubleApproxWords[phoneme]
}

func matchEnglish(phoneme string) bool {
	return singleEnglishWords[phoneme] || doubleEnglishWords[phoneme]
}

func anyApprox(planet string) bool {
	deSigged := strings.Replace(planet, "~", "", 1)
	parts := strings.Split(deSigged, "-")
	for i := 0; i < len(parts); i++ {
		part := parts[i]

		if matchApprox(part) {
			return true
		}
	}

	return false
}

func anyEnglish(planet string) bool {
	deSigged := strings.Replace(planet, "~", "", 1)
	parts := strings.Split(deSigged, "-")
	for i := 0; i < len(parts); i++ {
		part := parts[i]

		if matchEnglish(part) {
			return true
		}
	}

	return false
}

func onlyApprox(planet string) bool {
	deSigged := strings.Replace(planet, "~", "", 1)
	parts := strings.Split(deSigged, "-")

	return (matchApprox(parts[0]) || matchEnglish(parts[0])) && (matchApprox(parts[1]) || matchEnglish(parts[1]))
}

func onlyEnglish(planet string) bool {
	deSigged := strings.Replace(planet, "~", "", 1)
	parts := strings.Split(deSigged, "-")

	return matchEnglish(parts[0]) && matchEnglish(parts[1])
}

func doubles(planet string) bool {
	deSigged := strings.Replace(planet, "~", "", 1)
	parts := strings.Split(deSigged, "-")

	return parts[0] == parts[1]
}

// generates list of planets under parent
func planets(parent string, strategy Strategy) []string {
	hex, err := co.Patq2Hex(parent)
	if err != nil {
		fmt.Println(err)
	}

	var planets = make([]string, 0)

	for i := 1; i <= 0xFFFF; i++ {
		base := strconv.FormatInt(int64(i), 16)
		s := fmt.Sprintf("%04s", base) + hex
		p, err := co.Hex2Patp(s)
		if err != nil {
			fmt.Println(err)
		}

		switch strategy {
			case AnyApprox:
				if anyApprox(p) {
					planets = append(planets, p)
					// doesn't check english, so only continue if match
					continue
				}

			case AnyEnglish:
				if anyEnglish(p) {
					planets = append(planets, p)
				}
				continue

			case OnlyApprox:
				if onlyApprox(p) {
					planets = append(planets, p)
				}
				continue

			case OnlyEnglish:
				if onlyEnglish(p) {
					planets = append(planets, p)
				}
				continue

			case Doubles:
				if doubles(p) {
					planets = append(planets, p)
				}
				continue

			default:
				planets = append(planets, p)
		}
	}

	sort.Slice(planets, func(i, j int) bool {
		return planets[i] < planets[j]
	})

	return planets
}

// input validation
func validate(input string) error {
	// ensure valid patp
	isValid := co.IsValidPat(input)
	if !isValid {
		return errors.New("Invalid patp")
	}

	// ensure valid size
	size, err := co.Clan(input)
	if size != "star" || err != nil {
		return errors.New("Must be a star")
	}

	return nil
}

func main() {
	// prompt for parent
	parentPrompt := promptui.Prompt{
		Label:    "Which star? (e.g., ~marzod)",
		Validate: validate,
	}

	parent, err := parentPrompt.Run()

	if err != nil {
		fmt.Printf("bad input: %v\n", err)
		return
	}

	// prompt for strategy
	strategyPrompt := promptui.Select{
		Label: "Which set of planets?",
		Items: strategyPrompts,
	}

	k, _, err := strategyPrompt.Run()

	if err != nil {
		fmt.Printf("selection failed: %v\n", err)
		return
	}

	// find planets
	results := planets(parent, strategyPrompts[k].strategy)

	// write output
	pbytes := strings.Join(results, "\n")
	ioutil.WriteFile(fmt.Sprintf("./%s_planets.txt", parent), []byte(pbytes), 0600)

	fmt.Println("Done :)")
}
