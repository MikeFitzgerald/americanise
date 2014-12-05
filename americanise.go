/*
LESSONS

1. You can defer a function's execution to occur after the enclosing function has returned by using the defer keyword

2. panics can be used as a general-purpose exception-handling mechanism, but doing so is considered poor Go practice;
   good practice is returning an error value as the last (or only) value; calling function should check this

3. Go differentiates errors from exceptions whereas C++, Java, and Python use exceptions for both mechanisms

4. Go supports the use of bare returns which simply return the values of parameters; bare returns are considered poor Go style.

5. Subtle scoping issue when using named return values.  Avoid assignments using the := operator to avoid declaring a duplicate variable.
   This can be avoided by specifying the return values explicitly instead of using named return values.

*/

package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	inFilename, outFilename, err := filenamesFromCommandLine()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	inFile, outFile := os.Stdin, os.Stdout
	if inFilename != "" {
		if inFile, err = os.Open(inFilename); err != nil {
			log.Fatal(err)
		}
		// Defer the following function call until after the enclosing function (main) returns
		defer inFile.Close()
	}
	if outFilename != "" {
		if outFile, err = os.Create(outFilename); err != nil {
			log.Fatal(err)
		}
		// Defer the following function call until after the enclosing function (main) returns
		defer outFile.Close()
	}
	if err = americanise(inFile, outFile); err != nil {
		log.Fatal(err)
	}
}

func filenamesFromCommandLine() (inFilename, outFilename string, err error) {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		err = fmt.Errorf("usage: %s [<]infile.txt [>]outfile.txt", filepath.Base(os.Args[0]))
		return "", "", err
	}
	if len(os.Args) > 1 {
		inFilename = os.Args[1]
		if len(os.Args) > 2 {
			outFilename = os.Args[2]
		}
	}
	if inFilename != "" && inFilename == outFilename {
		log.Fatal("won't overwrite the infile")
	}
	return inFilename, outFilename, nil
}

var britishAmerican = "british-american.txt"

// Function doesn't care what it's reading.
// Supports anything that supports the io.Reader and io.Writer interfaces.
func americanise(inFile io.Reader, outFile io.Writer) (err error) {
	reader := bufio.NewReader(inFile)
	writer := bufio.NewWriter(outFile)

	// Anonymous function that flushes the writer's buffer before the americanise() function returns control to its caller.
	defer func() {
		if err == nil {
			err = writer.Flush()
		}
	}()

	// Create a function signature that takes a string and returns a func(string) as a string ???
	var replacer func(string) string
	if replacer, err = makeReplacerFunction(britishAmerican); err != nil {
		return err
	}
	wordRx := regexp.MustCompile("[A-Za-z]+")
	eof := false
	for !eof {
		var line string
		line, err = reader.ReadString('\n')
		if err == io.EOF {
			err = nil	// io.EOF isn't really an error
			eof = true	// this will end the loop at the next iteration
		} else if err != nil {
			return err // finish immediately for real errors
		}

		line = wordRx.ReplaceAllStringFunc(line, replacer)
		// If we had a very small replacer function, say, one that simply upper-cased the words it 
		// matched, we could have created it as an anonymous function when we called the 
		// replacement function like so:
		// line = wordRx.ReplaceAllStringFunc(line, 
		//     func(word string) string { return strings.ToUpper(word) })

		if _, err = writer.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}

func makeReplacerFunction(file string) (func(string) string, error) {
	rawBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	text := string(rawBytes)

	// Make a map type and return its REFERENCE
	usForBritish := make(map[string]string)
	
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 2 {
			usForBritish[fields[0]] = fields[1]
		}
	}

	return func(word string) string {
		if usWord, found := usForBritish[word]; found {
			return usWord
		}
		return word
	}, nil
}