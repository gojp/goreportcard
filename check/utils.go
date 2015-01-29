package check

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
)

// GoFiles returns a slice of Go filenames
// in a given directory.
func GoFiles(dir string) ([]string, error) {
	var filenames []string
	visit := func(fp string, fi os.FileInfo, err error) error {
		if strings.Contains(fp, "Godeps") {
			return nil
		}
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // not a file.  ignore.
		}
		ext := filepath.Ext(fi.Name())
		if ext == ".go" {
			filenames = append(filenames, fp)
		}
		return nil
	}

	err := filepath.Walk(dir, visit)

	return filenames, err
}

// lineCount returns the number of lines in a given file
func lineCount(filepath string) (int, error) {
	out, err := exec.Command("wc", "-l", filepath).Output()
	if err != nil {
		return 0, err
	}
	// wc output is like: 999 filename.go
	count, err := strconv.Atoi(strings.Split(strings.TrimSpace(string(out)), " ")[0])
	if err != nil {
		return 0, err
	}

	return count, nil
}

type Error struct {
	LineNumber  int    `json:"line_number"`
	ErrorString string `json:"error_string"`
}

type FileSummary struct {
	Filename string  `json:"filename"`
	FileURL  string  `json:"file_url"`
	Errors   []Error `json:"errors"`
}

// ByFilename implements sort.Interface for []Person based on
// the Age field.
type ByFilename []FileSummary

func (a ByFilename) Len() int           { return len(a) }
func (a ByFilename) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFilename) Less(i, j int) bool { return a[i].Filename < a[j].Filename }

func getFileSummary(filename, dir, cmd, out string) (FileSummary, error) {
	filename = strings.TrimPrefix(filename, "repos/src")
	githubLink := strings.TrimPrefix(dir, "repos/src")
	fileURL := "https://" + strings.TrimPrefix(dir, "repos/src/") + "/blob/master" + strings.TrimPrefix(filename, githubLink)
	fs := FileSummary{
		Filename: filename,
		FileURL:  fileURL,
	}
	split := strings.Split(string(out), "\n")
	for _, sp := range split[0 : len(split)-1] {
		parts := strings.Split(sp, ":")
		msg := sp
		if cmd != "gocyclo" {
			msg = parts[len(parts)-1]
		}
		e := Error{ErrorString: msg}
		switch cmd {
		case "golint", "gocyclo", "vet":
			ln, err := strconv.Atoi(strings.Split(sp, ":")[1])
			if err != nil {
				return fs, err
			}
			e.LineNumber = ln
		}

		fs.Errors = append(fs.Errors, e)
	}

	return fs, nil
}

// GoTool runs a given go command (for example gofmt, go tool vet)
// on a directory
func GoTool(dir string, filenames, command []string) (float64, []FileSummary, error) {
	var failed = []FileSummary{}
	for _, fi := range filenames {
		params := command[1:]
		params = append(params, fi)

		cmd := exec.Command(command[0], params...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return 0, []FileSummary{}, nil
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return 0, []FileSummary{}, nil
		}

		err = cmd.Start()
		if err != nil {
			return 0, []FileSummary{}, nil
		}

		out, err := ioutil.ReadAll(stdout)
		if err != nil {
			return 0, []FileSummary{}, nil
		}

		errout, err := ioutil.ReadAll(stderr)
		if err != nil {
			return 0, []FileSummary{}, nil
		}

		if string(out) != "" {
			fs, err := getFileSummary(fi, dir, command[0], string(out))
			if err != nil {
				return 0, []FileSummary{}, nil
			}
			failed = append(failed, fs)
		}

		// go vet logs to stderr
		if string(errout) != "" {
			cmd := command[0]
			if reflect.DeepEqual(command, []string{"go", "tool", "vet"}) {
				cmd = "vet"
			}
			fs, err := getFileSummary(fi, dir, cmd, string(errout))
			if err != nil {
				return 0, []FileSummary{}, nil
			}
			failed = append(failed, fs)
		}

		err = cmd.Wait()
		if exitErr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				// some commands exit 1 when files fail to pass (for example go vet)
				if status.ExitStatus() != 1 {
					return 0, failed, err
					// return 0, Error{}, err
				}
			}
		}

	}

	if len(filenames) == 1 {
		lc, err := lineCount(filenames[0])
		if err != nil {
			return 0, failed, err
		}

		var errors int
		if len(failed) != 0 {
			errors = len(failed[0].Errors)
		}

		return float64(lc-errors) / float64(lc), failed, nil
	}

	return float64(len(filenames)-len(failed)) / float64(len(filenames)), failed, nil
}
