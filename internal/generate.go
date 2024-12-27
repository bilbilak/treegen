package internal

import (
	"path/filepath"
	"strings"
)

const (
	gap         = "    "
	ascChild    = "├── "
	ascSubChild = "│   "
	ascEndChild = "└── "
	txtChild    = "|-- "
	txtSubChild = "|   "
	txtEndChild = "+-- "
)

func generate(lines []string) {
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
			level = currentLevel // could drop down multiple levels
		}

		nodePath := dir + nodeName(line)

		if strings.HasSuffix(nodePath, "/") {
			createDirectory(nodePath)
		} else {
			createFile(nodePath)
		}
	}
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
