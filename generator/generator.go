package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const (
	B           uint64 = 1
	Kb                 = B << 10
	Mb                 = Kb << 10
	Gb                 = Mb << 10
	Tb                 = Gb << 10
	MaxFileSize        = 100 * Gb

	MaxLineLength = 1 * Mb
	BufferSize    = 10 * Mb
)

type Generator interface {
	Run() error
}

func NewGenerator(outfile string, lineCount uint64, lineLength uint64) (Generator, error) {
	if lineLength > MaxLineLength || lineLength == 0 {
		return nil, fmt.Errorf("line length can not be zero or greater than %d (got %d)", MaxLineLength, lineLength)
	}

	fileSize := lineLength * lineCount
	if fileSize > MaxFileSize {
		return nil, fmt.Errorf("resulting file size is too big (%d, maximum is %d)", fileSize, MaxLineLength)
	}

	return &generator{
		outfile:    outfile,
		lineCount:  lineCount,
		lineLength: lineLength,
	}, nil
}

type generator struct {
	outfile    string
	lineCount  uint64
	lineLength uint64
}

func (g *generator) appendString(writer *bufio.Writer) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lettersCnt := int32(len(letters))

	// each random uint64 is capable of generating 4 random bytes
	// this can boost random string generation a lot
	random := rand.Uint64()
	shift := uint(0)
	for k := 0; k < int(g.lineLength); k++ {
		if shift == 64 {
			// we used all bytes from random number
			shift = 0
			random = rand.Uint64()
		}
		// use (shift/8) byte from random number
		idx := (random & (0xff << shift)) >> shift
		shift = shift + 8
		writer.WriteByte(letters[idx%uint64(lettersCnt)])
	}
	writer.WriteByte('\n')
}

func (g *generator) Run() error {
	file, err := os.Create(g.outfile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for line := uint64(0); line < g.lineCount; line++ {
		g.appendString(writer)
	}
	return nil
}
