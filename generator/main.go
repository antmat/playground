package main

import (
	"errors"
	"os"
	"strconv"
	"fmt"
)

func printErrorAndExit(err error) {
	fmt.Printf("Error: %s\n", err)
	fmt.Println("Usage: ./generator {outfile} {line_count} {line_length}")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 4 {
		printErrorAndExit(errors.New("invalid argument count"))
	}
	outfile := os.Args[1]
	lineCountStr := os.Args[2]
	lineLengthStr := os.Args[3]

	lineCount, err :=  strconv.ParseUint(lineCountStr, 10, 64)
	if err != nil {
		printErrorAndExit(fmt.Errorf("could not parse line_count - %s", err))
	}

	lineLength, err := strconv.ParseUint(lineLengthStr, 10, 64)
	if err != nil {
		printErrorAndExit(fmt.Errorf("could not parse line_length - %s", err))
	}

	g, err := NewGenerator(outfile, lineCount, lineLength)
	if err != nil {
		printErrorAndExit(err)
	}

	err = g.Run()
	if err != nil {
		printErrorAndExit(err)
	}
}