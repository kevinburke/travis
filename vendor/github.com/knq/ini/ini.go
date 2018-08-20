// Package ini provides a simple package to read/write/manipulate ini files.
//
// Mainly a frontend to http://github.com/knq/ini/parser
package ini

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/knq/ini/parser"
)

// Error is a ini error.
type Error string

// Error satisfies the error interface.
func (err Error) Error() string {
	return string(err)
}

// Error values.
const (
	ErrNoFilenameSupplied Error = "no filename supplied"
)

// ParseError is a ini parse error.
type ParseError struct {
	name string
	err  error
}

// Error satisfies the error interface.
func (err *ParseError) Error() string {
	return fmt.Sprintf("unable to parse %s: %v", err.name, err.err)
}

// File wraps parser.File with information about an ini file.
//
// File can be written to disk by calling File.Save.
type File struct {
	*parser.File        // ini file
	Filename     string // filename to read/write from/to
}

// NewFile creates a new File.
func NewFile() *File {
	return &File{
		File:     parser.NewFile(nil),
		Filename: "",
	}
}

// Save writes the ini file data to File.Filename.
//
// Returns error if File.Filename name was not set, or if an error was
// encountered during write. Simple wrapper around parser.File.Write.
func (f *File) Save() error {
	if f.Filename == "" {
		return ErrNoFilenameSupplied
	}
	return f.Write(f.Filename)
}

// Parse passes the filename/reader to ini.Parser.Parse.
func Parse(name, filename string, r io.Reader) (*File, error) {
	// sanitize data first (ensure file ends with '\n')
	buf, err := fixEnding(r)
	if err != nil {
		return nil, err
	}

	// pass through ini/parser package
	f, err := parser.Parse(name, buf)
	if err != nil {
		return nil, &ParseError{name, parser.LastError()}
	}

	// convert to *parser.File
	inifile, ok := f.(*parser.File)
	if !ok {
		return nil, &ParseError{name, parser.LastError()}
	}

	return &File{
		File:     inifile,
		Filename: filename,
	}, nil
}

// Load loads ini file from a io.Reader.
func Load(r io.Reader) (*File, error) {
	return Parse("<io.Reader>", "", r)
}

// LoadBytes loads ini file from a byte slice.
func LoadBytes(buf []byte) (*File, error) {
	return Parse("<buffer>", "", bytes.NewReader(buf))
}

// LoadString loads ini file from string.
func LoadString(str string) (*File, error) {
	return Parse("<string>", "", strings.NewReader(str))
}

// LoadFile loads ini data from a file with specified filename.
//
// If the filename doesn't exist, then an empty File is returned. The data can
// then be written to disk using File.Save, or parser.File.Write.
func LoadFile(filename string) (*File, error) {
	// check if the file exists, return a new file if it doesn't
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file := NewFile()
		file.Filename = filename
		return file, nil
	}

	// if file exists, read and parse it
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(filename, filename, f)
}

// fixEnding fixes the file data in r, ensuring the file ends with \n.
func fixEnding(r io.Reader) ([]byte, error) {
	// read
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// add '\n' to end if not present
	if len(buf) == 0 || buf[len(buf)-1] != '\n' {
		return append(buf, '\n'), nil
	}
	return buf, nil
}
