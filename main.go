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
	AnyEnglish
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
		"Any English", AnyEnglish,
	},
	{
		"Only English", OnlyEnglish,
	},
	{
		"Doubles", Doubles,
	},
}

// generates an English word map
func generateWords() map[string]bool {
	file, err := os.Open("words.txt")
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

// loaded words global
var words = generateWords()

// strategy filters
func all(planet string) bool {
	return true
}

func anyEnglish(planet string) bool {
	deSigged := strings.Replace(planet, "~", "", 1)
	parts := strings.Split(deSigged, "-")
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		exists := words[part]
		if exists {
			return true
		}
	}

	return false
}

func onlyEnglish(planet string) bool {
	deSigged := strings.Replace(planet, "~", "", 1)
	parts := strings.Split(deSigged, "-")

	return words[parts[0]] && words[parts[1]]
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
		case All:
			planets = append(planets, p)
			continue

		case Doubles:
			if doubles(p) {
				planets = append(planets, p)
			}
			continue

		case AnyEnglish:
			if anyEnglish(p) {
				planets = append(planets, p)
			}
			continue

		case OnlyEnglish:
			if onlyEnglish(p) {
				planets = append(planets, p)
			}
			continue

		default:
			continue
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
	ioutil.WriteFile("planets.txt", []byte(pbytes), 0600)

	fmt.Println("Done :)")
}
