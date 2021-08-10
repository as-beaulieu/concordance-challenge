# concordance-challenge

## How to run

1. [setup golang for your computer](https://golang.org/doc/install#download) 

from a terminal, locate the project directory

`cd untangle-challenge`

run the main.go file

`go run main.go`

## What is this doing?

The concordance-challenge program will take a predetermined `input.txt` file 
locally within the repo, and read it. While reading the text file, it is performing several
different actions:

- converting the input text file into a string object

- separating the string object into a slice of strings `[]string` based off of each sentence

- within each sentence, it is reading each word

- the program as a whole is creating a `concordance` - an alphabetical list of words, with
citations of the sentences in which the word was found.
  
- the concordance is printed into a text file `concordance.txt` located locally within the
project directory
  
## For fun

A copyright-expired book was found, and copied into a text file, `armageddon.txt` to test 
the capabilities of this program against a larger text sample. Either change the parameter
in the `fileToString()` function from `input.txt` to `armageddon.txt`, or rename the text
file itself to input.txt, and run the results.