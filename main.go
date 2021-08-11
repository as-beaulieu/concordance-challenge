package main

import (
	"concordance-challenge/concordance"
	"concordance-challenge/file"
	"fmt"
)

func main() {
	sentences, err := file.FileToString("input.txt")
	if err != nil {
		fmt.Println(err)
	}

	contents := concordance.SentencesToConcordance(sentences)

	file.PrintConcordance(contents)
}
