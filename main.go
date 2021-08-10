package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

func main() {
	sentences, err := fileToString("input.txt")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(sentences)

	contents := sentencesToConcordance(sentences)

	//fmt.Println(contents)

	printConcordance(contents)
}

type Details struct {
	Count     int
	Locations []int
}

type concordance map[string]Details

func printConcordance(c concordance) {
	encodedFile, err := os.Create("concordance.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	e := gob.NewEncoder(encodedFile)

	//Now convert the map to text
	b := new(bytes.Buffer)
	i := 0
	for key, value := range c {
		//a. a {2:1,1}
		fmt.Fprintf(b, "%v. %v {%v:%v} \n", i, key, value.Count, value.Locations)
		i++
	}

	if err := e.Encode(b.String()); err != nil {
		fmt.Println(err)
		return
	}
}

func sentencesToConcordance(sentences []string) concordance {
	contents := make(map[string]Details, 0)

	for sentenceIndex, sentence := range sentences {
		//fmt.Printf("Reading sentence #%d: %v \n", sentenceIndex, sentence)
		words := strings.Split(sentence, " ")
		//fmt.Printf("Reading words after space split: %v \n", words)
		for _, word := range words {
			lowerCaseWord := strings.ToLower(word)
			d, exist := contents[lowerCaseWord]
			if exist {
				d.Count++
				d.Locations = append(d.Locations, sentenceIndex+1)
				contents[lowerCaseWord] = d
			} else {
				newDetails := Details{
					Count:     1,
					Locations: []int{sentenceIndex + 1},
				}
				contents[lowerCaseWord] = newDetails
			}
		}
	}
	return contents
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
