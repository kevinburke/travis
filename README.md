# travis

A command line binary (and library) for accessing the Travis CI API, with
several nice features.

## Usage

```
The travis binary interacts with Travis CI.

Usage:

	travis command [arguments]

The commands are:

	open                Open the latest branch build in a browser.
	version             Print the current version
	wait                Wait for tests to finish on a branch.

Use "travis help [command]" for more information about a command.
```

More commands (enable, rebuild, download-artifacts) may be added in the future.

The corresponding library is available at
`github.com/kevinburke/travis/lib`. View the library documentation at
[godoc.org](https://godoc.org/github.com/kevinburke/travis/lib).

## Installation

Find your target operating system (darwin, windows, linux) and desired bin
directory, and modify the command below as appropriate:

    curl --silent --location --output=/usr/local/bin/travis https://github.com/kevinburke/travis/releases/download/0.3/travis-linux-amd64 && chmod 755 /usr/local/bin/travis

The latest version is 0.3.

If you have a Go development environment, you can also install via source code:

    go get -u github.com/kevinburke/travis

The corresponding library is available at
`github.com/kevinburke/travis/lib`. View the library documentation at
[godoc.org](https://godoc.org/github.com/kevinburke/travis/lib).
