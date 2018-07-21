// The travis binary interacts with Travis CI.
//
// Usage:
//
//	travis command [arguments]
//
// The commands are:
//
//	enable              Enable builds for this repository.
//	open                Open the latest branch build in a browser.
//	sync                Sync repos for the account.
//	version             Print the current version
//	wait                Wait for tests to finish on a branch.
//
// Use "travis help [command]" for more information about a command.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/kevinburke/bigtext"
	git "github.com/kevinburke/go-git"
	"github.com/kevinburke/remoteci"
	travis "github.com/kevinburke/travis/lib"
	"github.com/pkg/browser"
)

const help = `The travis binary interacts with Travis CI.

Usage: 

	travis command [arguments]

The commands are:

	enable              Enable builds for this repository.
	open                Open the latest branch build in a browser.
	sync                Sync repos for the account.
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

func newClient(org string) (*travis.Client, error) {
	token, err := travis.GetToken(org)
	if err != nil {
		return nil, err
	}
	return travis.NewClient(token), nil
}

var errNoBuilds = errors.New("travis: no builds")

func getLatestBuild(client *travis.Client, org, repo, branch string) (*travis.Build, error) {
	builds, err := getBuilds(client, org, repo, branch)
	if err != nil {
		return nil, err
	}
	if len(builds) == 0 {
		return nil, errNoBuilds
	}
	return builds[0], nil
}

func getBuild(client *travis.Client, id int64) (*travis.Build, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Builds.Get(ctx, id, "build.jobs", "job.config")
}

func getBuilds(client *travis.Client, org, repo, branch string) ([]*travis.Build, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	slug := url.PathEscape(org + "/" + repo)
	req, err := client.NewRequest("GET", "/repo/"+slug+"/builds?branch.name="+url.QueryEscape(branch), nil)
	if err != nil {
		return nil, err
	}
	// https://developer.travis-ci.org/authentication
	builds := make([]*travis.Build, 0)
	resp := &travis.ListResponse{
		Data: &builds,
	}
	req = req.WithContext(ctx)
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}
	return builds, nil
}

func doEnable(flags *flag.FlagSet, remoteStr string) {
	remote, err := git.GetRemoteURL(remoteStr)
	checkError(err, "getting remote URL")
	client, err := newClient(remote.Path)
	checkError(err, "getting token")
	slug := remote.Path + "/" + remote.RepoName
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Repos.Activate(ctx, slug); err != nil {
		failError(err, "activating repository")
	}
	fmt.Printf("%s/%s enabled\n", travis.WebHost, slug)
}

func doOpen(flags *flag.FlagSet, remoteStr string) {
	args := flags.Args()
	branch, err := getBranchFromArgs(args)
	checkError(err, "getting git branch")
	remote, err := git.GetRemoteURL(remoteStr)
	checkError(err, "getting remote URL")

	client, err := newClient(remote.Path)
	checkError(err, "getting token")
	latestBuild, err := getLatestBuild(client, remote.Path, remote.RepoName, branch)
	checkError(err, "getting latest build")
	if err := browser.OpenURL(latestBuild.WebURL()); err != nil {
		checkError(err, "opening url "+latestBuild.WebURL())
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
		failError(err, msg)
	}
}

func failError(err error, msg string) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", msg, err)
	os.Exit(1)
}

// isHttpError checks if the given error is a request timeout or a network
// failure - in those cases we want to just retry the request.
func isHttpError(err error) bool {
	if err == nil {
		return false
	}
	// some net.OpError's are wrapped in a url.Error
	if uerr, ok := err.(*url.Error); ok {
		err = uerr.Err
	}
	switch err := err.(type) {
	default:
		return false
	case *net.OpError:
		return err.Op == "dial" && err.Net == "tcp"
	case *net.DNSError:
		return true
	// Catchall, this needs to go last.
	case net.Error:
		return err.Timeout() || err.Temporary()
	}
}

// getMinTipLength compares two strings and returns the length of the
// shortest
func getMinTipLength(remoteTip string, localTip string) int {
	var minTipLength int
	if len(remoteTip) <= len(localTip) {
		minTipLength = len(remoteTip)
	} else {
		minTipLength = len(localTip)
	}
	return minTipLength
}

func shouldPrint(lastPrinted time.Time, duration time.Duration, latestBuild, previousBuild *travis.Build) bool {
	now := time.Now()
	var buildDuration time.Duration
	if previousBuild == nil {
		buildDuration = 5 * time.Minute
	} else {
		buildDuration = previousBuild.FinishedAt.Time.Sub(previousBuild.StartedAt)
	}
	var durToUse time.Duration
	timeRemaining := buildDuration - duration
	switch {
	case timeRemaining > 25*time.Minute:
		durToUse = 3 * time.Minute
	case timeRemaining > 8*time.Minute:
		durToUse = 2 * time.Minute
	case timeRemaining > 5*time.Minute:
		durToUse = 30 * time.Second
	case timeRemaining > 3*time.Minute:
		durToUse = 20 * time.Second
	case timeRemaining > time.Minute:
		durToUse = 15 * time.Second
	default:
		durToUse = 10 * time.Second
	}
	return lastPrinted.Add(durToUse).Before(now)
}

func doWait(branch, remoteStr string) error {
	remote, err := git.GetRemoteURL(remoteStr)
	if err != nil {
		return err
	}
	tip, err := git.Tip(branch)
	if err != nil {
		return err
	}
	client, err := newClient(remote.Path)
	if err != nil {
		return err
	}
	fmt.Println("Waiting for latest build on", branch, "to complete")
	var lastPrintedAt time.Time
	var previousBuild *travis.Build
	builds, err := getBuilds(client, remote.Path, remote.RepoName, branch)
	if err == nil {
		for i := 1; i < len(builds); i++ {
			if builds[i].State == "passed" {
				previousBuild = builds[i]
				break
			}
		}
	}
	for {
		latestBuild, err := getLatestBuild(client, remote.Path, remote.RepoName, branch)
		if err != nil {
			if isHttpError(err) {
				fmt.Printf("Caught network error: %s. Continuing\n", err.Error())
				lastPrintedAt = time.Now()
				time.Sleep(2 * time.Second)
				continue
			}
			if err == errNoBuilds {
				return fmt.Errorf("No results, are you sure there are tests for %s/%s?\n",
					remote.Path, remote.RepoName)
			}
			return err
		}
		if latestBuild.Commit == nil {
			return fmt.Errorf("Latest build on %s/%s is not a commit?\n",
				remote.Path, remote.RepoName)
		}
		c := bigtext.Client{
			Name:    fmt.Sprintf("%s (github.com/kevinburke/travis)", remote.RepoName),
			OpenURL: latestBuild.WebURL(),
		}
		maxTipLengthToCompare := getMinTipLength(latestBuild.Commit.SHA, tip)
		if latestBuild.Commit.SHA[:maxTipLengthToCompare] != tip[:maxTipLengthToCompare] {
			fmt.Printf("Latest build in Travis is %s, waiting for %s...\n",
				latestBuild.Commit.SHA[:maxTipLengthToCompare], tip[:maxTipLengthToCompare])
			lastPrintedAt = time.Now()
			time.Sleep(5 * time.Second)
			continue
		}
		var duration time.Duration
		if latestBuild.FinishedAt.Valid {
			duration = latestBuild.FinishedAt.Time.Sub(latestBuild.StartedAt).Round(time.Second)
		} else {
			duration = time.Since(latestBuild.StartedAt).Round(time.Second)
		}
		if latestBuild.State == "passed" {
			fmt.Printf("Build on %s succeeded!\n\n", branch)
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			build, err := client.Builds.Get(ctx, latestBuild.ID, "build.jobs", "job.config")
			if err == nil {
				stats, err := client.BuildSummary(ctx, build)
				if err == nil {
					fmt.Print(stats)
				} else {
					fmt.Printf("error fetching build summary: %v\n", err)
				}
			} else {
				fmt.Printf("error getting build: %v\n", err)
			}
			fmt.Printf("\nTests on %s took %s. Quitting.\n", branch, duration.String())
			c.Display(branch + " build complete!")
			break
		}
		if latestBuild.State == "failed" || latestBuild.State == "errored" {
			build, err := getBuild(client, latestBuild.ID)
			if err == nil {
				ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
				defer cancel()
				stats, err := client.BuildSummary(ctx, build)
				if err == nil {
					fmt.Print(stats)
				} else {
					fmt.Printf("error fetching build stats: %v\n", err)
				}
			} else {
				fmt.Printf("error getting build: %v\n", err)
			}
			fmt.Printf("\nURL: %s\n", latestBuild.WebURL())
			err = fmt.Errorf("Build on %s failed!\n\n", branch)
			c.Display("build failed")
			return err
		}
		if latestBuild.State == "started" {
			// Show more and more output as we approach the duration of the previous
			// successful build.
			if shouldPrint(lastPrintedAt, duration, latestBuild, previousBuild) {
				fmt.Printf("Build %d running (%s elapsed)\n", latestBuild.ID, duration.String())
				lastPrintedAt = time.Now()
			}
		} else if false { // queued build
			cost := remoteci.GetEffectiveCost(duration)
			centsPortion := cost % 100
			dollarPortion := cost / 100
			costStr := fmt.Sprintf("$%d.%.2d", dollarPortion, centsPortion)
			if lastPrintedAt.Add(12 * time.Second).Before(time.Now()) {
				fmt.Printf("State is %s (queued for %s, cost %s), trying again\n",
					latestBuild.State, duration.String(), costStr)
				lastPrintedAt = time.Now()
			}
		} else {
			fmt.Printf("State is %s, trying again\n", latestBuild.State)
			lastPrintedAt = time.Now()
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}

func doSync(flags *flag.FlagSet) {
	remote, err := git.GetRemoteURL("origin")
	checkError(err, "getting remote URL")

	client, err := newClient(remote.Path)
	checkError(err, "getting token")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	u, err := client.Users.Current(ctx)
	checkError(err, "getting user id")
	if err := client.Users.Sync(ctx, u.ID); err != nil {
		checkError(err, "syncing account")
	}
	fmt.Println(u.Login + " synced")
}

func main() {
	enableflags := flag.NewFlagSet("open", flag.ExitOnError)
	enableRemote := enableflags.String("remote", "origin", "Git remote to use")
	enableflags.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: enable [--remote=origin]

Enable Travis CI builds for this repository.

`)
		enableflags.PrintDefaults()
	}
	openflags := flag.NewFlagSet("open", flag.ExitOnError)
	openRemote := openflags.String("remote", "origin", "Git remote to use")
	syncflags := flag.NewFlagSet("sync", flag.ExitOnError)
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
		os.Exit(2)
	}
	subargs := args[1:]
	switch flag.Arg(0) {
	case "enable":
		enableflags.Parse(subargs)
		doEnable(enableflags, *enableRemote)
	case "open":
		openflags.Parse(subargs)
		doOpen(openflags, *openRemote)
	case "sync":
		syncflags.Parse(subargs)
		doSync(syncflags)
	case "version":
		fmt.Fprintf(os.Stderr, "travis version %s\n", travis.Version)
		os.Exit(1)
	case "wait":
		waitflags.Parse(subargs)
		args := waitflags.Args()
		branch, err := getBranchFromArgs(args)
		checkError(err, "getting git branch")
		err = doWait(branch, *waitRemote)
		checkError(err, "waiting for branch")
	default:
		fmt.Fprintf(os.Stderr, "travis: unknown command %q\n\n", flag.Arg(0))
		usage()
		os.Exit(2)
	}
}
