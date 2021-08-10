package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
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
	start := time.Now()
	sentences, err := fileToString("armageddon.txt")
	if err != nil {
		fmt.Println(err)
	}

	contents := sentencesToConcordance(sentences)

	printConcordance(contents)
	fmt.Printf("Task completed in %v \n", time.Since(start))
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v kb", bTokb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v kb", bTokb(m.TotalAlloc))
	fmt.Printf("\tSys = %v kb", bTokb(m.Sys))
	fmt.Printf("\tNum of GC cycles = %v\n", m.NumGC)
}

func bTokb(b uint64) uint64 {
	return b / 1024
}

func printConcordance(c concordance) error {
	encodedFile, err := os.Create("index.txt")
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
		//fmt.Fprintf(b, "%v. %v {%v:%v} \n", concordance_index[i], key, value.Count, locationsDisplay)
		fmt.Fprintf(b, "%v. %v {%v:%v} \n", concordanceIndex(i), key, value.Count, locationsDisplay)
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

//does this have to be int32?
func concordanceIndex(i int) (index string) {
	i--
	if firstLetter := i / 26; firstLetter > 0 {
		index += concordanceIndex(firstLetter)
		index += string(rune('a' + i%26))
	} else {
		index += string(rune('a' + i))
	}
	return
}

func sentencesToConcordance(sentences []string) concordance {
	contents := make(map[string]Details, 0)

	for sentenceIndex, sentence := range sentences {
		words := strings.Split(sentence, " ")
		//Need to scrub other non alphabetical characters " _ , ?
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
