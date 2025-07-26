package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"unicode"
)

var (
	FilterActions   FilterCategory = "actions"
	FilterLanguages FilterCategory = "language"
	FilterTrivial   FilterCategory = "trivial"
)

type FilterCategory string

type FiltersInput struct {
	Actions   []string `json:"actions"`
	Languages []string `json:"languages"`
	Trivial   []string `json:"trivial"`
}

type FilterParameters struct {
	actions   map[string]bool
	languages map[string]bool
	trivial   map[string]bool
}

type Processes struct {
	fileProcess      chan string
	actionsProcess   chan string
	languagesProcess chan string
	othersProcess    chan string
}

func main() {
	processes := Processes{
		fileProcess:      make(chan string),
		actionsProcess:   make(chan string),
		languagesProcess: make(chan string),
		othersProcess:    make(chan string),
	}

	var wg sync.WaitGroup
	var actionsList, languagesList, othersList map[string]int

	// only want keywords, so no reason to keep non-context words (and, the, etc.)
	// we'll grab a list of filters from a json file
	fmt.Println("gathering filters")
	filters, err := ReadFilters()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(5)
	go func() {
		defer wg.Done()
		ReadFile(processes.fileProcess)
	}()

	go FilterWords(&wg, processes, filters)

	// "fan out" - pattern which spreads work out to multiple paths
	go func() {
		defer wg.Done()
		actionsList = BuildList(processes.actionsProcess)
	}()

	go func() {
		defer wg.Done()
		languagesList = BuildList(processes.languagesProcess)
	}()

	go func() {
		defer wg.Done()
		othersList = BuildList(processes.othersProcess)
	}()

	wg.Wait()

	fmt.Println("actionsList: ", actionsList)
	fmt.Println("languagesList: ", languagesList)
	fmt.Println("othersList: ", othersList)
	fmt.Println("Done")
}

func FilterWords(wg *sync.WaitGroup, processes Processes, filters FilterParameters) {
	defer close(processes.actionsProcess)
	defer close(processes.languagesProcess)
	defer close(processes.othersProcess)
	for line := range processes.fileProcess {
		spaceSplit := strings.Split(line, " ")
		for _, word := range spaceSplit {
			lowerWord := strings.ToLower(word)
			filter := ApplyFilters(lowerWord, filters)
			switch {
			case filter == FilterActions:
				processes.actionsProcess <- lowerWord
			case filter == FilterLanguages:
				processes.languagesProcess <- lowerWord
			case filter == FilterTrivial:
				continue
			default:
				processes.othersProcess <- lowerWord
			}
		}
	}

	wg.Done()
}

func BuildList(listProcess chan string) map[string]int {
	list := make(map[string]int)
	for word := range listProcess {
		if _, here := list[word]; here {
			list[word]++
		} else {
			list[word] = 1
		}
	}

	return list
}

func ApplyFilters(word string, filters FilterParameters) FilterCategory {
	if _, here := filters.actions[word]; here {
		return FilterActions
	}
	if _, here := filters.languages[word]; here {
		return FilterLanguages
	}
	if _, here := filters.trivial[word]; here {
		return FilterTrivial
	}
	return ""
}

func ReadFile(fileProcess chan string) {
	defer close(fileProcess)
	// for safety against larger files, opting for streaming file line by line
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var cleaned strings.Builder
		for _, r := range line {
			if r == '-' {
				cleaned.WriteRune(' ')
			}
			if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
				cleaned.WriteRune(r)
			}
		}
		fileProcess <- cleaned.String()
	}

	// Check for any errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during scanning: %v", err)
	}
}

func ReadFilters() (FilterParameters, error) {
	file, err := os.Open("filters.json")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return FilterParameters{}, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return FilterParameters{}, err
	}

	var filters FiltersInput
	err = json.Unmarshal(bytes, &filters)
	if err != nil {
		fmt.Printf("Error parsing file: %v", err)
		return FilterParameters{}, err
	}

	actionKeys := make(map[string]bool)
	for _, action := range filters.Actions {
		actionKeys[action] = true
	}
	languagesKeys := make(map[string]bool)
	for _, language := range filters.Languages {
		languagesKeys[language] = true
	}
	trivialKeys := make(map[string]bool)
	for _, trivialWord := range filters.Trivial {
		trivialKeys[trivialWord] = true
	}

	return FilterParameters{
		actions:   actionKeys,
		languages: languagesKeys,
		trivial:   trivialKeys,
	}, nil
}
