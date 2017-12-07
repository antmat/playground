package main

import (
	"bufio"
	"io"
	"os"
)

type Merger interface {
	Merge() error
}

type merger struct {
	paths      []string
	outPath    string
	readers    []*bufio.Reader
	topStrings []*string
}

func NewMerger(paths []string, outPath string) Merger {
	return &merger{paths: paths, outPath: outPath}
}

func (m *merger) Merge() error {
	// Always create file even if it is empty
	out, err := os.Create(m.outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if len(m.paths) == 0 {
		return nil
	}
	writer := bufio.NewWriter(out)
	defer writer.Flush()

	files := make([]*os.File, 0)
	defer func() {
		for _, file := range files {
			file.Close()
		}
	}()

	for _, path := range m.paths {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		files = append(files, file)
		m.readers = append(m.readers, bufio.NewReader(file))
	}

	m.initializeTopStrings()

	for {
		top, err := m.readTopString()
		if err != nil {
			return err
		}
		if top == nil {
			break
		}
		writer.WriteString(*top)
	}

	return nil
}

func (m *merger) readString(fileIdx int) error {
	reader := m.readers[fileIdx]
	if reader == nil {
		return nil
	}

	// TODO: This can result in OOM for huge lines,
	// Consider using some kind of limited buffer reader or even iobuf.Scanner
	str, err := reader.ReadString('\n')
	if err != nil {
		if len(str) != 0 {
			str = str + "\n"
		}
		if err == io.EOF {
			m.readers[fileIdx] = nil
		} else {
			return err
		}
	}
	if len(str) != 0 {
		m.topStrings[fileIdx] = &str
	} else {
		m.topStrings[fileIdx] = nil
	}
	return nil
}

func (m *merger) initializeTopStrings() error {
	m.topStrings = make([]*string, len(m.readers))
	for idx := range m.readers {
		if err := m.readString(idx); err != nil {
			return err
		}
	}
	return nil
}

func (m *merger) readTopString() (*string, error) {
	var min *string
	minIdx := 0
	for idx, str := range m.topStrings {
		if min == nil {
			min = str
			minIdx = idx
		} else {
			if str != nil && *str < *min {
				min = str
				minIdx = idx
			}
		}
	}
	if min == nil {
		return min, nil
	}
	if err := m.readString(minIdx); err != nil {
		return nil, err
	}
	return min, nil
}
