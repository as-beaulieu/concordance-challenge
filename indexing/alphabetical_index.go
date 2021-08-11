package indexing

func AlphabeticalIndex(i int) (index string) {
	i--
	if firstLetter := i / 26; firstLetter > 0 {
		index += AlphabeticalIndex(firstLetter)
		index += string(rune('a' + i%26))
	} else {
		index += string(rune('a' + i))
	}
	return
}
