package main

import (
	"concordance-challenge/file-concordance/concordance"
	file2 "concordance-challenge/file-concordance/file"
	"fmt"
)

func main() {
	sentences, err := file2.FileToString("input.txt")
	if err != nil {
		fmt.Println(err)
	}

	contents := concordance.SentencesToConcordance(sentences)

	file2.PrintConcordance(contents)
}
