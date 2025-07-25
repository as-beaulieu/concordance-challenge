package file

import (
	"bytes"
	"concordance-challenge/file-concordance/concordance"
	"concordance-challenge/file-concordance/indexing"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func PrintConcordance(c concordance.Concordance) error {
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
		fmt.Fprintf(b, "%v. %v {%v:%v} \n", indexing.AlphabeticalIndex(i+1), key, value.Count, locationsDisplay)
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
