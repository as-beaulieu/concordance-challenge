package main

import (
	"fmt"
	"strings"
	"sync"
)

// For this challenge you will need to implement a code that will count the frequency of each word in the textData variable.
// To do that, you must consider using channels and goroutines, like producer(goroutine1) and consumer(goroutine2).
// The final result will be the individual number of each word.
// Pre-requisites:
// 2 goroutines
// 1 goroutine must be responsible to manage the text
// 1 goroutine to count the words

func main() {
	processing := make(chan string)
	var wg sync.WaitGroup

	textData := "This is a document with some words Here are more words in another document This document repeats some words"

	wg.Add(2)
	go processText(textData, processing, &wg)

	go countWords(processing, &wg)

	wg.Wait()
}

func processText(text string, processing chan string, wg *sync.WaitGroup) {
	split := strings.Split(text, " ")
	for _, item := range split {
		processing <- item
	}
	defer close(processing)
	wg.Done()
	return
}

func countWords(collection chan string, wg *sync.WaitGroup) {
	concordance := make(map[string]int, 0)
	for word := range collection {
		if _, ok := concordance[word]; !ok {
			concordance[word] = 1
		} else {
			concordance[word]++
		}
		fmt.Println("processed! ", word)
	}
	fmt.Println("Done! WordCount: ", concordance)

	wg.Done()
	return
}
