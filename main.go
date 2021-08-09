package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	convert, err := fileToString("input.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(convert)
}

func fileToString(fileName string) (string, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	str := string(file)

	return str, nil
}
