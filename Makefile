SHELL = /bin/bash -o pipefail

BENCHSTAT := $(GOPATH)/bin/benchstat
BUMP_VERSION := $(GOPATH)/bin/bump_version
MEGACHECK := $(GOPATH)/bin/megacheck
RELEASE := $(GOPATH)/bin/github-release
UNAME = $(shell uname -s)

test:
	go test ./...

race-test: lint
	go test -race ./...

lint: | $(MEGACHECK)
	go vet ./...
	go list ./... | grep -v vendor | xargs $(MEGACHECK)

bench: | $(BENCHSTAT)
	go list ./... | grep -v vendor | xargs go test -benchtime=2s -bench=. -run='^$$' 2>&1 | $(BENCHSTAT) /dev/stdin

$(BUMP_VERSION):
	go get -u github.com/kevinburke/bump_version

$(BENCHSTAT):
	go get golang.org/x/perf/cmd/benchstat

$(RELEASE):
	go get -u github.com/aktau/github-release

$(GOPATH)/bin:
	mkdir -p $(GOPATH)/bin

$(MEGACHECK): | $(GOPATH)/bin
ifeq ($(UNAME),Darwin)
	curl --silent --location --output $(MEGACHECK) https://github.com/kevinburke/go-tools/releases/download/2018-11-18/megacheck-darwin-amd64
else
	curl --silent --location --output $(MEGACHECK) https://github.com/kevinburke/go-tools/releases/download/2018-11-18/megacheck-linux-amd64
endif
	chmod +x $(MEGACHECK)

release: race-test | $(BUMP_VERSION) $(RELEASE)
ifndef version
	@echo "Please provide a version"
	exit 1
endif
ifndef GITHUB_TOKEN
	@echo "Please set GITHUB_TOKEN in the environment"
	exit 1
endif
	$(BUMP_VERSION) --version=$(version) lib/travis.go
	git push origin --tags
	mkdir -p releases/$(version)
	GOOS=linux GOARCH=amd64 go build -o releases/$(version)/travis-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o releases/$(version)/travis-darwin-amd64 .
	GOOS=windows GOARCH=amd64 go build -o releases/$(version)/travis-windows-amd64 .
	# These commands are not idempotent, so ignore failures if an upload repeats
	$(RELEASE) release --user kevinburke --repo travis --tag $(version) || true
	$(RELEASE) upload --user kevinburke --repo travis --tag $(version) --name travis-linux-amd64 --file releases/$(version)/travis-linux-amd64 || true
	$(RELEASE) upload --user kevinburke --repo travis --tag $(version) --name travis-darwin-amd64 --file releases/$(version)/travis-darwin-amd64 || true
	$(RELEASE) upload --user kevinburke --repo travis --tag $(version) --name travis-windows-amd64 --file releases/$(version)/travis-windows-amd64 || true
