package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	version     = "1.0.0"
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

	if flag.NArg() > 0 {
		for _, path := range flag.Args() {
			content := readFromFile(path)
			generate(content)
		}
	} else {
		content := readFromStdIn()
		generate(content)
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

func readFromFile(path string) string {
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("File not found:", path)
		os.Exit(1)
	}

	content, err := io.ReadAll(file)

	if err != nil {
		fmt.Println("File not readable:", path)
		os.Exit(1)
	}

	_ = file.Close()

	return string(content)
}

func readFromStdIn() string {
	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		helpMessage()
		os.Exit(2)
	}

	content, _ := io.ReadAll(os.Stdin)

	return string(content)
}

func generate(tree string) {
	lines := strings.Split(tree, "\n")
	level := 0

	dir := lines[0]
	createDirectory(dir)

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
			createDirectory(nodePath)
		} else {
			createFile(nodePath)
		}
	}

	fmt.Println("Directory structure created successfully")
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

func createDirectory(path string) {
	err := os.MkdirAll(path, 0755)

	if err != nil {
		fmt.Println("Directory not created:", path)
		os.Exit(1)
	}
}

func createFile(path string) {
	file, err := os.Create(path)

	if err != nil {
		fmt.Println("File not created:", path)
		os.Exit(1)
	}

	_ = file.Close()
}
