package export

import (
	"concordance-challenge/filtered/models"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func WriteFinalWordCountCSV(path string, fw models.FinalWordCount) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating csv file: %w", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	if err := w.Write([]string{"Category", "Word", "Count"}); err != nil {
		return fmt.Errorf("writing header: %w", err)
	}

	writeCategory := func(category string, entries []models.WordCount) error {
		for _, wc := range entries {
			record := []string{
				category,
				wc.Word,
				strconv.Itoa(wc.Count),
			}
			if err := w.Write(record); err != nil {
				return fmt.Errorf("writing record for %s: %w", category, err)
			}
		}
		return nil
	}

	if err := writeCategory("Actions", fw.Actions); err != nil {
		return err
	}
	if err := writeCategory("Frameworks", fw.Frameworks); err != nil {
		return err
	}
	if err := writeCategory("Languages", fw.Langauges); err != nil {
		return err
	}
	if err := writeCategory("Leadership", fw.Leadership); err != nil {
		return err
	}
	if err := writeCategory("Others", fw.Others); err != nil {
		return err
	}
	if err := writeCategory("Tools", fw.Tools); err != nil {
		return err
	}
	if err := writeCategory("Trivial", fw.Trivial); err != nil {
		return err
	}

	if err := w.Error(); err != nil {
		return fmt.Errorf("csv writer error: %w", err)
	}
	return nil
}
