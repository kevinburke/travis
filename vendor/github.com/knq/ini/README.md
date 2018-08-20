# About ini [![Build Status][t-badge]][t-link] [![Coverage][c-badge]][c-link]

`ini` is a simple [Go][go-project] package for manipulating [ini files][wiki-ini].

`ini` is mostly a simple wrapper around the [`ini/parser` package](/parser)
also contained in this repository. `ini/parser` was implemented by generating a
[Pigeon][pigeon] parser from a [PEG grammar][wiki-peg].

With the correct configuration, the `ini` package is able to read [git
config][git-config] files, very simple [TOML][toml] files, and [Java
Properties][wiki-dotproperties] files.

## Why Another ini Package? ##

Prior to writing this package, a number of existing Go ini packages/parsers
were investigated. The packages available at the time did not possess the
complete feature set needed: specifically, the available packages did not work
well with badly formatted files (and their parsers were not easily fixable),
they would erase any comments/spacing when writing out a modified ini file, and
they were not written in [idiomatic Go][go-idiomatic].

As such, it was necessary to author a new package that could work with a
variety of badly formatted ini files, in idiomatic Go, and provide a simple
interface to reading/writing/manipulating ini files.

## Installing ##

Install in the usual Go fashion:

    go get -u github.com/knq/ini

## Using ##

`ini` can be used similarly to the following:

```go
// examples/test2/main.go
package main

import (
	"fmt"
	"log"

	"github.com/knq/ini"
)

var (
	data = `
	firstkey = one

	[some section]
	key = blah ; comment

	[another section]
	key = blah`

	gitconfig = `
	[difftool "gdmp"]
	cmd = ~/gdmp/x "$LOCAL" "$REMOTE"
	`
)

func main() {
	f, err := ini.LoadString(data)
	if err != nil {
		log.Fatal(err)
	}

	s := f.GetSection("some section")

	fmt.Printf("some section.key: %s\n", s.Get("key"))
	s.SetKey("key2", "another value")
	f.Write("out.ini")

	// create a gitconfig parser
	g, err := ini.LoadString(gitconfig)
	if err != nil {
		log.Fatal(err)
	}

	// setup gitconfig name/key manipulation functions
	g.SectionManipFunc = ini.GitSectionManipFunc
	g.SectionNameFunc = ini.GitSectionNameFunc

	fmt.Printf("difftool.gdmp.cmd: %s\n", g.GetKey("difftool.gdmp.cmd"))
}
```

Please see [the GoDoc API page][godoc] for a full API listing.

## TODO

* convert to github.com/alecthomas/participle parser

[c-badge]: https://coveralls.io/repos/github/knq/ini/badge.svg?branch=master
[c-link]: https://coveralls.io/github/knq/ini?branch=master
[git-config]: http://git-scm.com/docs/git-config
[godoc]: http://godoc.org/github.com/knq/ini
[go-idiomatic]: https://golang.org/doc/effective_go.html
[go-project]: http://www.golang.org/project/
[pigeon]: https://github.com/mna/pigeon/
[t-badge]: https://travis-ci.org/knq/ini.svg
[t-link]: https://travis-ci.org/knq/ini
[toml]: https://github.com/toml-lang/toml
[wiki-dotproperties]: https://en.wikipedia.org/wiki/.properties
[wiki-ini]: https://en.wikipedia.org/wiki/INI_file
[wiki-peg]: https://en.wikipedia.org/wiki/Parsing_expression_grammar
