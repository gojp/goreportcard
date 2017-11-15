package check

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

//TodoFinder is th check for the todos in comments
type TodoFinder struct {
	Dir       string
	Filenames []string
}

//Name returns the name of the display name of the command
func (g TodoFinder) Name() string {
	return "todofinder"
}

//Weight returns the weight this check has in the overall average
func (g TodoFinder) Weight() float64 {
	return 0
}

var regex = "(\\/\\/.*\\b(?i)todo(?-i)\\b)"
var todoRegex, err = regexp.Compile(regex)

//Percentage returns files with todo comments
func (g TodoFinder) Percentage() (todosCount float64, fs []FileSummary, err error) {

	for _, file := range g.Filenames {
		todosInFile, _ := findTodosInFile(file)
		todosCount += float64(len(todosInFile))
		if len(todosInFile) != 0 {
			filename := strings.TrimPrefix(file, "_repos/src")
			fs = append(fs, FileSummary{Filename: makeFilename(filename),
				FileURL: fileURL(g.Dir, strings.TrimPrefix(file, "_repos/src")),
				Errors:  todosInFile})

		}
	}
	if err != nil {
		return 0.0, []FileSummary{}, err
	}
	return float64(todosCount), fs, err
}

//Description returns the description of TodoFinder
func (g TodoFinder) Description() string {
	return "Todofinder finds todos in comments"
}

func readLinesFromFile(filename string) ([]string, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0660)
	defer f.Close()

	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(f)
	lines := []string{}
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines, sc.Err()
}

func findTodosInFile(filename string) (todos []Error, err error) {
	fileLines, err := readLinesFromFile(filename)
	if err != nil {
		return nil, err
	}
	for i, line := range fileLines {
		if todoRegex.FindString(line) != "" {
			newTodo := Error{i + 1, line}
			todos = append(todos, newTodo)
		}
	}
	return todos, nil
}
