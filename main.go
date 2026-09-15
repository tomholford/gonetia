package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/deelawn/urbit-gob/co"
	"github.com/manifoldco/promptui"
)

// Strategy selects which wordlist filter to apply.
type Strategy int

const (
	All Strategy = iota
	AnyApprox
	AnyEnglish
	OnlyApprox
	OnlyEnglish
	Doubles
	Alliteration
)

func generateWords(fileName string) map[string]bool {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	words := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words[scanner.Text()] = true
	}

	err = scanner.Err()
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return words
}

// loaded words globals
var singleEnglishWords = generateWords("./wordlists/name/english-single.txt")
var doubleEnglishWords = generateWords("./wordlists/name/english-double.txt")
var singleApproxWords = generateWords("./wordlists/name/approx-single.txt")
var doubleApproxWords = generateWords("./wordlists/name/approx-double.txt")

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

func alliteration(planet string) bool {
	deSigged := strings.Replace(planet, "~", "", 1)
	parts := strings.Split(deSigged, "-")

	return parts[0][0] == parts[1][0]
}

func makePlanets(parent string) ([]string, error) {
	fmt.Printf("Making planet list for %s ...\n", parent)

	hex, err := co.Patq2Hex(parent)
	if err != nil {
		return nil, fmt.Errorf("convert star to hex: %w", err)
	}

	planets := make([]string, 0, 0xFFFF)

	for i := 1; i <= 0xFFFF; i++ {
		base := strconv.FormatInt(int64(i), 16)
		s := fmt.Sprintf("%04s", base) + hex
		p, err := co.Hex2Patp(s)
		if err != nil {
			return nil, fmt.Errorf("convert hex to patp: %w", err)
		}

		planets = append(planets, p)
	}

	sort.Slice(planets, func(i, j int) bool {
		return planets[i] < planets[j]
	})

	return planets, nil
}

func filterPlanets(planets []string, strategy Strategy) []string {
	output := make([]string, 0)

	for _, p := range planets {
		switch strategy {
		case AnyApprox:
			if anyApprox(p) {
				output = append(output, p)
			}
		case AnyEnglish:
			if anyEnglish(p) {
				output = append(output, p)
			}
		case OnlyApprox:
			if onlyApprox(p) {
				output = append(output, p)
			}
		case OnlyEnglish:
			if onlyEnglish(p) {
				output = append(output, p)
			}
		case Doubles:
			if doubles(p) {
				output = append(output, p)
			}
		case Alliteration:
			if alliteration(p) {
				output = append(output, p)
			}
		default:
			output = append(output, p)
		}
	}

	return output
}

func validate(input string) error {
	isValid := co.IsValidPat(input)
	if !isValid {
		return errors.New("invalid patp")
	}

	size, err := co.Clan(input)
	if size != "star" || err != nil {
		return errors.New("must be a star")
	}

	return nil
}

func writeResults(parent string, strategy string, results []string) error {
	fmt.Printf("Writing output for %s ...\n", strategy)

	path := fmt.Sprintf("./output/%s/%s_planets.txt", parent, strategy)
	return os.WriteFile(path, []byte(strings.Join(results, "\n")), 0o644)
}

func main() {
	var parent string

	argLength := len(os.Args[1:])
	if argLength > 0 {
		arg := os.Args[1]
		aErr := validate(arg)
		if aErr == nil {
			parent = arg
		}
	}

	if parent == "" {
		parentPrompt := promptui.Prompt{
			Label:    "Which star? (e.g., ~marzod)",
			Validate: validate,
		}

		var pErr error
		parent, pErr = parentPrompt.Run()

		if pErr != nil {
			fmt.Fprintf(os.Stderr, "bad input: %v\n", pErr)
			os.Exit(1)
		}
	}

	planets, err := makePlanets(parent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	anyApproxPlanets := filterPlanets(planets, AnyApprox)
	onlyApproxPlanets := filterPlanets(planets, OnlyApprox)
	anyEnglishPlanets := filterPlanets(planets, AnyEnglish)
	onlyEnglishPlanets := filterPlanets(planets, OnlyEnglish)
	doublesPlanets := filterPlanets(planets, Doubles)
	alliterationPlanets := filterPlanets(planets, Alliteration)

	deSiggedParent := strings.Replace(parent, "~", "", 1)
	if err := os.MkdirAll(fmt.Sprintf("./output/%s", deSiggedParent), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir: %v\n", err)
		os.Exit(1)
	}

	outputs := []struct {
		name    string
		planets []string
	}{
		{"any_approx", anyApproxPlanets},
		{"only_approx", onlyApproxPlanets},
		{"any_english", anyEnglishPlanets},
		{"only_english", onlyEnglishPlanets},
		{"doubles", doublesPlanets},
		{"alliteration", alliterationPlanets},
	}
	for _, out := range outputs {
		if err := writeResults(deSiggedParent, out.name, out.planets); err != nil {
			fmt.Fprintf(os.Stderr, "write %s: %v\n", out.name, err)
			os.Exit(1)
		}
	}

	fmt.Println("Done :)")
}
