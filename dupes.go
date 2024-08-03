package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// Read the big file from stdin
	bigFileLines, err := readLinesFromStdin()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		os.Exit(1)
	}

	// Read the small file from the command line argument
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <small_file>\n", os.Args[0])
		os.Exit(1)
	}
	smallFileName := os.Args[1]
	smallFile, err := os.Open(smallFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening small file: %v\n", err)
		os.Exit(1)
	}
	defer smallFile.Close()

	// Check for duplicates
	duplicates := 0
	scanner := bufio.NewScanner(smallFile)
	for smallLineNum := 1; scanner.Scan(); smallLineNum++ {
		line := strings.TrimSpace(scanner.Text())
		if bigLineNums, found := bigFileLines[line]; found {
			duplicates++
			fmt.Printf("Duplicate found on line %d of the small file: %s\n", smallLineNum, line)
			fmt.Printf("  Corresponds to line(s) in the big file: %v\n", bigLineNums)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading small file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total duplicates found: %d\n", duplicates)
}

func readLinesFromStdin() (map[string][]int, error) {
	lines := make(map[string][]int)
	reader := bufio.NewReader(os.Stdin)
	for lineNum := 1; ; lineNum++ {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		lines[line] = append(lines[line], lineNum)
	}
	return lines, nil
}

