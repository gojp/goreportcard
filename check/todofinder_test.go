package check

import (
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	_, err := readLinesFromFile("testfiles/notexists.go")

	if err == nil {
		t.Errorf("Doesn't throw error!")
	}
}
func TestTodoPercentage(t *testing.T) {
	repoDir := "../_repos/src/github.com/docker"
	filenames, _, _ := GoFiles(repoDir)
	p := TodoFinder{repoDir, filenames}
	todos, _, error := p.Percentage()
	fmt.Println(todos, error)
}
func TestRegex(t *testing.T) {
	returnedTodos, _ := findTodosInFile("testfiles/todo_cases.txt")
	t.Log(returnedTodos)
	currentCases := 6
	if len(returnedTodos) != currentCases {
		t.Errorf("Expected: %v ; Got: %v ", currentCases, len(returnedTodos))
	}
}
