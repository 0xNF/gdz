// Package gdzerr contains methods for producing a text file of the errors that this run of gdz encountered
package gdzerr

import (
	"fmt"
	"os"

	"github.com/0xNF/gdz/internal/fs"
)

// Timestamp of the current run
var Timestamp string = ""

// ErrFile is a TXT File to write to if we encounter any errors
var ErrFile *os.File = nil

// GetErrorTxt Produces a new Error file
func GetErrorTxt(timestamp string) (*os.File, error) {
	var txt string = fmt.Sprintf("gdz_errors.%s.txt", timestamp)
	if !fs.Exists(txt) {
		return os.Create(txt)
	}
	return os.OpenFile(txt, os.O_RDWR, 777)
}

// WriteErr writes the contents of msg to the given file
func WriteErr(msg string) error {
	if ErrFile == nil {
		var err error = nil
		ErrFile, err = GetErrorTxt(Timestamp)
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(ErrFile, msg)
	return nil
}

// CloseError Closes the Error file if it exits
func CloseError() {
	if ErrFile != nil {
		ErrFile.Close()
	}
}
