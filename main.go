package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	concordance_index = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj", "kk", "ll", "mm", "nn", "oo", "pp", "qq", "rr", "ss", "tt", "uu", "vv", "ww", "xx", "yy", "zz",
	}
)

type (
	Details struct {
		Count     int
		Locations []int
	}

	concordance map[string]Details
)

func main() {
	sentences, err := fileToString("input.txt")
	if err != nil {
		fmt.Println(err)
	}

	contents := sentencesToConcordance(sentences)

	printConcordance(contents)
}

func printConcordance(c concordance) error {
	encodedFile, err := os.Create("concordance.txt")
	if err != nil {
		return err
	}

	b := new(bytes.Buffer)

	keys := make([]string, 0, len(c))
	for key := range c {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for i, key := range keys {
		value := c[key]
		locations := value.Locations
		valuesText := make([]string, 0, len(locations))
		for valueIdx := range locations {
			number := locations[valueIdx]
			text := strconv.Itoa(number)
			valuesText = append(valuesText, text)
		}
		locationsDisplay := strings.Join(valuesText, ",")
		fmt.Fprintf(b, "%v. %v {%v:%v} \n", concordance_index[i], key, value.Count, locationsDisplay)
	}

	wrote, err := encodedFile.WriteString(b.String())
	if err != nil {
		return err
	}

	fmt.Printf("wrote %d bytes \n", wrote)

	if err := encodedFile.Sync(); err != nil {
		return err
	}

	return nil
}

func sentencesToConcordance(sentences []string) concordance {
	contents := make(map[string]Details, 0)

	for sentenceIndex, sentence := range sentences {
		words := strings.Split(sentence, " ")
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
