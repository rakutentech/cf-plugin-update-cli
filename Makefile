COMMIT = $$(git describe --always)

default: run

run: uninstall install
	DEBUG_PLUGIN=1 cf update-cli

deps:
	go get -u -v .

build: deps
	go build -ldflags "-X main.GitCommit=$(COMMIT)" -o bin/cf-plugin-update-cli

install: build
	cf install-plugin bin/cf-plugin-update-cli -f
	cf plugins

uninstall:
	cf uninstall-plugin 'update-cli'

test: vet 
	go test -v

vet:
	@go get golang.org/x/tools/cmd/vet
	go tool vet *.go

lint:
	@go get github.com/golang/lint/golint
	golint ./...

# cover shows test coverages
cover:
	@go get golang.org/x/tools/cmd/cover		
	godep go test -coverprofile=cover.out
	go tool cover -html cover.out
	rm cover.out
