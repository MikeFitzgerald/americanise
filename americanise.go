/*
LESSONS

1. You can defer a function's execution to occur after the enclosing function has returned by using the defer keyword

2. panics can be used as a general-purpose exception-handling mechanism, but doing so is considered poor Go practice;
   good practice is returning an error value as the last (or only) value; calling function should check this

3. Go differentiates errors from exceptions whereas C++, Java, and Python use exceptions for both mechanisms

*/

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	return inFilename, outFilename, nil
}