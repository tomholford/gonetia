package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/deelawn/urbit-gob/co"
	"github.com/manifoldco/promptui"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
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
func makePlanets(parent string) []string {
	fmt.Sprintf("Making planet list for %s ...\n", parent)

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

		planets = append(planets, p)
	}

	sort.Slice(planets, func(i, j int) bool {
		return planets[i] < planets[j]
	})

	return planets
}

func filterPlanets(planets []string, strategy Strategy) []string {
	var output = make([]string, 0)

	for i := 1; i < len(planets); i++ {
		p := planets[i]

		switch strategy {
		case AnyApprox:
			if anyApprox(p) {
				output = append(output, p)
				// doesn't check english, so only continue if match
				continue
			}

		case AnyEnglish:
			if anyEnglish(p) {
				output = append(output, p)
			}
			continue

		case OnlyApprox:
			if onlyApprox(p) {
				output = append(output, p)
			}
			continue

		case OnlyEnglish:
			if onlyEnglish(p) {
				output = append(output, p)
			}
			continue

		case Doubles:
			if doubles(p) {
				output = append(output, p)
			}
			continue

		default:
			output = append(output, p)
		}
	}

	return output
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

func writeResults(parent string, strategy string, results []string) error {
	fmt.Sprintf("Writing output for %s ...\n", strategy)

	pbytes := strings.Join(results, "\n")
	err := ioutil.WriteFile(fmt.Sprintf("./output/%s/%s_planets.txt", parent, strategy), []byte(pbytes), 0755)

	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func main() {
	var parent string

	// get input from args...
	argLength := len(os.Args[1:])
	if argLength > 0 {
		arg := os.Args[1]
		aErr := validate(arg)
		if aErr == nil {
			parent = arg
		}
	}

	// ... or prompt for parent
	if parent == "" {
		parentPrompt := promptui.Prompt{
			Label:    "Which star? (e.g., ~marzod)",
			Validate: validate,
		}

		var pErr error
		parent, pErr = parentPrompt.Run()

		if pErr != nil {
			fmt.Printf("bad input: %v\n", pErr)
			return
		}
	}

	// find planets
	planets := makePlanets(parent)

	// filter for each strategy
	anyApproxPlanets := filterPlanets(planets, AnyApprox)
	onlyApproxPlanets := filterPlanets(planets, OnlyApprox)
	anyEnglishPlanets := filterPlanets(planets, AnyEnglish)
	onlyEnglishPlanets := filterPlanets(planets, OnlyEnglish)
	doublesPlanets := filterPlanets(planets, Doubles)

	// write output
	_, oErr := os.Stat("output")
	if os.IsNotExist(oErr) {
		os.Mkdir("./output", 0755)
	}
	deSiggedParent := strings.Replace(parent, "~", "", 1)
	e := os.Mkdir(fmt.Sprintf("./output/%s", deSiggedParent), 0755)
	if e != nil {
		fmt.Println(e)
	}

	writeResults(deSiggedParent, "any_approx", anyApproxPlanets)
	writeResults(deSiggedParent, "only_approx", onlyApproxPlanets)
	writeResults(deSiggedParent, "any_english", anyEnglishPlanets)
	writeResults(deSiggedParent, "only_english", onlyEnglishPlanets)
	writeResults(deSiggedParent, "doubles", doublesPlanets)

	fmt.Println("Done :)")
}
