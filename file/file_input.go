package file

import (
	"bufio"
	"os"
	"strings"
)

func FileToString(fileName string) (sentences []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}

	defer file.Close()

	s := bufio.NewScanner(file)
	for s.Scan() {
		r := s.Text()
		if len(r) > 0 {
			t := strings.TrimFunc(r, func(x rune) bool {
				return x == '!' || x == '.' || x == '?'
			})
			sentences = append(sentences, t)
		}
	}

	return
}
