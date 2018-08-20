package parser

import "fmt"

// Section in a File.
type Section struct {
	file *File

	pos position

	name    string
	ws      string
	comment *Comment

	keys []string
}

// NewSection creates a new section.
func NewSection(pos position, name, ws string, comment *Comment) *Section {
	var keys []string

	return &Section{
		pos: pos,

		name:    name,
		ws:      ws,
		comment: comment,

		keys: keys,
	}
}

// String returns a formatted section.
func (s Section) String() string {
	comment := ""
	if s.comment != nil {
		comment = s.comment.String()
	}

	return fmt.Sprintf("[%s]%s%s", s.name, s.ws, comment)
}

// RawName returns the raw (unmanipulated) Section name.
func (s *Section) RawName() string {
	return s.name
}

// Name returns the Section name.
//
// Pasess name through File's SectionNameFunc.
func (s *Section) Name() string {
	return s.file.SectionNameFunc(s.name)
}

// RawKeys returns the raw (unmanipulated) keys defined in Section.
func (s *Section) RawKeys() []string {
	return s.keys
}

// Keys returns the keys defined in Section.
//
// Keys are passed through File's KeyManipFunc.
func (s *Section) Keys() []string {
	keys := make([]string, len(s.keys))
	for i, k := range s.keys {
		keys[i] = s.file.KeyManipFunc(k)
	}

	return keys
}

// getInsertLocation determines insert location in a Section, which is the
// first blank line after a non-blank.
func (s *Section) getInsertLocation(idx int) int {
	for i := idx; i >= 0; i-- {
		if s.file.lines[i].item != nil {
			return i + 1
		}
	}

	return -1
}

// getKey returns the KeyValuePair and its line position, or nil and the
// position the key should be inserted at.
func (s *Section) getKey(key string) (*KeyValuePair, int) {
	// loop over lines and find the key
	lastSectionName := ""
	for lastIdx, l := range s.file.lines {
		switch l.item.(type) {
		case *Section:
			if lastSectionName == s.name {
				// must be entering a new section; so not found, return
				return nil, s.getInsertLocation(lastIdx - 1)
			}

			sect, _ := l.item.(*Section)
			lastSectionName = sect.name

		case *KeyValuePair:
			kvp, _ := l.item.(*KeyValuePair)
			//fmt.Printf(">>> compare: %s//%s :: %s//%s\n", lastSectionName, s.name, kvp.key, key)
			if lastSectionName == s.name && s.file.KeyCompFunc(kvp.key, key) {
				return kvp, lastIdx
			}
		}
	}

	// if we get here, then must be last section of file
	return nil, s.getInsertLocation(len(s.file.lines) - 1)
}

// GetRaw returns the raw (unmanipulated) value for a key.
func (s *Section) GetRaw(key string) string {
	k, _ := s.getKey(key)
	if k != nil {
		return k.value
	}

	return ""
}

// Get returns the value for a key.
//
// The value is passed through ValueManipFunc.
func (s *Section) Get(key string) string {
	return s.file.ValueManipFunc(s.GetRaw(key))
}

// SetKeyValueRaw sets a key's value to the raw (unmanipulated) value.
//
// If key already present, then it's value is overwritten. If key doesn't
// exist, then it is added to the end of the Section.
func (s *Section) SetKeyValueRaw(key, value string) {
	// get position
	k, pos := s.getKey(key)

	// key is present, set value
	if k != nil {
		k.value = value
		return
	}

	// key doesn't exist, create it...

	// grab default whitespace and line ending
	ws := DefaultLeadingKeyWhitespace

	// set no ws if empty section
	if s.name == "" {
		ws = ""
	}

	le := DefaultLineEnding
	if len(s.file.lines) > 0 {
		// take line ending from first line if present
		le = s.file.lines[0].le
	}

	// create the key and line
	k = NewKeyValuePair(position{}, key, "", value, nil)
	line := NewLine(position{}, ws, k, le)

	// insert line into s.file.lines
	if pos < 0 {
		// must be inserting into empty section where there are no keys present
		s.file.lines = append([]*Line{line}, s.file.lines...)
	} else {
		// copy whitespace from previous line if its a kvp
		if _, ok := s.file.lines[pos-1].item.(*KeyValuePair); ok {
			line.ws = s.file.lines[pos-1].ws
		}

		s.file.lines = append(
			s.file.lines[:pos],
			append(
				[]*Line{line},
				s.file.lines[pos:]...,
			)...,
		)
	}

	// add key to s.keys
	s.keys = append(s.keys, k.key)
}

// SetKey sets a key to the provided value.
//
// If key already present, then it's value is overwritten. If key doesn't
// exist, then it is added to the end of the Section.
//
// Passes key through KeyManipFunc and value through ValueManipFunc.
func (s *Section) SetKey(key, value string) {
	s.SetKeyValueRaw(s.file.KeyManipFunc(key), s.file.ValueManipFunc(value))
}

// RemoveKey removes a key and its value from Section.
//
// If there is a comment on the line, it will not be removed.
func (s *Section) RemoveKey(key string) {
	k, pos := s.getKey(key)
	if k != nil {
		s.file.lines = append(s.file.lines[:pos], s.file.lines[pos+1:]...)

		// find place in s.keys
		idx := 0
		for ; idx < len(s.keys); idx++ {
			if s.file.KeyCompFunc(key, s.keys[idx]) {
				break
			}
		}

		// remove from s.keys
		s.keys = append(s.keys[:idx], s.keys[idx+1:]...)
	}
}
