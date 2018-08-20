package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Helper function taken from pigeon source / examples
func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}

	return v.([]interface{})
}

var g = &grammar{
	rules: []*rule{
		{
			name: "File",
			pos:  position{line: 17, col: 1, offset: 236},
			expr: &actionExpr{
				pos: position{line: 17, col: 9, offset: 244},
				run: (*parser).callonFile1,
				expr: &seqExpr{
					pos: position{line: 17, col: 9, offset: 244},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 17, col: 9, offset: 244},
							label: "lines",
							expr: &zeroOrMoreExpr{
								pos: position{line: 17, col: 15, offset: 250},
								expr: &ruleRefExpr{
									pos:  position{line: 17, col: 15, offset: 250},
									name: "Line",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 17, col: 21, offset: 256},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "Line",
			pos:  position{line: 33, col: 1, offset: 589},
			expr: &actionExpr{
				pos: position{line: 33, col: 9, offset: 597},
				run: (*parser).callonLine1,
				expr: &seqExpr{
					pos: position{line: 33, col: 9, offset: 597},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 33, col: 9, offset: 597},
							label: "ws",
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 12, offset: 600},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 33, col: 14, offset: 602},
							label: "item",
							expr: &zeroOrOneExpr{
								pos: position{line: 33, col: 19, offset: 607},
								expr: &choiceExpr{
									pos: position{line: 33, col: 20, offset: 608},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 33, col: 20, offset: 608},
											name: "Comment",
										},
										&ruleRefExpr{
											pos:  position{line: 33, col: 30, offset: 618},
											name: "Section",
										},
										&ruleRefExpr{
											pos:  position{line: 33, col: 40, offset: 628},
											name: "KeyValuePair",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 33, col: 55, offset: 643},
							label: "le",
							expr: &ruleRefExpr{
								pos:  position{line: 33, col: 58, offset: 646},
								name: "LineEnd",
							},
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 42, col: 1, offset: 864},
			expr: &actionExpr{
				pos: position{line: 42, col: 12, offset: 875},
				run: (*parser).callonComment1,
				expr: &seqExpr{
					pos: position{line: 42, col: 12, offset: 875},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 42, col: 12, offset: 875},
							label: "cs",
							expr: &choiceExpr{
								pos: position{line: 42, col: 16, offset: 879},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 42, col: 16, offset: 879},
										val:        ";",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 42, col: 22, offset: 885},
										val:        "#",
										ignoreCase: false,
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 42, col: 27, offset: 890},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 42, col: 35, offset: 898},
								name: "CommentVal",
							},
						},
					},
				},
			},
		},
		{
			name: "Section",
			pos:  position{line: 50, col: 1, offset: 1111},
			expr: &actionExpr{
				pos: position{line: 50, col: 12, offset: 1122},
				run: (*parser).callonSection1,
				expr: &seqExpr{
					pos: position{line: 50, col: 12, offset: 1122},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 50, col: 12, offset: 1122},
							val:        "[",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 50, col: 16, offset: 1126},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 50, col: 21, offset: 1131},
								name: "SectionName",
							},
						},
						&litMatcher{
							pos:        position{line: 50, col: 33, offset: 1143},
							val:        "]",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 50, col: 37, offset: 1147},
							label: "ws",
							expr: &ruleRefExpr{
								pos:  position{line: 50, col: 40, offset: 1150},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 50, col: 42, offset: 1152},
							label: "comment",
							expr: &zeroOrOneExpr{
								pos: position{line: 50, col: 50, offset: 1160},
								expr: &ruleRefExpr{
									pos:  position{line: 50, col: 50, offset: 1160},
									name: "Comment",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "KeyValuePair",
			pos:  position{line: 59, col: 1, offset: 1388},
			expr: &actionExpr{
				pos: position{line: 59, col: 17, offset: 1404},
				run: (*parser).callonKeyValuePair1,
				expr: &seqExpr{
					pos: position{line: 59, col: 17, offset: 1404},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 59, col: 17, offset: 1404},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 59, col: 21, offset: 1408},
								name: "Key",
							},
						},
						&litMatcher{
							pos:        position{line: 59, col: 25, offset: 1412},
							val:        "=",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 59, col: 29, offset: 1416},
							label: "ws",
							expr: &ruleRefExpr{
								pos:  position{line: 59, col: 32, offset: 1419},
								name: "_",
							},
						},
						&labeledExpr{
							pos:   position{line: 59, col: 34, offset: 1421},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 59, col: 38, offset: 1425},
								name: "Value",
							},
						},
						&labeledExpr{
							pos:   position{line: 59, col: 44, offset: 1431},
							label: "comment",
							expr: &zeroOrOneExpr{
								pos: position{line: 59, col: 52, offset: 1439},
								expr: &ruleRefExpr{
									pos:  position{line: 59, col: 52, offset: 1439},
									name: "Comment",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CommentVal",
			pos:  position{line: 68, col: 1, offset: 1700},
			expr: &actionExpr{
				pos: position{line: 68, col: 15, offset: 1714},
				run: (*parser).callonCommentVal1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 68, col: 15, offset: 1714},
					expr: &seqExpr{
						pos: position{line: 68, col: 16, offset: 1715},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 68, col: 16, offset: 1715},
								expr: &ruleRefExpr{
									pos:  position{line: 68, col: 17, offset: 1716},
									name: "LineEnd",
								},
							},
							&anyMatcher{
								line: 68, col: 25, offset: 1724,
							},
						},
					},
				},
			},
		},
		{
			name: "SectionName",
			pos:  position{line: 76, col: 1, offset: 1891},
			expr: &actionExpr{
				pos: position{line: 76, col: 16, offset: 1906},
				run: (*parser).callonSectionName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 76, col: 16, offset: 1906},
					expr: &charClassMatcher{
						pos:        position{line: 76, col: 16, offset: 1906},
						val:        "[^#;\\r\\n[\\]]",
						chars:      []rune{'#', ';', '\r', '\n', '[', ']'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "Key",
			pos:  position{line: 84, col: 1, offset: 2084},
			expr: &actionExpr{
				pos: position{line: 84, col: 8, offset: 2091},
				run: (*parser).callonKey1,
				expr: &oneOrMoreExpr{
					pos: position{line: 84, col: 8, offset: 2091},
					expr: &charClassMatcher{
						pos:        position{line: 84, col: 8, offset: 2091},
						val:        "[^#;=\\r\\n[\\]]",
						chars:      []rune{'#', ';', '=', '\r', '\n', '[', ']'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "Value",
			pos:  position{line: 92, col: 1, offset: 2262},
			expr: &choiceExpr{
				pos: position{line: 92, col: 10, offset: 2271},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 92, col: 10, offset: 2271},
						name: "QuotedValue",
					},
					&actionExpr{
						pos: position{line: 92, col: 24, offset: 2285},
						run: (*parser).callonValue3,
						expr: &ruleRefExpr{
							pos:  position{line: 92, col: 24, offset: 2285},
							name: "SimpleValue",
						},
					},
				},
			},
		},
		{
			name: "QuotedValue",
			pos:  position{line: 100, col: 1, offset: 2455},
			expr: &actionExpr{
				pos: position{line: 100, col: 16, offset: 2470},
				run: (*parser).callonQuotedValue1,
				expr: &seqExpr{
					pos: position{line: 100, col: 16, offset: 2470},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 100, col: 16, offset: 2470},
							val:        "\"",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 100, col: 20, offset: 2474},
							expr: &ruleRefExpr{
								pos:  position{line: 100, col: 20, offset: 2474},
								name: "Char",
							},
						},
						&litMatcher{
							pos:        position{line: 100, col: 26, offset: 2480},
							val:        "\"",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 100, col: 30, offset: 2484},
							name: "_",
						},
					},
				},
			},
		},
		{
			name: "Char",
			pos:  position{line: 108, col: 1, offset: 2650},
			expr: &choiceExpr{
				pos: position{line: 108, col: 9, offset: 2658},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 108, col: 9, offset: 2658},
						exprs: []interface{}{
							&notExpr{
								pos: position{line: 108, col: 9, offset: 2658},
								expr: &choiceExpr{
									pos: position{line: 108, col: 11, offset: 2660},
									alternatives: []interface{}{
										&litMatcher{
											pos:        position{line: 108, col: 11, offset: 2660},
											val:        "\"",
											ignoreCase: false,
										},
										&litMatcher{
											pos:        position{line: 108, col: 17, offset: 2666},
											val:        "\\",
											ignoreCase: false,
										},
									},
								},
							},
							&anyMatcher{
								line: 108, col: 23, offset: 2672,
							},
						},
					},
					&actionExpr{
						pos: position{line: 108, col: 27, offset: 2676},
						run: (*parser).callonChar8,
						expr: &seqExpr{
							pos: position{line: 108, col: 27, offset: 2676},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 108, col: 27, offset: 2676},
									val:        "\\",
									ignoreCase: false,
								},
								&choiceExpr{
									pos: position{line: 108, col: 33, offset: 2682},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 108, col: 33, offset: 2682},
											val:        "[\\\\/bfnrt\"]",
											chars:      []rune{'\\', '/', 'b', 'f', 'n', 'r', 't', '"'},
											ignoreCase: false,
											inverted:   false,
										},
										&seqExpr{
											pos: position{line: 108, col: 47, offset: 2696},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 108, col: 47, offset: 2696},
													val:        "u",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 108, col: 51, offset: 2700},
													name: "HexDigit",
												},
												&ruleRefExpr{
													pos:  position{line: 108, col: 60, offset: 2709},
													name: "HexDigit",
												},
												&ruleRefExpr{
													pos:  position{line: 108, col: 69, offset: 2718},
													name: "HexDigit",
												},
												&ruleRefExpr{
													pos:  position{line: 108, col: 78, offset: 2727},
													name: "HexDigit",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 116, col: 1, offset: 2909},
			expr: &actionExpr{
				pos: position{line: 116, col: 13, offset: 2921},
				run: (*parser).callonHexDigit1,
				expr: &charClassMatcher{
					pos:        position{line: 116, col: 13, offset: 2921},
					val:        "[0-9a-f]i",
					ranges:     []rune{'0', '9', 'a', 'f'},
					ignoreCase: true,
					inverted:   false,
				},
			},
		},
		{
			name: "SimpleValue",
			pos:  position{line: 124, col: 1, offset: 3092},
			expr: &actionExpr{
				pos: position{line: 124, col: 16, offset: 3107},
				run: (*parser).callonSimpleValue1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 124, col: 16, offset: 3107},
					expr: &charClassMatcher{
						pos:        position{line: 124, col: 16, offset: 3107},
						val:        "[^;#\\r\\n]",
						chars:      []rune{';', '#', '\r', '\n'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "LineEnd",
			pos:  position{line: 132, col: 1, offset: 3282},
			expr: &choiceExpr{
				pos: position{line: 132, col: 12, offset: 3293},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 132, col: 12, offset: 3293},
						val:        "\r\n",
						ignoreCase: false,
					},
					&actionExpr{
						pos: position{line: 132, col: 21, offset: 3302},
						run: (*parser).callonLineEnd3,
						expr: &litMatcher{
							pos:        position{line: 132, col: 21, offset: 3302},
							val:        "\n",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 140, col: 1, offset: 3443},
			expr: &actionExpr{
				pos: position{line: 140, col: 19, offset: 3461},
				run: (*parser).callon_1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 140, col: 19, offset: 3461},
					expr: &charClassMatcher{
						pos:        position{line: 140, col: 19, offset: 3461},
						val:        "[ \\t]",
						chars:      []rune{' ', '\t'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 148, col: 1, offset: 3597},
			expr: &notExpr{
				pos: position{line: 148, col: 8, offset: 3604},
				expr: &anyMatcher{
					line: 148, col: 9, offset: 3605,
				},
			},
		},
	},
}

func (c *current) onFile1(lines interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf("\n\n\n>> File: %s // '%s'", c.pos, string(c.text))

	// convert iface to []*Line
	lsSlice := toIfaceSlice(lines)
	ls := make([]*Line, len(lsSlice))
	for i, l := range lsSlice {
		ls[i] = l.(*Line)
	}

	return NewFile(ls), nil
}

func (p *parser) callonFile1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFile1(stack["lines"])
}

func (c *current) onLine1(ws, item, le interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Line: %s // '%s'", c.pos, string(c.text))
	it, _ := item.(Item)
	return NewLine(c.pos, ws.(string), it, le.(string)), nil
}

func (p *parser) callonLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLine1(stack["ws"], stack["item"], stack["le"])
}

func (c *current) onComment1(cs, comment interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Comment: %s // '%s'\n", c.pos, string(c.text))
	return NewComment(c.pos, string(cs.([]byte)), comment.(string)), nil
}

func (p *parser) callonComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onComment1(stack["cs"], stack["comment"])
}

func (c *current) onSection1(name, ws, comment interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Section: %s // '%s'\n", c.pos, name)
	com, _ := comment.(*Comment)
	return NewSection(c.pos, name.(string), ws.(string), com), nil
}

func (p *parser) callonSection1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSection1(stack["name"], stack["ws"], stack["comment"])
}

func (c *current) onKeyValuePair1(key, ws, val, comment interface{}) (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> KeyValuePair: %s // '%s': '%s'\n", c.pos, key, val)
	com, _ := comment.(*Comment)
	return NewKeyValuePair(c.pos, key.(string), ws.(string), val.(string), com), nil
}

func (p *parser) callonKeyValuePair1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onKeyValuePair1(stack["key"], stack["ws"], stack["val"], stack["comment"])
}

func (c *current) onCommentVal1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> CommentVal: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonCommentVal1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCommentVal1()
}

func (c *current) onSectionName1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> SectionName: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonSectionName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSectionName1()
}

func (c *current) onKey1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Key: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonKey1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onKey1()
}

func (c *current) onValue3() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Value: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonValue3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onValue3()
}

func (c *current) onQuotedValue1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> QuotedValue: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonQuotedValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onQuotedValue1()
}

func (c *current) onChar8() (interface{}, error) {
	// " // ignore
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> Char: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonChar8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onChar8()
}

func (c *current) onHexDigit1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> HexDigit: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonHexDigit1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onHexDigit1()
}

func (c *current) onSimpleValue1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> SimpleValue: %s // '%s'\n", c.pos, string(c.text))
	return string(c.text), nil
}

func (p *parser) callonSimpleValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSimpleValue1()
}

func (c *current) onLineEnd3() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> LineEnd: %s\n", c.pos)
	return string(c.text), nil
}

func (p *parser) callonLineEnd3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineEnd3()
}

func (c *current) on_1() (interface{}, error) {
	lastPosition = c.pos
	lastText = string(c.text)

	//fmt.Printf(">> _ %s\n", c.pos)
	return string(c.text), nil
}

func (p *parser) callon_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.on_1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEntrypoint is returned when the specified entrypoint rule
	// does not exit.
	errInvalidEntrypoint = errors.New("invalid entrypoint")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Entrypoint creates an Option to set the rule name to use as entrypoint.
// The rule name must have been specified in the -alternate-entrypoints
// if generating the parser with the -optimize-grammar flag, otherwise
// it may have been optimized out. Passing an empty string sets the
// entrypoint to the first rule in the grammar.
//
// The default is to start parsing at the first rule in the grammar.
func Entrypoint(ruleName string) Option {
	return func(p *parser) Option {
		oldEntrypoint := p.entrypoint
		p.entrypoint = ruleName
		if ruleName == "" {
			p.entrypoint = g.rules[0].name
		}
		return Entrypoint(oldEntrypoint)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// AllowInvalidUTF8 creates an Option to allow invalid UTF-8 bytes.
// Every invalid UTF-8 byte is treated as a utf8.RuneError (U+FFFD)
// by character class matchers and is matched by the any matcher.
// The returned matched value, c.text and c.offset are NOT affected.
//
// The default is false.
func AllowInvalidUTF8(b bool) Option {
	return func(p *parser) Option {
		old := p.allowInvalidUTF8
		p.allowInvalidUTF8 = b
		return AllowInvalidUTF8(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// InitState creates an Option to set a key to a certain value in
// the global "state" store.
func InitState(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.state[key]
		p.cur.state[key] = value
		return InitState(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// state is a store for arbitrary key,value pairs that the user wants to be
	// tied to the backtracking of the parser.
	// This is always rolled back if a parsing rule fails.
	state storeDict

	// globalStore is a general store for the user to store arbitrary key-value
	// pairs that they need to manage and that they do not want tied to the
	// backtracking of the parser. This is only modified by the user and never
	// rolled back by the parser. It is always up to the user to keep this in a
	// consistent state.
	globalStore storeDict
}

type storeDict map[string]interface{}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type stateCodeExpr struct {
	pos position
	run func(*parser) error
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			state:       make(storeDict),
			globalStore: make(storeDict),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
		// start rule is rule [0] unless an alternate entrypoint is specified
		entrypoint: g.rules[0].name,
		emptyState: make(storeDict),
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64
	// entrypoint for the parser
	entrypoint string

	allowInvalidUTF8 bool

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}

	// emptyState contains an empty storeDict, which is used to optimize cloneState if global "state" store is not used.
	emptyState storeDict
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError && n == 1 { // see utf8.DecodeRune
		if !p.allowInvalidUTF8 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// Cloner is implemented by any value that has a Clone method, which returns a
// copy of the value. This is mainly used for types which are not passed by
// value (e.g map, slice, chan) or structs that contain such types.
//
// This is used in conjunction with the global state feature to create proper
// copies of the state to allow the parser to properly restore the state in
// the case of backtracking.
type Cloner interface {
	Clone() interface{}
}

// clone and return parser current state.
func (p *parser) cloneState() storeDict {
	if p.debug {
		defer p.out(p.in("cloneState"))
	}

	if len(p.cur.state) == 0 {
		if len(p.emptyState) > 0 {
			p.emptyState = make(storeDict)
		}
		return p.emptyState
	}

	state := make(storeDict, len(p.cur.state))
	for k, v := range p.cur.state {
		if c, ok := v.(Cloner); ok {
			state[k] = c.Clone()
		} else {
			state[k] = v
		}
	}
	return state
}

// restore parser current state to the state storeDict.
// every restoreState should applied only one time for every cloned state
func (p *parser) restoreState(state storeDict) {
	if p.debug {
		defer p.out(p.in("restoreState"))
	}
	p.cur.state = state
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	startRule, ok := p.rules[p.entrypoint]
	if !ok {
		p.addErr(errInvalidEntrypoint)
		return nil, p.errs.err()
	}

	p.read() // advance to first rune
	val, ok = p.parseRule(startRule)
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *stateCodeExpr:
		val, ok = p.parseStateCodeExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		state := p.cloneState()
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		p.restoreState(state)

		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	state := p.cloneState()

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn == utf8.RuneError && p.pt.w == 0 {
		// EOF - see utf8.DecodeRune
		p.failAt(false, p.pt.position, ".")
		return nil, false
	}
	start := p.pt
	p.read()
	p.failAt(true, start.position, ".")
	return p.sliceFrom(start), true
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError && p.pt.w == 0 { // see utf8.DecodeRune
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		state := p.cloneState()

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
		p.restoreState(state)
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	state := p.cloneState()

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	p.restoreState(state)

	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	state := p.cloneState()
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restoreState(state)
	p.restore(pt)

	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	state := p.cloneState()
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restoreState(state)
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseStateCodeExpr(state *stateCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseStateCodeExpr"))
	}

	err := state.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
