package travis

import (
	"context"
	"encoding/json"
	"testing"
	"time"
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
	j := new(Job)
	if err := json.Unmarshal(jobResponse, j); err != nil {
		t.Fatal(err)
	}
}

func TestParseLog(t *testing.T) {
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
