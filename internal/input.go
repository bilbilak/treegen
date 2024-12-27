package internal

import (
	"bufio"
	"os"
)

func ProcessInput(args []string) bool {
	processed := false

	if len(args) > 0 {
		for _, path := range args {
			tree := readFromFile(path)

			if len(tree) > 0 {
				generate(tree)
				processed = true
			}
		}
	}

	tree := readFromStdIn()

	if len(tree) > 0 {
		generate(tree)
		processed = true
	}

	return processed
}

func readFromFile(path string) []string {
	file, err := os.Open(path)

	if err != nil {
		FatalError(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		FatalError(err)
	}

	return lines
}

func readFromStdIn() []string {
	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return []string{}
	}

	scanner := bufio.NewScanner(os.Stdin)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		FatalError(err)
	}

	return lines
}
