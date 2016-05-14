package check

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

var (
	skipDirs     = []string{"/Godeps/", "/vendor/", "/third_party/"}
	skipSuffixes = []string{".pb.go", ".pb.gw.go", ".generated.go", "bindata.go"}
)

func addSkipDirs(params []string) []string {
	for _, dir := range skipDirs {
		params = append(params, fmt.Sprintf("--skip=%s", dir))
	}
	return params
}

// GoFiles returns a slice of Go filenames
// in a given directory.
func GoFiles(dir string) ([]string, error) {
	var filenames []string
	visit := func(fp string, fi os.FileInfo, err error) error {
		for _, skip := range skipDirs {
			if strings.Contains(fp, skip) {
				return nil
			}
		}
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if fi.IsDir() {
			return nil // not a file.  ignore.
		}
		fiName := fi.Name()
		for _, skip := range skipSuffixes {
			if strings.HasSuffix(fiName, skip) {
				return nil
			}
		}
		ext := filepath.Ext(fiName)
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

// Error contains the line number and the reason for
// an error output from a command
type Error struct {
	LineNumber  int    `json:"line_number"`
	ErrorString string `json:"error_string"`
}

// FileSummary contains the filename, location of the file
// on GitHub, and all of the errors related to the file
type FileSummary struct {
	Filename string  `json:"filename"`
	FileURL  string  `json:"file_url"`
	Errors   []Error `json:"errors"`
}

// AddError adds an Error to FileSummary
func (fs *FileSummary) AddError(out string) error {
	s := strings.SplitN(out, ":", 2)
	msg := strings.SplitAfterN(s[1], ":", 3)[2]

	e := Error{ErrorString: msg}
	ls := strings.Split(s[1], ":")
	ln, err := strconv.Atoi(ls[0])
	if err != nil {
		return err
	}
	e.LineNumber = ln

	fs.Errors = append(fs.Errors, e)

	return nil
}

// GoTool runs a given go command (for example gofmt, go tool vet)
// on a directory
func GoTool(dir string, filenames, command []string) (float64, []FileSummary, error) {
	params := command[1:]
	params = addSkipDirs(params)
	params = append(params, dir+"/...")

	cmd := exec.Command(command[0], params...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, []FileSummary{}, err
	}

	err = cmd.Start()
	if err != nil {
		return 0, []FileSummary{}, err
	}

	out := bufio.NewScanner(stdout)

	githubLink := strings.TrimPrefix(dir, "repos/src")

	// the same file can appear multiple times out of order
	// in the output, so we can't go line by line, have to store
	// a map of filename to FileSummary
	fsMap := map[string]FileSummary{}
	var failed = []FileSummary{}
outer:
	for out.Scan() {
		filename := strings.Split(out.Text(), ":")[0]
		filename = strings.TrimPrefix(filename, "repos/src")
		for _, skip := range skipSuffixes {
			if strings.HasSuffix(filename, skip) {
				continue outer
			}
		}
		fileURL := "https://" + strings.TrimPrefix(dir, "repos/src/") + "/blob/master" + strings.TrimPrefix(filename, githubLink)
		fs := fsMap[filename]
		if fs.Filename == "" {
			fs.Filename = filename
			if strings.HasPrefix(filename, "/github.com") {
				sp := strings.Split(filename, "/")
				if len(sp) > 3 {
					fs.Filename = strings.Join(sp[3:], "/")
				}

			}
			fs.FileURL = fileURL
		}
		err = fs.AddError(out.Text())
		if err != nil {
			return 0, []FileSummary{}, err
		}
		fsMap[filename] = fs
	}
	if err := out.Err(); err != nil {
		return 0, []FileSummary{}, err
	}

	for _, v := range fsMap {
		failed = append(failed, v)
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
