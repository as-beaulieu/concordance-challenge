package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var (
	FilterActions     FilterCategory = "actions"
	FilterFrameworks  FilterCategory = "frameworks"
	FilterLanguages   FilterCategory = "language"
	FilterLeaderships FilterCategory = "leadership"
	FilterTools       FilterCategory = "tools"
	FilterTrivial     FilterCategory = "trivial"
)

type FilterCategory string

type FiltersInput struct {
	Actions    []string `json:"actions"`
	Frameworks []string `json:"frameworks"`
	Languages  []string `json:"languages"`
	Leadership []string `json:"leadership"`
	Tools      []string `json:"tools"`
	Trivial    []string `json:"trivial"`
}

type FilterParameters struct {
	actions    map[string]bool
	frameworks map[string]bool
	languages  map[string]bool
	leadership map[string]bool
	tools      map[string]bool
	trivial    map[string]bool
}

type Processes struct {
	fileProcess       chan string
	actionsProcess    chan string
	frameworksProcess chan string
	languagesProcess  chan string
	leadershipProcess chan string
	othersProcess     chan string
	toolsProcess      chan string
	trivialProcess    chan string
}

type WordCount struct {
	Word  string
	Count int
}

type FinalWordCount struct {
	Actions    []WordCount `json:"actions"`
	Frameworks []WordCount `json:"frameworks"`
	Langauges  []WordCount `json:"langauges"`
	Leadership []WordCount `json:"leadership"`
	Others     []WordCount `json:"others"`
	Tools      []WordCount `json:"tools"`
	Trivial    []WordCount `json:"trivial"`
}

func main() {
	processes := Processes{
		fileProcess:       make(chan string),
		actionsProcess:    make(chan string),
		frameworksProcess: make(chan string),
		languagesProcess:  make(chan string),
		leadershipProcess: make(chan string),
		othersProcess:     make(chan string),
		toolsProcess:      make(chan string),
		trivialProcess:    make(chan string),
	}

	var wg sync.WaitGroup
	var actionsList,
		frameworksList,
		languagesList,
		leadershipList,
		othersList,
		toolsList,
		trivialList map[string]int

	// only want keywords, so no reason to keep non-context words (and, the, etc.)
	// we'll grab a list of filters from a json file
	fmt.Println("gathering filters")
	filters, err := ReadFilters()
	if err != nil {
		log.Fatal(err)
	}

	//wg.Add(6)
	go func() {
		wg.Add(1)
		defer wg.Done()
		ReadFile(processes.fileProcess)
	}()

	wg.Add(1)
	go FilterWords(&wg, processes, filters)

	// "fan out" - pattern which spreads work out to multiple paths
	go func() {
		wg.Add(1)
		defer wg.Done()
		actionsList = BuildList(processes.actionsProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		languagesList = BuildList(processes.frameworksProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		languagesList = BuildList(processes.languagesProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		leadershipList = BuildList(processes.leadershipProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		othersList = BuildList(processes.othersProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		toolsList = BuildList(processes.toolsProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		trivialList = BuildList(processes.trivialProcess)
	}()

	wg.Wait()

	fmt.Println("lists generation complete, now sorting")

	finalReport := FinalWordCount{
		Actions:    MapToSortedSlice(actionsList),
		Frameworks: MapToSortedSlice(frameworksList),
		Langauges:  MapToSortedSlice(languagesList),
		Leadership: MapToSortedSlice(leadershipList),
		Others:     MapToSortedSlice(othersList),
		Tools:      MapToSortedSlice(toolsList),
		Trivial:    MapToSortedSlice(trivialList),
	}

	if err := WriteFinalWordCountCSV("report.csv", finalReport); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
}

func WriteFinalWordCountCSV(path string, fw FinalWordCount) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating csv file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"Category", "Word", "Count"}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	writeCategory := func(category string, entries []WordCount) error {
		for _, wc := range entries {
			record := []string{
				category,
				wc.Word,
				strconv.Itoa(wc.Count),
			}
			if err := w.Write(record); err != nil {
				return fmt.Errorf("writing record for %s: %w", category, err)
			}
		}
		return nil
	}

	if err := writeCategory("Actions", fw.Actions); err != nil {
		return err
	}
	if err := writeCategory("Frameworks", fw.Frameworks); err != nil {
		return err
	}
	if err := writeCategory("Languages", fw.Langauges); err != nil {
		return err
	}
	if err := writeCategory("Leadership", fw.Leadership); err != nil {
		return err
	}
	if err := writeCategory("Others", fw.Others); err != nil {
		return err
	}
	if err := writeCategory("Tools", fw.Tools); err != nil {
		return err
	}
	if err := writeCategory("Trivial", fw.Trivial); err != nil {
		return err
	}

	if err := w.Error(); err != nil {
		return fmt.Errorf("csv writer error: %w", err)
	}
	return nil
}

func MapToSortedSlice(list map[string]int) []WordCount {
	wordCounts := make([]WordCount, 0, len(list))
	for word, count := range list {
		wordCounts = append(wordCounts, WordCount{Word: word, Count: count})
	}
	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].Count == wordCounts[j].Count {
			return wordCounts[i].Word < wordCounts[j].Word
		}
		return wordCounts[i].Count > wordCounts[j].Count
	})
	return wordCounts
}

func FilterWords(wg *sync.WaitGroup, processes Processes, filters FilterParameters) {

	defer close(processes.actionsProcess)
	defer close(processes.frameworksProcess)
	defer close(processes.languagesProcess)
	defer close(processes.leadershipProcess)
	defer close(processes.othersProcess)
	defer close(processes.toolsProcess)
	defer close(processes.trivialProcess)

	for line := range processes.fileProcess {
		spaceSplit := strings.Split(line, " ")
		for _, word := range spaceSplit {
			lowerWord := strings.ToLower(word)
			filter := ApplyFilters(lowerWord, filters)
			switch {
			case filter == FilterActions:
				processes.actionsProcess <- lowerWord
			case filter == FilterFrameworks:
				processes.frameworksProcess <- lowerWord
			case filter == FilterLanguages:
				processes.languagesProcess <- lowerWord
			case filter == FilterLeaderships:
				processes.leadershipProcess <- lowerWord
			case filter == FilterTools:
				processes.toolsProcess <- lowerWord
			case filter == FilterTrivial:
				processes.trivialProcess <- lowerWord
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
	if _, here := filters.frameworks[word]; here {
		return FilterLanguages
	}
	if _, here := filters.languages[word]; here {
		return FilterLanguages
	}
	if _, here := filters.leadership[word]; here {
		return FilterLeaderships
	}
	if _, here := filters.tools[word]; here {
		return FilterTools
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
	frameworksKeys := make(map[string]bool)
	for _, framework := range filters.Frameworks {
		frameworksKeys[framework] = true
	}
	languagesKeys := make(map[string]bool)
	for _, language := range filters.Languages {
		languagesKeys[language] = true
	}
	leadershipKeys := make(map[string]bool)
	for _, leadership := range filters.Leadership {
		leadershipKeys[leadership] = true
	}
	toolsKeys := make(map[string]bool)
	for _, tool := range filters.Tools {
		toolsKeys[tool] = true
	}
	trivialKeys := make(map[string]bool)
	for _, trivialWord := range filters.Trivial {
		trivialKeys[trivialWord] = true
	}

	return FilterParameters{
		actions:    actionKeys,
		frameworks: frameworksKeys,
		languages:  languagesKeys,
		leadership: leadershipKeys,
		tools:      toolsKeys,
		trivial:    trivialKeys,
	}, nil
}
