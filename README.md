# travis

A command line binary (and library) for accessing the Travis CI API, with
several nice features. Running `travis wait` at the command line will yield the
following step-by-step output:

```
$ travis wait
Waiting for latest build on master to complete
Build 366972330 running (3m49s elapsed)
Build on master succeeded!

Step                                         1.10.x  master
=============================================================
git.checkout                                 330ms   590ms
GIMME_OUTPUT="$(gimme 1.10 | tee -a $HOME/.… 3.9s    2m13s
cache.1:Installing caching utilities         30ms    30ms
cache.1:attempting to download cache archiv… 4.77s   5.56s
cache.1:adding /home/travis/gopath/pkg to c… 2.22s   2.84s
before_script                                17.84s  13.24s
make race-test                               5.54s   19.9s
cache.2:nothing changed, not updating cache  1.32s   1.79s

Tests on master took 3m52s. Quitting.
```

You don't need to configure anything besides your API token.

Failed builds will display in your terminal in red. Soon, we will print the
output from failed build steps to the screen.

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
