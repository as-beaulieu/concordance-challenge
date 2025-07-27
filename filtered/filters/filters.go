package filters

import (
	"concordance-challenge/filtered/models"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

func ApplyFilters(word string, filters models.FilterParameters) models.FilterCategory {
	if _, here := filters.Actions[word]; here {
		return models.FilterActions
	}
	if _, here := filters.Frameworks[word]; here {
		return models.FilterLanguages
	}
	if _, here := filters.Languages[word]; here {
		return models.FilterLanguages
	}
	if _, here := filters.Leadership[word]; here {
		return models.FilterLeaderships
	}
	if _, here := filters.Tools[word]; here {
		return models.FilterTools
	}
	if _, here := filters.Trivial[word]; here {
		return models.FilterTrivial
	}
	return ""
}

func FilterWords(processes models.Processes, filters models.FilterParameters) {

	defer close(processes.ActionsProcess)
	defer close(processes.FrameworksProcess)
	defer close(processes.LanguagesProcess)
	defer close(processes.LeadershipProcess)
	defer close(processes.OthersProcess)
	defer close(processes.ToolsProcess)
	defer close(processes.TrivialProcess)

	for line := range processes.FileProcess {
		spaceSplit := strings.Split(line, " ")
		for _, word := range spaceSplit {
			lowerWord := strings.ToLower(word)
			filter := ApplyFilters(lowerWord, filters)
			switch {
			case filter == models.FilterActions:
				processes.ActionsProcess <- lowerWord
			case filter == models.FilterFrameworks:
				processes.FrameworksProcess <- lowerWord
			case filter == models.FilterLanguages:
				processes.LanguagesProcess <- lowerWord
			case filter == models.FilterLeaderships:
				processes.LeadershipProcess <- lowerWord
			case filter == models.FilterTools:
				processes.ToolsProcess <- lowerWord
			case filter == models.FilterTrivial:
				processes.TrivialProcess <- lowerWord
			default:
				processes.OthersProcess <- lowerWord
			}
		}
	}
}

func ReadFilters() (models.FilterParameters, error) {
	file, err := os.Open("filters.json")
	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		return models.FilterParameters{}, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return models.FilterParameters{}, err
	}

	var filters models.FiltersInput
	err = json.Unmarshal(bytes, &filters)
	if err != nil {
		fmt.Printf("Error parsing file: %v", err)
		return models.FilterParameters{}, err
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

	return models.FilterParameters{
		Actions:    actionKeys,
		Frameworks: frameworksKeys,
		Languages:  languagesKeys,
		Leadership: leadershipKeys,
		Tools:      toolsKeys,
		Trivial:    trivialKeys,
	}, nil
}
