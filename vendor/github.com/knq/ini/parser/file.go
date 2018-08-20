package parser

import (
	"bytes"
	"fmt"
	"os"
)

// File represents parsed ini data.
type File struct {
	// lines in file.
	lines []*Line

	// sections in file.
	sections []*Section

	// Manipulation function used on section name for AddSection,
	// RenameSection.
	SectionManipFunc func(string) string

	// Function used to normalize and format section name for presentation.
	SectionNameFunc func(string) string

	// Comparison function used to find section in File. Set this to override
	// default comparison behavior.
	SectionCompFunc func(string, string) bool

	// Manipulation function used on key in File.
	KeyManipFunc func(string) string

	// Comparison function used to find key in File.
	KeyCompFunc func(string, string) bool

	// Manipulation function used when setting value in File.
	ValueManipFunc func(string) string

	// Function is used to split a key name (such as section.key).
	NameSplitFunc func(string) (string, string)
}

// NewFile creates a new ini.File from provided lines.
func NewFile(lines []*Line) *File {
	// create
	ret := &File{
		lines:    lines,
		sections: make([]*Section, 0),

		// copy manipulation funcs
		SectionManipFunc: SectionManipFunc,
		SectionNameFunc:  SectionNameFunc,
		SectionCompFunc:  nil,
		KeyManipFunc:     KeyManipFunc,
		KeyCompFunc:      KeyCompFunc,
		ValueManipFunc:   ValueManipFunc,
		NameSplitFunc:    NameSplitFunc,
	}

	// create default section
	lastSection := NewSection(position{}, "", "", nil)
	lastSection.file = ret
	ret.sections = append(ret.sections, lastSection)

	// loop over lines and build sections/keys
	for _, l := range lines {
		switch l.item.(type) {
		case *Section:
			// get section
			lastSection, _ = l.item.(*Section)

			// save data
			lastSection.file = ret
			ret.sections = append(ret.sections, lastSection)

		case *KeyValuePair:
			// save kvp
			kvp, _ := l.item.(*KeyValuePair)
			lastSection.keys = append(lastSection.keys, kvp.key)
		}
	}

	return ret
}

// String returns formatted ini file data.
//
// Satisfies fmt.Stringer interface.
func (f File) String() string {
	var buf bytes.Buffer
	for _, l := range f.lines {
		buf.WriteString(l.String())
	}
	return buf.String()
}

// Write to filename.
func (f *File) Write(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(f.String())
	return err
}

// LineCount returns the line count.
func (f *File) LineCount() int {
	return len(f.lines)
}

// RawSectionNames returns all section names from File in a raw (unmanipulated)
// format.
func (f *File) RawSectionNames() []string {
	names := make([]string, len(f.sections))
	for i, s := range f.sections {
		names[i] = s.RawName()
	}
	return names
}

// SectionNames returns all section names from File.
//
// Section names are passed through SectionNameFunc.
func (f *File) SectionNames() []string {
	names := make([]string, len(f.sections))
	for i, s := range f.sections {
		names[i] = s.Name()
	}
	return names
}

// AddSectionRaw adds a Section to File with a raw (unmanipulated) name.
//
// Returns the created Section.
func (f *File) AddSectionRaw(name string) *Section {
	// if its "", then avoid retrieving ...
	if f.sectionNameComp(name, "") {
		return f.GetSection("")
	}

	// create section
	s := NewSection(position{}, name, "", nil)
	s.file = f

	// add section data to file
	f.sections = append(f.sections, s)

	if len(f.lines) > 0 && f.lines[len(f.lines)-1].item == nil {
		// if it's a blank line on the last line, then put it there
		f.lines[len(f.lines)-1].item = s
	} else {
		// default line ending
		le := DefaultLineEnding
		if len(f.lines) > 0 {
			// take line ending from first line if present
			le = f.lines[0].le
		}

		// create the line and append to end
		l := NewLine(position{}, "", s, le)
		f.lines = append(f.lines, l)
	}

	return s
}

// AddSection adds a Section to File.
//
// Section name is passed through file's SectionManipFunc.
//
// Returns the created Section.
func (f *File) AddSection(name string) *Section {
	return f.AddSectionRaw(f.SectionManipFunc(name))
}

// sectionNameComp compares provided Section names to determine if they are
// equal.
//
// Uses f.SectionCompFunc if present, otherwise compares result of
// SectionNameFunc(a) == SectionNameFunc(b).
func (f *File) sectionNameComp(a, b string) bool {
	if f.SectionCompFunc != nil {
		return f.SectionCompFunc(a, b)
	}

	return f.SectionNameFunc(a) == f.SectionNameFunc(b)
}

// getSection Get a section and its starting line number.
func (f *File) getSection(name string) (*Section, int) {
	n := f.SectionManipFunc(name)

	// blank section isn't actually defined ...
	if f.sectionNameComp(n, "") {
		return f.sections[0], 0
	}

	// loop through lines and find section
	for idx, line := range f.lines {
		if s, ok := line.item.(*Section); ok && f.sectionNameComp(n, s.name) {
			return s, idx
		}
	}

	return nil, -1
}

// GetSection returns a Section with provided name from File.
func (f *File) GetSection(name string) *Section {
	s, _ := f.getSection(name)
	return s
}

// SetMap sets all section and key values from provided map.
//
// Replaces values if the key already exists, or adds them otherwise.
func (f *File) SetMap(values map[string]map[string]string) {
	for name, keys := range values {
		section := f.GetSection(name)
		if section == nil {
			section = f.AddSection(name)
		}

		for k, v := range keys {
			section.SetKey(k, v)
		}
	}
}

// GetMap returns all sections and key values as map.
func (f *File) GetMap() map[string]map[string]string {
	ret := make(map[string]map[string]string)

	for _, section := range f.sections {
		s := make(map[string]string)
		for _, key := range section.keys {
			s[f.KeyManipFunc(key)] = f.ValueManipFunc(section.GetRaw(key))
		}

		ret[section.Name()] = s
	}

	return ret
}

// SetMapFlat sets section and key values from a flat map.
func (f *File) SetMapFlat(values map[string]string) {
	for key, value := range values {
		f.SetKey(key, value)
	}
}

// GetMapFlat retrieves all sections and keys and values as flat map.
func (f *File) GetMapFlat() map[string]string {
	ret := make(map[string]string)

	for _, section := range f.sections {
		name := section.Name()
		if section.name != "" {
			name = fmt.Sprintf("%s%s", name, DefaultNameKeySeparator)
		}

		for _, key := range section.keys {
			ret[fmt.Sprintf("%s%s", name, key)] = f.ValueManipFunc(section.GetRaw(key))
		}
	}

	return ret
}

// RenameSectionRaw renames a Section in File using raw (unmanipulated) names.
func (f *File) RenameSectionRaw(name, value string) {
	s := f.GetSection(name)
	s.name = value
}

// RenameSection renames a Section in File.
//
// Value will be passed through the File's SectionManipFunc.
func (f *File) RenameSection(name, value string) {
	f.RenameSectionRaw(name, f.SectionManipFunc(value))
}

// RemoveSection removes a Section and all related lines from File.
func (f *File) RemoveSection(name string) {
	section, start := f.getSection(name)
	if section == nil {
		return
	}

	// save copy of line ending
	le := f.lines[0].le

	// find next section
	end := start + 1
	for ; end < len(f.lines); end++ {
		if _, ok := f.lines[end].item.(*Section); ok {
			break
		}
	}

	// remove from f.lines
	f.lines = append(f.lines[:start], f.lines[end:]...)

	// if we removed all lines, then put a blank line back in
	if len(f.lines) < 1 {
		line := NewLine(position{}, "", nil, le)
		f.lines = []*Line{line}
	}

	// find in f.sections
	pos := -1
	for idx, s := range f.sections {
		if section == s {
			pos = idx
			break
		}
	}

	// remove from f.sections
	f.sections = append(f.sections[:pos], f.sections[pos+1:]...)
}

// SetKey sets a key's value in File with name in form of section.key.
//
// If no section is specified, then the empty (first) section is used.
//
// Uses File's NameSplitFunc to split the key.
func (f *File) SetKey(key, value string) {
	name, k := f.NameSplitFunc(key)

	// get the section
	section := f.GetSection(name)
	if section == nil {
		section = f.AddSection(name)
	}

	section.SetKey(k, value)
}

// GetKey retrieves a stored key's value from File with name in form of
// section.key.
//
// If no section is specified, then the key will be looked for in the empty
// (first) section.
//
// Uses File's NameSplitFunc to split the key.
func (f *File) GetKey(key string) string {
	name, k := f.NameSplitFunc(key)

	// get the section
	section := f.GetSection(name)
	if section == nil {
		return ""
	}

	return section.Get(k)
}

// RemoveKey removes a key and its value from file with name in form of
// section.key.
//
// If no section is specified, then the empty (first) section is used.
//
// Uses File's NameSplitFunc to split the key.
func (f *File) RemoveKey(key string) {
	name, k := f.NameSplitFunc(key)

	// get the section
	section := f.GetSection(name)
	if section == nil {
		return
	}

	section.RemoveKey(k)
}
