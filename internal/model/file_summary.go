package model

import (
	"strconv"
	"strings"

	"github.com/yeqown/log"
)

// FileSummary contains the filename, location of the file
// on GitHub, and all of the errors related to the file
type FileSummary struct {
	Filename string  `json:"filename"`
	FileURL  string  `json:"file_url"`
	Errors   []Error `json:"errors"`
}

// AddError adds an Error to FileSummary
// out format like this: handlers/linter.go:68:3: ineffectual assignment to `err` (ineffassign)
// if `out` not in format, just output Error without any deal
func (fs *FileSummary) AddError(out string) error {
	s := strings.Split(out, ":")
	// log.Infof("out=%s, s=%v\n", out, s)
	// msg := strings.SplitAfterN(s[1], ":", 3)[2]

	if len(s) != 4 {
		log.Infof("out=%s, len(s)=%d", out, len(s))
		// return errors.New("invalid error output format")
		return nil
	}

	// e := Error{ErrorString: msg}
	// ls := strings.Split(s[1], ":")
	// ln, err := strconv.Atoi(ls[0])
	// if err != nil {
	// 	return fmt.Errorf("AddError: could not parse %q - %v", out, err)
	// }
	// e.LineNumber = ln

	errmsg := s[3]
	lineNo, _ := strconv.Atoi(s[1])

	fs.Errors = append(fs.Errors, Error{
		LineNumber:  lineNo,
		ErrorString: errmsg,
	})

	return nil
}
