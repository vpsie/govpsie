# govpsie SDK — common developer tasks.
.PHONY: default build test vet fmt tidy check

default: check

build:        ## Compile all packages
	go build ./...

test:         ## Run unit tests
	go test ./... -count=1

vet:          ## Run go vet
	go vet ./...

fmt:          ## Format all Go files
	gofmt -w .

tidy:         ## Tidy the module
	go mod tidy

check: fmt vet build test  ## Format, vet, build and test
