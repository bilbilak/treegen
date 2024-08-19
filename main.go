package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	version     = "1.0.2"
	bold        = "\033[1m"
	reset       = "\033[0m"
	gap         = "    "
	ascChild    = "├── "
	ascSubChild = "│   "
	ascEndChild = "└── "
	txtChild    = "|-- "
	txtSubChild = "|   "
	txtEndChild = "+-- "
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version number")
	helpFlag := flag.Bool("help", false, "Print help message")

	flag.Usage = func() {
		fmt.Println("Use --help for usage instructions")
	}
	flag.Parse()

	if *versionFlag {
		fmt.Println("treegen", version)
		return
	}

	if *helpFlag {
		helpMessage()
		return
	}

	var content string
	var err error
	if flag.NArg() > 0 {
		for _, path := range flag.Args() {
			content, err = readFromFile(path)
			if err != nil {
				fmt.Fprint(os.Stderr, err.Error())
				os.Exit(1)
			}
		}
	} else {
		content, err = readFromStdIn()
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	err = generate(content)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func helpMessage() {
	fmt.Println(fmt.Sprintf("\n%sUsage:%s treegen [OPTIONS] [STDIN|FILE...]", bold, reset))
	fmt.Println(fmt.Sprintf("\n%sOptions:%s", bold, reset))
	fmt.Println("\t--help       Print help message")
	fmt.Println("\t--version    Print version number")
	fmt.Println(fmt.Sprintf("\n%sExamples:%s", bold, reset))
	fmt.Println("\t$ treegen tree_structure.txt")
	fmt.Println("\t$ cat tree_structure.txt | treegen")
	fmt.Println("\t$ treegen < tree_structure.txt")
	fmt.Println("\t$ treegen <<-EOF")
	fmt.Println("\t    ...")
	fmt.Println("\tEOF\n")
}

func readFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(file)
	defer file.Close() // defer is more idomatic

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func readFromStdIn() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		helpMessage()
		return "", errors.New("")
	}

	content, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func generate(tree string) error {
	lines := strings.Split(tree, "\n")
	level := 0

	dir := lines[0]
	err := createDirectory(dir)
	if err != nil {
		return err
	}

	for i, line := range lines[1:] {
		currentLevel := nodeLevel(line)

		if currentLevel > level {
			parentDir := nodeName(lines[i]) // name of previous line
			dir = dir + parentDir
			level++
		} else if currentLevel < level {
			dir = strings.TrimRight(dir, "/")
			dir = moveUpDirectories(dir, level-currentLevel)
			dir = dir + "/"
			level = currentLevel // could drop down many levels
		}

		nodePath := dir + nodeName(line)

		if strings.HasSuffix(nodePath, "/") {
			err := createDirectory(nodePath)
			if err != nil {
				// we're in a loop here, so we should probably clean up whatever
				// stuff got created if one of these pops
				return err
			}
			continue
		}

		err := createFile(nodePath)
		if err != nil {
			// we're in a loop here, so we should probably clean up whatever
			// stuff got created if one of these pops
			return err
		}
	}

	fmt.Println("Directory structure created successfully")

	return nil
}

func nodeLevel(line string) int {
	level := 0

	for {
		if strings.HasPrefix(line, ascSubChild) {
			level++
			line = line[len(ascSubChild):]
		} else if strings.HasPrefix(line, txtSubChild) {
			level++
			line = line[len(txtSubChild):]
		} else if strings.HasPrefix(line, gap) {
			level++
			line = line[len(gap):]
		} else {
			break
		}
	}

	return level
}

func nodeName(line string) string {
	return strings.TrimLeft(line, "└├ ─│|+-")
}

func moveUpDirectories(path string, n int) string {
	for i := 0; i < n; i++ {
		path = filepath.Dir(path)
	}

	return path
}

func createDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

func createFile(path string) error {
	file, err := os.Create(path)
	defer file.Close()

	return err
}
