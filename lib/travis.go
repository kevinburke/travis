// Package travis implements a Go client for talking to the Travis CI API.
package travis

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	types "github.com/kevinburke/go-types"
	"github.com/kevinburke/rest"
	colorable "github.com/mattn/go-colorable"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sync/errgroup"
)

// GetToken looks in a config file for the Travis API token. organization is
// your Github username ("kevinburke") or organization ("golang").
func GetToken(organization string) (string, error) {
	var filename string
	var f io.ReadCloser
	var err error
	checkedLocations := make([]string, 1)
	if cfg, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		filename = filepath.Join(cfg, "travis")
		f, err = os.Open(filename)
		checkedLocations[0] = filename
	} else {
		var homeDir string
		user, userErr := user.Current()
		if userErr == nil {
			homeDir = user.HomeDir
		} else {
			homeDir = os.Getenv("HOME")
		}
		filename = filepath.Join(homeDir, "cfg", "travis")
		f, err = os.Open(filename)
		checkedLocations[0] = filename
		if err != nil { //fallback
			rcFilename := filepath.Join(homeDir, ".travis")
			f, err = os.Open(rcFilename)
			checkedLocations = append(checkedLocations, rcFilename)
		}
	}
	if err != nil {
		err = fmt.Errorf(`Couldn't find a config file in %s.

Add a configuration file with your Travis token, like this:

[organizations]

    [organizations.kevinburke]
    token = "aabbccddeeff00"

Go to https://travis-ci.org/profile/<your-username> if you need to find your token.
`, strings.Join(checkedLocations, " or "))
		return "", err
	}
	defer f.Close()
	var c FileConfig
	_, err = toml.DecodeReader(bufio.NewReader(f), &c)
	if err != nil {
		return "", err
	}
	org, err := getCaseInsensitiveOrg(organization, c.Organizations)
	if err != nil {
		return "", err
	}
	return org.Token, nil
}

// The client Version.
const Version = "0.4"

// The Host for the API.
const Host = "https://api.travis-ci.org"

// The hostname for viewing builds in a browser.
const WebHost = "https://travis-ci.org"

// Client is a HTTP client for interacting with the Travis API.
type Client struct {
	*rest.Client
	token string

	// For interacting with Build resources.
	Builds *BuildService
	// For interacting with Job resources.
	Jobs *JobService
}

type BuildService struct {
	client *Client
}

type JobService struct {
	client *Client
}

// Get retrieves the build with the given ID, or an error. include is a list of
// resources to load eagerly.
func (b *BuildService) Get(ctx context.Context, id int64, include ...string) (*Build, error) {
	path := "/build/" + strconv.FormatInt(id, 10)
	includes := strings.Join(include, ",")
	if includes != "" {
		path += "?include=" + includes
	}
	req, err := b.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	build := new(Build)
	if err := b.client.Do(req, build); err != nil {
		return nil, err
	}
	return build, nil
}

// Log represents a Travis Log object.
//
// https://developer.travis-ci.org/resource/log#Log
type Log struct {
	Type           string          `json:"@type"`
	HREF           string          `json:"@href"`
	Representation string          `json:"@representation"`
	Permissions    map[string]bool `json:"@permissions"`
	RawLogHREF     string          `json:"@raw_log_href"`
	ID             int64           `json:"id"`
	Content        string          `json:"content"`
	LogParts       []*LogPart      `json:"log_parts"`
}

// LogPart represents a log part.
type LogPart struct {
	Content string `json:"content"`
	Final   bool   `json:"final"`
	Number  int    `json:"number"`
}

// GetLog retrieves the job log for the job with the given ID, or an error.
// include is a list of resources to eager load.
func (j *JobService) GetLog(ctx context.Context, id int64, include ...string) (*Log, error) {
	path := "/job/" + strconv.FormatInt(id, 10) + "/log"
	includes := strings.Join(include, ",")
	if includes != "" {
		path += "?include=" + includes
	}
	req, err := j.client.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	build := new(Log)
	if err := j.client.Do(req, build); err != nil {
		return nil, err
	}
	return build, nil
}

// NewClient creates a new Client.
func NewClient(token string) *Client {
	c := &Client{
		Client: rest.NewClient("", "", Host),
		token:  token,
	}
	c.Builds = &BuildService{client: c}
	c.Jobs = &JobService{client: c}
	return c
}

// NewRequest creates a new HTTP request to hit the given endpoint.
func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := c.Client.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "travis-go/"+Version+" (github.com/kevinburke/travis) "+req.Header.Get("User-Agent"))
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Travis-API-Version", "3")
	return req, nil
}

// getCaseInsensitiveOrg finds the key in the list of orgs. This is a case
// insensitive match, so if key is "ShyP" and orgs has a key named "sHYp",
// that will count as a match.
func getCaseInsensitiveOrg(key string, orgs map[string]organization) (organization, error) {
	for k, _ := range orgs {
		lower := strings.ToLower(k)
		if _, ok := orgs[lower]; !ok {
			orgs[lower] = orgs[k]
			delete(orgs, k)
		}
	}
	lowerKey := strings.ToLower(key)
	if o, ok := orgs[lowerKey]; ok {
		return o, nil
	} else {
		return organization{}, fmt.Errorf(`Couldn't find organization %s in the config.

Go to https://travis-ci.org/profile/<your-username> if you need to create a token.
`, key)
	}
}

// FileConfig represents the structure of your ~/cfg/travis config file.
type FileConfig struct {
	Organizations map[string]organization
}

type organization struct {
	Token string
}

// Pagination contains details on paging through API responses.
type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Count   int  `json:"count"`
	IsFirst bool `json:"is_first"`
	IsLast  bool `json:"is_last"`
}

// ListResponse represents a Travis response for a list of resources.
type ListResponse struct {
	Type           string     `json:"@type"`
	HREF           string     `json:"@href"`
	Representation string     `json:"@representation"`
	Pagination     Pagination `json:"@pagination"`
	// Set this to whatever data you want to deserialize before calling
	// json.Unmarshal/client.Do.
	Data interface{}
}

// Build represents a Build in Travis CI.
//
// https://developer.travis-ci.org/resource/build#Build
type Build struct {
	Type           string         `json:"@type"`
	HREF           string         `json:"@href"`
	Representation string         `json:"@representation"`
	ID             int64          `json:"id"`
	Number         string         `json:"number"`
	State          string         `json:"state"`
	PreviousState  string         `json:"previous_state"`
	Duration       int64          `json:"duration"`
	StartedAt      time.Time      `json:"started_at"`
	FinishedAt     types.NullTime `json:"finished_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Branch         Branch         `json:"branch"`
	Repository     Repository     `json:"repository"`
	Commit         *Commit        `json:"commit"`
	Jobs           []*Job         `json:"jobs"`
}

// Job represents a Job in Travis CI. A Build has one or more Jobs.
//
// https://developer.travis-ci.org/resource/job#Job
type Job struct {
	Type           string          `json:"@type"`
	HREF           string          `json:"@href"`
	Representation string          `json:"@representation"`
	Permissions    map[string]bool `json:"@permissions"`
	ID             int64           `json:"id"`
	AllowFailure   bool            `json:"allow_failure"`
	Number         string          `json:"number"`
	State          string          `json:"state"`
	StartedAt      time.Time       `json:"started_at"`
	FinishedAt     types.NullTime  `json:"finished_at"`
	Queue          string          `json:"queue"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`

	Config *Config `json:"config"`
}

// Not documented, but represents your Travis CI config in JSON form.
type Config struct {
	Language     string   `json:"language"`
	BeforeScript []string `json:"before_script"`
	Script       []string `json:"script"`
	// "true", "false", "required"
	Sudo   string `json:"sudo"`
	OS     string `json:"os"`
	Group  string `json:"group"`
	Extras map[string]interface{}
}

const stepWidth = 45

var stepPadding = fmt.Sprintf("%%-%ds", stepWidth)

func isatty() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}

// BuildStatistics returns statistics about a build as a multiline string.
func (c *Client) BuildStatistics(ctx context.Context, b *Build) (string, error) {
	if len(b.Jobs) == 0 {
		return "(no jobs)", nil
	}
	steps := make([][]*Step, len(b.Jobs))
	group, errctx := errgroup.WithContext(ctx)
	for i := range b.Jobs {
		i := i
		group.Go(func() error {
			log, err := c.Jobs.GetLog(errctx, b.Jobs[i].ID)
			if err != nil {
				return err
			}
			steps[i] = ParseLog(log.Content)
			return nil
		})
	}
	if err := group.Wait(); err != nil {
		return "", err
	}
	for j := len(steps[0]) - 1; j >= 0; j-- {
		longStep := false
		for i := range b.Jobs {
			if j >= len(steps[i]) {
				longStep = true
				break
			}
			longStep = steps[i][j].End.Sub(steps[i][j].Start) > 10*time.Millisecond
		}
		if longStep {
			continue
		}
		// we can delete "j"
		for i := range b.Jobs {
			copy(steps[i][j:], steps[i][j+1:])
			steps[i][len(steps[i])-1] = nil // or the zero value of T
			steps[i] = steps[i][:len(steps[i])-1]
		}
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(stepPadding, "Step"))
	l := stepWidth
	for i := range steps {
		if b.Jobs[i].Config == nil || b.Jobs[i].Config.Language == "" {
			buf.WriteString(fmt.Sprintf("%-8d", i))
		} else {
			lang := b.Jobs[i].Config.Language
			valI, ok := b.Jobs[i].Config.Extras[lang]
			if !ok {
				buf.WriteString(fmt.Sprintf("%-8d", i))
				continue
			}
			val, ok := valI.(string)
			if !ok {
				buf.WriteString(fmt.Sprintf("%-8d", i))
				continue
			}
			if len(val) > 8-2 {
				val = fmt.Sprintf("%s… ", val[:(8-2)])
			} else {
				val = fmt.Sprintf("%-8s", val)
			}
			buf.WriteString(val)
		}
		l += 8
	}
	buf.WriteString(fmt.Sprintf("\n%s\n", strings.Repeat("=", l)))
	// sorta backwards iteration, but eh
	// TODO: strip steps where runtime was < 10ms on all builders
	for i := range steps[0] {
		stepName := strings.Replace(steps[0][i].Name, "\n", "\\n", -1)
		if len(stepName) > stepWidth-2 {
			stepName = fmt.Sprintf("%s… ", stepName[:(stepWidth-2)])
		} else {
			stepName = fmt.Sprintf(stepPadding, stepName)
		}
		buf.WriteString(stepName)
		for j := range steps {
			if i >= len(steps[j]) {
				fmt.Fprintf(&buf, "%-8s", "")
				continue
			}
			runtime := steps[j][i].End.Sub(steps[j][i].Start)
			var dur time.Duration
			if runtime > time.Minute {
				dur = runtime.Round(time.Second)
			} else {
				dur = runtime.Round(10 * time.Millisecond)
			}
			if b.Jobs[j].State == "failed" || b.Jobs[j].State == "errored" && isatty() && i == len(steps[j])-1 {
				// color the output red
				fmt.Fprintf(&buf, "\033[38;05;160m%-8s\033[0m", dur.String())
				continue
			}
			fmt.Fprintf(&buf, "%-8s", dur.String())
		}
		buf.WriteString("\n")
	}
	return buf.String(), nil
}

func (c *Config) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "sudo":
			t, ok := v.(bool)
			if ok {
				if t {
					c.Sudo = "true"
				} else {
					c.Sudo = "false"
				}
				continue
			}
			fallthrough
		case "language", "os", "group":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf("could not convert %s to string: %v", k, v)
			}
			switch k {
			case "language":
				c.Language = s
			case "os":
				c.OS = s
			case "group":
				c.Group = s
			case "sudo":
				c.Sudo = s
			}
		case "before_script", "script":
			beforeIArr, ok := v.([]interface{})
			if !ok {
				return fmt.Errorf("could not convert before_script to array: %v", v)
			}
			for i := range beforeIArr {
				s, ok := beforeIArr[i].(string)
				if !ok {
					return fmt.Errorf("could not convert %s item to string: %v", k, beforeIArr[i])
				}
				switch k {
				case "before_script":
					if c.BeforeScript == nil {
						c.BeforeScript = make([]string, len(beforeIArr))
					}
					c.BeforeScript[i] = s
				case "script":
					if c.Script == nil {
						c.Script = make([]string, len(beforeIArr))
					}
					c.Script[i] = s
				}
			}
		default:
			if c.Extras == nil {
				c.Extras = make(map[string]interface{})
			}
			c.Extras[k] = v
		}
	}
	return nil
}

// WebURL returns the URL for viewing this build in a web browser.
func (b Build) WebURL() string {
	return fmt.Sprintf("%s/%s/builds/%d", WebHost, b.Repository.Slug, b.ID)
}

// Commit represents a Git commit in Travis CI.
//
// https://developer.travis-ci.org/resource/commit#Commit
type Commit struct {
	Type           string    `json:"@type"`
	HREF           string    `json:"@href"`
	Representation string    `json:"@representation"`
	SHA            string    `json:"sha"`
	Ref            string    `json:"ref"`
	Message        string    `json:"message"`
	CompareURL     string    `json:"compare_url"`
	CommittedAt    time.Time `json:"committed_at"`
}

// Branch represents a Git branch in Travis CI.
//
// https://developer.travis-ci.org/resource/branch#Branch
type Branch struct {
	Type           string `json:"@type"`
	HREF           string `json:"@href"`
	Representation string `json:"@representation"`
	Name           string `json:"name"`
}

// Repository represents a repository in Travis CI.
//
// https://developer.travis-ci.org/resource/repository#Repository
type Repository struct {
	Type           string `json:"@type"`
	HREF           string `json:"@href"`
	Representation string `json:"@representation"`
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
}

func (r *ListResponse) UnmarshalJSON(b []byte) error {
	r2 := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &r2); err != nil {
		return err
	}
	if err := json.Unmarshal(r2["@type"], &r.Type); err != nil {
		return err
	}
	if err := json.Unmarshal(r2["@href"], &r.HREF); err != nil {
		return err
	}
	if err := json.Unmarshal(r2["@representation"], &r.Representation); err != nil {
		return err
	}
	if err := json.Unmarshal(r2["@pagination"], &r.Pagination); err != nil {
		return err
	}
	if err := json.Unmarshal(r2[r.Type], &r.Data); err != nil {
		return err
	}
	return nil
}

var foldStart = "travis_fold:start:"
var timeStart = "travis_time:start:"
var foldEnd = "travis_fold:end:"
var timeEnd = "travis_time:end:"

func getStep(text string) (*Step, bool, string) {
	timeStartIdx := strings.Index(text, timeStart)
	if timeStartIdx == -1 {
		return nil, false, text
	}
	text = text[timeStartIdx+len(timeStart):]
	escapes := strings.Index(text, "\r\x1b[0K")
	if escapes == -1 {
		return nil, false, text
	}
	text = text[escapes+len("\r\x1b[0K"):]
	endIdx := strings.IndexByte(text, '\n')
	var name string
	if endIdx == -1 {
		name = text
	} else {
		name = text[:endIdx]
		text = text[endIdx:]
	}
	name = strings.TrimSpace(strings.TrimPrefix(name, "$ "))
	step := &Step{
		Name: stripANSI(name),
	}
	if step.Name == "" {
		step.Name = "(no name)"
	}
	endTimeIdx := strings.Index(text, timeEnd)
	if endTimeIdx == -1 {
		return nil, false, text
	}
	text = text[endTimeIdx+len(timeEnd):]
	lineIdx := strings.Index(text, "start=")
	if lineIdx == -1 {
		return nil, false, text
	}
	text = text[lineIdx+len("start="):]
	commaIdx := strings.IndexByte(text, ',')
	if commaIdx == -1 {
		return nil, false, text
	}
	start, err := strconv.ParseInt(text[:commaIdx], 10, 64)
	if err != nil {
		return nil, false, text
	}
	// start is in nanoseconds
	step.Start = time.Unix(0, start).UTC()

	lineIdx = strings.Index(text, "finish=")
	if lineIdx == -1 {
		return nil, false, text
	}
	text = text[lineIdx+len("finish="):]
	commaIdx = strings.IndexByte(text, ',')
	if commaIdx == -1 {
		return nil, false, text
	}
	end, err := strconv.ParseInt(text[:commaIdx], 10, 64)
	if err != nil {
		return nil, false, text
	}
	// start is in nanoseconds
	step.End = time.Unix(0, end).UTC()
	if len(text) <= commaIdx+1 {
		return nil, false, text
	}
	text = text[commaIdx+1:]
	return step, true, text
}

// Step represents a step of a build. These get parsed out of the log files;
// it's not clear that it's possible to get them any other way.
type Step struct {
	Name       string
	Start, End time.Time
}

func stripANSI(ansi string) string {
	var buf strings.Builder
	w := colorable.NewNonColorable(&buf)
	if _, err := io.WriteString(w, ansi); err != nil {
		panic(err)
	}
	return buf.String()
}

// ParseLog parses a log file, returning the names of each step in the log, and
// the amount of time each step took.
func ParseLog(log string) []*Step {
	steps := make([]*Step, 0)
	var currentStep *Step
	for {
		foldStartIdx := strings.Index(log, foldStart)
		timeStartIdx := strings.Index(log, timeStart)
		if foldStartIdx == -1 && timeStartIdx == -1 {
			break
		}
		ok := false
		if foldStartIdx == -1 || (timeStartIdx < foldStartIdx) {
			currentStep, ok, log = getStep(log)
			if !ok {
				break
			}
			steps = append(steps, currentStep)
			continue
		}
		// in a folded step
		log = log[foldStartIdx+len(foldStart):]
		endIdx := strings.IndexByte(log, '\r')
		if endIdx == -1 {
			break
		}
		foldNameANSI := strings.TrimSpace(log[:endIdx])
		foldName := stripANSI(foldNameANSI)
		log = log[endIdx:]
		foldSteps := make([]*Step, 0)
		for {
			timeStartIdx := strings.Index(log, timeStart)
			foldEndIdx := strings.Index(log, foldEnd)
			if foldEndIdx == -1 {
				break
			}
			if timeStartIdx >= 0 && timeStartIdx < foldEndIdx {
				currentStep, ok, log = getStep(log)
				if !ok {
					break
				}
				currentStep.Name = foldName + ":" + currentStep.Name
				foldSteps = append(foldSteps, currentStep)
				continue
			}
			// no more timings
			if len(foldSteps) == 1 {
				foldSteps[0].Name = foldName
			}
			steps = append(steps, foldSteps...)
			log = log[foldEndIdx+len(foldEnd):]
			break
		}
	}
	return steps
}
