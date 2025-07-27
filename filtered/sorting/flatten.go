package sorting

import (
	"concordance-challenge/filtered/models"
	"sort"
)

func MapToSortedSlice(list map[string]int) []models.WordCount {
	wordCounts := make([]models.WordCount, 0, len(list))
	for word, count := range list {
		wordCounts = append(wordCounts, models.WordCount{Word: word, Count: count})
	}
	sort.Slice(wordCounts, func(i, j int) bool {
		if wordCounts[i].Count == wordCounts[j].Count {
			return wordCounts[i].Word < wordCounts[j].Word
		}
		return wordCounts[i].Count > wordCounts[j].Count
	})
	return wordCounts
}
