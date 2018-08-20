package ini

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/knq/ini/parser"
)

// GitSectionManipFunc is a helper method to manipulate sections in ini files
// in a Gitconfig compatible way and provides subsection functionality.
//
// Use it by setting File.SectionManipFunc.
//
// Example:
//		f := ini.LoadString(`...`)
//		f.SectionManipFunc = ini.GitSectionManipFunc
//		f.SectionNameFunc = ini.GitSectionNameFunc
func GitSectionManipFunc(name string) string {
	n, sub := parser.NameSplitFunc(name)
	if n == "" {
		n = sub
		sub = ""
	}

	// clean up name
	n = strings.TrimSpace(strings.ToLower(n))

	// if there's a subsection in the name
	s := ""
	if sub != "" {
		s = fmt.Sprintf(" \"%s\"", sub)
	}
	return n + s
}

// spaceOrTabRE is regexp used for cleaning git section names.
var spaceOrTabRE = regexp.MustCompile(`[ \t]+`)

// GitSectionNameFunc is a helper method to manipulate section names in ini
// files in a Gitconfig compatible way and provides subsection functionality.
//
// Effectively inverse of GitSectionManipFunc.
//
// Use this by setting File.SectionNameFunc.
//
// Example:
//		f := ini.LoadString(`...`)
//		f.SectionManipFunc = ini.GitSectionManipFunc
//		f.SectionNameFunc = ini.GitSectionNameFunc
func GitSectionNameFunc(name string) string {
	// remove " from string
	n := strings.Replace(strings.TrimSpace(name), "\"", "", -1)

	// replace any space or tab with .
	return spaceOrTabRE.ReplaceAllString(n, parser.DefaultNameKeySeparator)
}
