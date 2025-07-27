package main

import (
	"concordance-challenge/filtered/export"
	"concordance-challenge/filtered/filters"
	"concordance-challenge/filtered/input"
	"concordance-challenge/filtered/models"
	"concordance-challenge/filtered/sorting"
	"fmt"
	"log"
	"sync"
)

func main() {
	processes := models.Processes{
		FileProcess:       make(chan string),
		ActionsProcess:    make(chan string),
		FrameworksProcess: make(chan string),
		LanguagesProcess:  make(chan string),
		LeadershipProcess: make(chan string),
		OthersProcess:     make(chan string),
		ToolsProcess:      make(chan string),
		TrivialProcess:    make(chan string),
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
	filter, err := filters.ReadFilters()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		input.ReadFile(processes.FileProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		filters.FilterWords(processes, filter)
	}()

	// "fan out" - pattern which spreads work out to multiple paths
	go func() {
		wg.Add(1)
		defer wg.Done()
		actionsList = BuildList(processes.ActionsProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		languagesList = BuildList(processes.FrameworksProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		languagesList = BuildList(processes.LanguagesProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		leadershipList = BuildList(processes.LeadershipProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		othersList = BuildList(processes.OthersProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		toolsList = BuildList(processes.ToolsProcess)
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		trivialList = BuildList(processes.TrivialProcess)
	}()

	wg.Wait()

	fmt.Println("lists generation complete, now sorting")

	finalReport := models.FinalWordCount{
		Actions:    sorting.MapToSortedSlice(actionsList),
		Frameworks: sorting.MapToSortedSlice(frameworksList),
		Langauges:  sorting.MapToSortedSlice(languagesList),
		Leadership: sorting.MapToSortedSlice(leadershipList),
		Others:     sorting.MapToSortedSlice(othersList),
		Tools:      sorting.MapToSortedSlice(toolsList),
		Trivial:    sorting.MapToSortedSlice(trivialList),
	}

	if err := export.WriteFinalWordCountCSV("report.csv", finalReport); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done")
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
