package concordance

import "strings"

type (
	Details struct {
		Count     int
		Locations []int
	}

	Concordance map[string]Details
)

func SentencesToConcordance(sentences []string) Concordance {
	contents := make(map[string]Details, 0)

	for sentenceIndex, sentence := range sentences {
		words := strings.Split(sentence, " ")
		for _, word := range words {
			cleanedWord := trimWord(word)
			d, exist := contents[cleanedWord]
			if exist {
				d.Count++
				d.Locations = append(d.Locations, sentenceIndex+1)
				contents[cleanedWord] = d
			} else {
				newDetails := Details{
					Count:     1,
					Locations: []int{sentenceIndex + 1},
				}
				contents[cleanedWord] = newDetails
			}
		}
	}
	return contents
}

func trimWord(word string) string {
	trimmedWordLeft := strings.TrimLeftFunc(word, func(r rune) bool {
		return r == '"' || r == '_' || r == ',' || r == '(' || r == ')' || r == '*' || r == '-'
	})
	trimmedWordRight := strings.TrimRightFunc(trimmedWordLeft, func(r rune) bool {
		return r == '"' || r == '_' || r == ',' || r == '(' || r == ')' || r == '*' || r == '-'
	})
	return strings.ToLower(trimmedWordRight)
}
