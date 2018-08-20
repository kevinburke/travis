// Package parser is a Pigeon-based ini file parser.
//
// Please see http://godoc.org/github.com/knq/ini for the frontend package.
package parser

//go:generate ./gen.sh

import (
	"fmt"
	"strings"
)

var (
	// DefaultLeadingKeyWhitespace is the default leading whitespace for
	// non-blank section keys.
	DefaultLeadingKeyWhitespace = "\t"

	// DefaultLineEnding is the default line ending for new lines added to
	// file.
	DefaultLineEnding = "\n"

	// DefaultNameKeySeparator is the default separator token for section.name
	// style keys.
	DefaultNameKeySeparator = "."

	// last position
	lastPosition position

	// last text
	lastText string
)

// SectionManipFunc manipulates a Section name.
//
// This function is used when a section name is created or altered.
//
// Override on a per-File basis by setting File.SectionManipFunc.
func SectionManipFunc(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// SectionNameFunc formats a Section name.
//
// This function is used to format (normalize) Section names.
//
// Override on a per-File basis by setting File.SectionNameFunc.
func SectionNameFunc(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// KeyManipFunc manipulates key names.
//
// Takes a key name and returns the value that used. By default does
// strings.TrimSpace(strings.ToLower(key)).
//
// This function is used when a key is created or altered.
//
// Override on a per-File basis by setting File.KeyManipFunc.
func KeyManipFunc(key string) string {
	return strings.TrimSpace(strings.ToLower(key))
}

// KeyCompFunc is used to compare key names on get/set.
//
// Passes keys a, b through KeyManipFunc and returns string equality.
//
// This function is used when key names are compared.
//
// Override on a per-File basis by setting File.KeyCompFunc.
func KeyCompFunc(a, b string) bool {
	return KeyManipFunc(a) == KeyManipFunc(b)
}

// NameSplitFunc splits Section names.
//
// Splits names based on DefaultNameKeySeparator.
//
// Returns section, key.
//
// This function is used to split keys when being retrieved or set on a File.
//
// Override on a per-File basis by setting File.NameSplitFunc.
func NameSplitFunc(name string) (string, string) {
	idx := strings.LastIndex(name, DefaultNameKeySeparator)

	// no section name
	if idx < 0 {
		return "", name
	}

	return name[:idx], name[idx+1:]
}

// ValueManipFunc manipulates values when setting a key value.
//
// Override on a per-File basis by setting File.ValueManipFunc.
func ValueManipFunc(value string) string {
	return strings.TrimSpace(value)
}

// LastError returns the last error encountered (if any) during parsing.
func LastError() error {
	return fmt.Errorf("error on line %d:%d near '%s'", lastPosition.line, lastPosition.col, lastText)
}

// Item is the shared interface for Comment, Section, and KeyValuePair.
type Item interface {
	String() string
}

// Line in a File.
type Line struct {
	pos position

	ws   string // Leading whitespace (if any)
	item Item   // A Comment, Section, or KeyValuePair
	le   string
}

// NewLine creates a new line.
func NewLine(pos position, ws string, item Item, le string) *Line {
	return &Line{
		pos: pos,

		ws:   ws,
		item: item,
		le:   le,
	}
}

// String returns a formatted line.
func (l Line) String() string {
	item := ""
	if l.item != nil {
		item = l.item.String()
	}

	return fmt.Sprintf("%s%s%s", l.ws, item, l.le)
}

// Comment in a File.
type Comment struct {
	pos position

	cs      string // comment separator
	comment string // actual comment
}

// NewComment creates a new Comment.
func NewComment(pos position, cs string, comment string) *Comment {
	return &Comment{
		pos: pos,

		cs:      cs,
		comment: comment,
	}
}

// String returns a formatted comment.
func (c Comment) String() string {
	return fmt.Sprintf("%s%s", c.cs, c.comment)
}

// KeyValuePair in a File.
type KeyValuePair struct {
	//section *Section

	pos position

	key   string
	ws    string
	value string

	comment *Comment
}

// NewKeyValuePair creates a new key value pair.
func NewKeyValuePair(pos position, key, ws, value string, comment *Comment) *KeyValuePair {
	return &KeyValuePair{
		pos: pos,

		key:     key,
		ws:      ws,
		value:   value,
		comment: comment,
	}
}

// String returns a formatted key value pair.
func (kvp KeyValuePair) String() string {
	comment := ""
	if kvp.comment != nil {
		comment = kvp.comment.String()
	}

	return fmt.Sprintf("%s=%s%s%s", kvp.key, kvp.ws, kvp.value, comment)
}
