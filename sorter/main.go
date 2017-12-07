package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func printErrorAndExit(err error) {
	fmt.Printf("Error: %s\n", err)
	fmt.Println("Usage: ./sorter {infile} {outfile} {tmp_dir} {tmp_file_size}")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 5 {
		printErrorAndExit(errors.New("invalid argument count"))
	}
	infile := os.Args[1]
	outfile := os.Args[2]
	tmpdir := os.Args[3]
	tmpFileSizeStr := os.Args[4]

	tmpFileSize, err := strconv.ParseUint(tmpFileSizeStr, 10, 64)
	if err != nil {
		printErrorAndExit(fmt.Errorf("could not parse tmp_file_size: %s", err))
	}

	sorter, err := NewSorter(infile, outfile, tmpdir, tmpFileSize)
	if err != nil {
		printErrorAndExit(fmt.Errorf("could not create sorter: %s", err))
	}
	if err := sorter.Sort(); err != nil {
		printErrorAndExit(fmt.Errorf("could not sort: %s", err))
	}
}
