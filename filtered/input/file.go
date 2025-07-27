package input

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"
)

func ReadFile(fileProcess chan string) {
	defer close(fileProcess)
	// for safety against larger files, opting for streaming file line by line
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		var cleaned strings.Builder
		for _, r := range line {
			if r == '-' {
				cleaned.WriteRune(' ')
			}
			if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
				cleaned.WriteRune(r)
			}
		}
		fileProcess <- cleaned.String()
	}

	// Check for any errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error during scanning: %v", err)
	}
}
