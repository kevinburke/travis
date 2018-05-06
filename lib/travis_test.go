package travis

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/kevinburke/rest"
)

func TestMarshalBuilds(t *testing.T) {
	r := new(ListResponse)
	builds := make([]Build, 0)
	r.Data = &builds
	if err := json.Unmarshal(buildResponse, r); err != nil {
		t.Fatal(err)
	}
	if r.Type != "builds" {
		t.Errorf("bad decode: want 'builds' got %q", r.Type)
	}
	if r.Representation != "standard" {
		t.Errorf("bad decode: want 'builds' got %q", r.Type)
	}
	if r.Pagination.Limit != 25 {
		t.Errorf("bad decode: want 'limit 25' got %q", r.Pagination.Limit)
	}
	if len(builds) != 2 {
		t.Errorf("bad builds: got %d want 2", len(builds))
	}
	if builds[0].ID != 366635873 {
		t.Errorf("bad build id: got %d want 366635873", builds[0].ID)
	}
	if builds[0].Branch.Name != "master" {
		t.Errorf("bad branch name: got %s want master", builds[0].Branch.Name)
	}
}

func TestMarshalJob(t *testing.T) {
	t.Parallel()
	j := new(Job)
	if err := json.Unmarshal(jobResponse, j); err != nil {
		t.Fatal(err)
	}
}

var userContent = []byte(`{
  "@type": "user",
  "@href": "/user/6151",
  "@representation": "standard",
  "@permissions": {
    "read": true,
    "sync": true
  },
  "id": 6151,
  "login": "kevinburke",
  "name": "Kevin Burke",
  "github_id": 234019,
  "avatar_url": "https://avatars1.githubusercontent.com/u/234019?v=4",
  "education": false,
  "is_syncing": false,
  "synced_at": "2018-05-15T00:21:36Z"
}`)

func TestMarshalUser(t *testing.T) {
	t.Parallel()
	u := new(User)
	if err := json.Unmarshal(userContent, u); err != nil {
		t.Fatal(err)
	}
	if u.ID != 6151 {
		t.Errorf("bad ID: want %d got %d", 6151, u.ID)
	}
}

func TestParseLog(t *testing.T) {
	t.Parallel()
	steps := ParseLog(logContent)
	if len(steps) != 13 {
		t.Errorf("parseSteps: want 13 steps, got %d", len(steps))
	}
}

func TestParseLogFailure(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/failure.txt")
	if err != nil {
		t.Fatal(err)
	}
	steps := ParseLog(string(data))
	if len(steps) != 14 {
		t.Errorf("parseSteps: want 14 steps, got %d", len(steps))
	}
	if steps[11].ReturnCode != 2 {
		t.Errorf("step rc: want 2 got %d", steps[11].ReturnCode)
	}
}

func TestParseLogFailureOtherway(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/failure-before-script.txt")
	if err != nil {
		t.Fatal(err)
	}
	steps := ParseLog(string(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(steps) != 11 {
		t.Errorf("parseSteps: want 11 steps, got %d", len(steps))
	}
	if steps[10].ReturnCode != 1 {
		t.Errorf("step rc: want 1 got %d", steps[10].ReturnCode)
	}
	wantOutput := "$ which barbang"
	if steps[10].Output != wantOutput {
		t.Errorf("step output: want %q, got %q", wantOutput, steps[10].Output)
	}
}

var stepsSink []*Step

func BenchmarkParseLog(b *testing.B) {
	data, err := ioutil.ReadFile("testdata/failure.txt")
	if err != nil {
		b.Fatal(err)
	}
	s := string(data)
	b.ResetTimer()
	b.SetBytes(int64(len(data)))
	for i := 0; i < b.N; i++ {
		stepsSink = ParseLog(s)
	}
}

func TestMatchExitLine(t *testing.T) {
	toMatch := "\r\n\u001b[32;1mThe command \"make race-test\" exited with 2.\u001b[0m"
	match := commandExitRx.FindStringSubmatch(toMatch)
	if match == nil {
		t.Errorf("no match")
	}
	if match[2] != "make race-test" {
		t.Errorf("bad match 2: want make race-test got %s", match[2])
	}
	if match[4] != "2" {
		t.Errorf("bad match 3: want 2 got %s", match[4])
	}

	toMatch = "\r\n\x1b[31;1mThe command \"which barbang\" failed and exited with 1 during .\x1b[0m"
	match = commandExitRx.FindStringSubmatch(toMatch)
	if match == nil {
		t.Errorf("no match")
	}
	if match[2] != "which barbang" {
		t.Errorf("bad match 2: want make race-test got %s", match[2])
	}
	if match[4] != "1" {
		t.Errorf("bad match 4: want 1 got %s", match[4])
	}
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

func TestBuildStatistics(t *testing.T) {
	if testing.Short() {
		t.Skip("skip HTTP request in short mode")
	}
	t.Parallel()
	token, err := GetToken("kevinburke")
	if err != nil {
		t.Skip("token not found")
	}
	c := NewClient(token)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	build, err := c.Builds.Get(ctx, 366686564, "build.jobs", "job.config")
	if err != nil {
		if isHttpError(err) {
			t.Skip("no network")
		}
		t.Fatal(err)
	}
	s, err := c.BuildSummary(ctx, build)
	if err != nil {
		t.Fatal(err)
	}
	if len(s) == 0 {
		t.Errorf("zero length stats")
	}
}

func TestErrorParsing(t *testing.T) {
	token, err := GetToken("kevinburke")
	if err != nil {
		t.Skip("token not found")
	}
	c := NewClient(token)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(403)
		w.Write([]byte(`{
  "@type": "error",
  "error_type": "insufficient_access",
  "error_message": "operation requires activate access to repository",
  "resource_type": "repository",
  "permission": "activate",
  "repository": {
    "@type": "repository",
    "@href": "/repo/891",
    "@representation": "minimal",
    "id": 891,
    "name": "rails",
    "slug": "rails/rails"
  }
}`))
	}))
	defer s.Close()
	c.Client.Base = s.URL
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = c.Repos.Activate(ctx, "rails/rails")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	rerr, ok := err.(*rest.Error)
	if !ok {
		t.Errorf("err not a rest.Error: %v", err)
	}
	if rerr.Status != 403 {
		t.Errorf("bad status: want 403 got %d", rerr.Status)
	}
	want := "operation requires activate access to repository"
	if rerr.Title != want {
		t.Errorf("bad title: want %s got %s", want, rerr.Title)
	}
}

var jsonConfig = []byte(`{
        "go": "master",
        "os": "linux",
        "dist": "trusty",
        "cache": {
          "directories": [
            "$GOPATH/pkg"
          ]
        },
        "group": "stable",
        "script": "make diff race-test",
        ".result": "configured",
        "language": "go",
        "go_import_path": "github.com/kevinburke/multi-emailer"
      }`)

func TestUnmarshalConfig(t *testing.T) {
	c := new(Config)
	if err := json.Unmarshal(jsonConfig, &c); err != nil {
		t.Fatal(err)
	}
	if c.Script[0] != "make diff race-test" {
		t.Errorf("bad script: %v", c.Script)
	}
}
