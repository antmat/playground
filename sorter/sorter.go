package main

import (
	"fmt"

	"bufio"
	"io"
	"os"
	"sort"

	"github.com/pborman/uuid"
	"github.com/shirou/gopsutil/mem"
)

const (
	B              uint64 = 1
	Kb                    = B << 10
	Mb                    = Kb << 10
	Gb                    = Mb << 10
	Tb                    = Gb << 10
	MinTmpFileSize        = 1 * Mb
	MaxTmpFileSize        = 1 * Tb

	TmpFileBaseName = ".sorter.tmp."
)

type Sorter interface {
	Sort() error
}

func NewSorter(infile string, outfile string, tmpDir string, tmpFileSize uint64) (Sorter, error) {
	v, _ := mem.VirtualMemory()

	if tmpFileSize > v.Total {
		return nil, fmt.Errorf("tmp_file_size should be less than system RAM amount")
	}

	if tmpFileSize < MinTmpFileSize {
		return nil, fmt.Errorf("tmp_file_size should be more than %d", MinTmpFileSize)
	}

	return &sorter{
		infile:          infile,
		outfile:         outfile,
		tmpDir:          tmpDir,
		tmpFileSize:     tmpFileSize,
		tmpFileIndex:    0,
		tmpUniqueSuffix: uuid.New(),
	}, nil
}

type sorter struct {
	infile string

	outfile string

	tmpDir          string
	tmpFileSize     uint64
	tmpFileIndex    uint
	tmpUniqueSuffix string
}

func (s *sorter) Sort() error {
	if err := s.splitAndSort(); err != nil {
		return fmt.Errorf("could not sort chunk: %s", err)
	}
	if err := s.mergeSortedFiles(); err != nil {
		return fmt.Errorf("could not open tmp file: %s", err)
	}
	return nil
}

func (s *sorter) openNextTmpFile() (*os.File, *bufio.Writer, error) {
	filename := s.tmpFileName(s.tmpFileIndex)
	s.tmpFileIndex++
	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	return f, bufio.NewWriter(f), nil
}

func (s *sorter) tmpFileName(idx uint) string {
	return s.tmpDir + "/" + TmpFileBaseName + s.tmpUniqueSuffix + "." + fmt.Sprint(idx)
}

func (s *sorter) splitAndSort() error {
	file, err := os.Open(s.infile)
	if err != nil {
		return err
	}

	defer file.Close()
	reader := bufio.NewReader(file)

	for {
		part, err := s.readAndSortPart(reader)
		if part != nil && len(part) > 0 {
			if err := s.writeTmpFile(part); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sorter) readAndSortPart(reader *bufio.Reader) ([]string, error) {
	strings := make([]string, 0)
	defer func() {
		sort.Strings(strings)
	}()

	read := uint64(0)

	for read < s.tmpFileSize {
		str, err := reader.ReadString('\n')
		if len(str) != 0 {
			if err != nil {
				str = str + "\n"
			}
			strings = append(strings, str)
			read += uint64(len(str))
		}
		if err != nil {
			return strings, err
		}
	}
	return strings, nil
}

func (s *sorter) writeTmpFile(sortedData []string) error {
	file, writer, err := s.openNextTmpFile()
	if err != nil {
		return err
	}
	defer file.Close()
	defer writer.Flush()
	for _, str := range sortedData {
		_, err := writer.WriteString(str)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *sorter) mergeSortedFiles() error {
	if s.tmpFileIndex == 1 {
		return os.Rename(s.tmpFileName(0), s.outfile)
	}
	files := make([]string, 0, s.tmpFileIndex)
	for i := uint(0); i < s.tmpFileIndex; i++ {
		files = append(files, s.tmpFileName(i))
	}

	merger := NewMerger(files, s.outfile)
	if err := merger.Merge(); err != nil {
		return err
	}
	return s.removeTmpFiles()
}

func (s *sorter) removeTmpFiles() error {
	for i := uint(0); i < s.tmpFileIndex; i++ {
		if err := os.Remove(s.tmpFileName(i)); err != nil {
			return err
		}
	}
	return nil
}
