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

type FiltersInput struct {
	Trivial []string `json:"trivial"`
}

type FilterParameters struct {
	trivial map[string]bool
}

func main() {
	fileProcess := make(chan string)
	listProcess := make(chan string)

	var wg sync.WaitGroup
	var result map[string]int

	// only want keywords, so no reason to keep non-context words (and, the, etc.)
	// we'll grab a list of filters from a json file
	fmt.Println("gathering filters")
	filters, err := ReadFilters()
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(3)
	go ReadFile(&wg, fileProcess)

	go FilterWords(&wg, fileProcess, listProcess, filters)

	go func() {
		result = BuildList(listProcess)
		wg.Done()
	}()

	wg.Wait()

	fmt.Println(result)
	fmt.Println("Done")
}

func FilterWords(wg *sync.WaitGroup, fileProcess, listProcess chan string, filters FilterParameters) {
	defer close(listProcess)
	for line := range fileProcess {
		spaceSplit := strings.Split(line, " ")
		for _, word := range spaceSplit {
			if skip := ApplyFilters(word, filters); skip {
				continue
			}
			listProcess <- word
		}
	}

	wg.Done()
}

func BuildList(listProcess chan string) map[string]int {
	list := make(map[string]int)
	for word := range listProcess {
		lowerWord := strings.ToLower(word)
		if _, here := list[lowerWord]; here {
			list[lowerWord]++
		} else {
			list[lowerWord] = 1
		}
	}

	return list
}

func ApplyFilters(word string, filters FilterParameters) bool {
	if _, here := filters.trivial[word]; here {
		return true
	}
	return false
}

func ReadFile(wg *sync.WaitGroup, fileProcess chan string) {
	defer close(fileProcess)
	// for safety against larger files, opting for streaming file line by line
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	// Create a new Scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate through each line
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
	wg.Done()
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

	trivialKeys := make(map[string]bool)
	for _, trivialWord := range filters.Trivial {
		trivialKeys[trivialWord] = true
	}

	return FilterParameters{
		trivial: trivialKeys,
	}, nil
}
