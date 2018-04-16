BUMP_VERSION := $(GOPATH)/bin/bump_version
RELEASE := $(GOPATH)/bin/github-release

test:
	go test ./...

race-test:
	go test -race ./...

$(BUMP_VERSION):
	go get -u github.com/kevinburke/bump_version

$(RELEASE):
	go get -u github.com/aktau/github-release

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
