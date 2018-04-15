package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/kevinburke/go-circle/wait"
	git "github.com/kevinburke/go-git"
	travis "github.com/kevinburke/travis/lib"
	"github.com/pkg/browser"
)

const help = `The travis binary interacts with Travis CI.

Usage: 

	travis command [arguments]

The commands are:

	open                Open the latest branch build in a browser.
	version             Print the current version
	wait                Wait for tests to finish on a branch.

Use "travis help [command]" for more information about a command.
`

func usage() {
	fmt.Fprintf(os.Stderr, help)
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func doOpen(flags *flag.FlagSet) {
	args := flags.Args()
	branch, err := getBranchFromArgs(args)
	checkError(err, "getting git branch")
	remote, err := git.GetRemoteURL("origin")
	checkError(err, "getting remote URL")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	token, err := travis.GetToken(remote.Path)
	checkError(err, "finding token")
	client := travis.NewClient(token)
	slug := url.PathEscape(remote.Path + "/" + remote.RepoName)
	req, err := client.NewRequest("GET", "/repo/"+slug+"/builds?branch.name="+url.QueryEscape(branch), nil)
	checkError(err, "creating HTTP client")
	// https://developer.travis-ci.org/authentication
	builds := make([]*travis.Build, 0)
	resp := &travis.Response{
		Data: &builds,
	}
	req = req.WithContext(ctx)
	if err := client.Do(req, resp); err != nil {
		checkError(err, "fetching recent builds")
	}
	latestBuild := builds[0]
	u := fmt.Sprintf("%s/%s/builds/%d", travis.WebHost, latestBuild.Repository.Slug, latestBuild.ID)
	if err := browser.OpenURL(u); err != nil {
		checkError(err, "opening url "+u)
	}
}

// Given a set of command line args, return the git branch or an error. Returns
// the current git branch if no argument is specified
func getBranchFromArgs(args []string) (string, error) {
	if len(args) == 0 {
		return git.CurrentBranch()
	} else {
		return args[0], nil
	}
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %s: %v\n", msg, err)
		os.Exit(1)
	}
}

func main() {
	openflags := flag.NewFlagSet("open", flag.ExitOnError)
	waitflags := flag.NewFlagSet("wait", flag.ExitOnError)
	waitRemote := waitflags.String("remote", "origin", "Git remote to use")
	waitflags.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: wait [refspec]

Wait for builds to complete, then print a descriptive output on success or
failure. By default, waits on the current branch, otherwise you can pass a
branch to wait for.

`)
		waitflags.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		usage()
		return
	}
	subargs := args[1:]
	switch flag.Arg(0) {
	case "open":
		openflags.Parse(subargs)
		doOpen(openflags)
	case "version":
		fmt.Fprintf(os.Stderr, "travis version %s\n", travis.Version)
		os.Exit(1)
	case "wait":
		waitflags.Parse(subargs)
		args := waitflags.Args()
		branch, err := getBranchFromArgs(args)
		checkError(err, "getting git branch")
		err = wait.Wait(branch, *waitRemote)
		checkError(err, "waiting for branch")
	default:
		fmt.Fprintf(os.Stderr, "travis: unknown command %q\n\n", flag.Arg(0))
		usage()
	}
}
