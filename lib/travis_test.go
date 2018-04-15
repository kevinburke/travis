package travis

import (
	"encoding/json"
	"testing"
)

func TestMarshalBuilds(t *testing.T) {
	r := new(Response)
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
