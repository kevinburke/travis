package travis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

func TestParseLog(t *testing.T) {
	t.Parallel()
	steps := ParseLog(logContent)
	if len(steps) != 13 {
		t.Errorf("parseSteps: want 13 steps, got %d", len(steps))
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
		t.Fatal(err)
	}
	s, err := c.BuildStatistics(ctx, build)
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
