package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {
	// Read the small and big files from the command line arguments
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <small_file> <big_file>\n", os.Args[0])
		os.Exit(1)
	}
	smallFileName := os.Args[1]
	bigFileName := os.Args[2]

	smallFile, err := os.Open(smallFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening small file: %v\n", err)
		os.Exit(1)
	}
	defer smallFile.Close()

	// Read the small file into memory
	smallFileLines := make(map[string][]int)
	scanner := bufio.NewScanner(smallFile)
	for smallLineNum := 1; scanner.Scan(); smallLineNum++ {
		line := strings.TrimSpace(scanner.Text())
		smallFileLines[line] = append(smallFileLines[line], smallLineNum)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading small file: %v\n", err)
		os.Exit(1)
	}

	// Open the big file
	bigFile, err := os.Open(bigFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening big file: %v\n", err)
		os.Exit(1)
	}
	defer bigFile.Close()

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup
	lineChan := make(chan string)

	// Start a goroutine to read the big file
	go func() {
		defer close(lineChan)
		scanner := bufio.NewScanner(bigFile)
		for bigLineNum := 1; scanner.Scan(); bigLineNum++ {
			line := strings.TrimSpace(scanner.Text())
			lineChan <- line
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading big file: %v\n", err)
			os.Exit(1)
		}
	}()

	// Start multiple goroutines to process lines from the big file
	duplicates := 0
	var mu sync.Mutex
	for i := 0; i < 4; i++ { // Adjust the number of goroutines as needed
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lineChan {
				if smallLineNums, found := smallFileLines[line]; found {
					mu.Lock()
					duplicates++
					fmt.Printf("Duplicate found: %s\n", line)
					fmt.Printf("  Corresponds to line(s) in the small file: %v\n", smallLineNums)
					mu.Unlock()
				}
			}
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Printf("Total duplicates found: %d\n", duplicates)
}
