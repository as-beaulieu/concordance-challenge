package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	sentences, err := fileToString("input.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sentences)

	contents := make(map[string]details, 0)

	for sentenceIndex, sentence := range sentences {
		fmt.Printf("Reading sentence #%d: %v \n", sentenceIndex, sentence)
		words := strings.Split(sentence, " ")
		fmt.Printf("Reading words after space split: %v \n", words)
		for _, word := range words {
			lowerCaseWord := strings.ToLower(word)
			d, exist := contents[lowerCaseWord]
			if exist {
				//found
				fmt.Println("HEY I WAS FOUND!!!: ", lowerCaseWord)
				d.count++
				d.locations = append(d.locations, sentenceIndex+1)
				contents[lowerCaseWord] = d
			} else {
				//not found
				newDetails := details{
					count:     1,
					locations: []int{sentenceIndex + 1},
				}
				contents[lowerCaseWord] = newDetails
			}
		}
	}

	fmt.Println(contents)
}

type details struct {
	count     int
	locations []int
}

func fileToString(fileName string) (sentences []string, err error) {
	//file, err := ioutil.ReadFile(fileName)
	//if err != nil {
	//	return "", err
	//}
	//
	//str := string(file)
	//
	//return str, nil

	file, err := os.Open(fileName)
	if err != nil {
		return
	}

	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		r := s.Text()
		if len(r) > 0 {
			t := strings.Trim(r, ".")
			sentences = append(sentences, t)
		}
	}

	return
}
