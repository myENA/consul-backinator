#!/usr/bin/env bash
set -e

## ensure we have the golint tool
## https://github.com/golang/lint
if ! golint="$(type -p "${GOPATH}/bin/golint")"; then
	echo -n "Installing golint ... "
	go get -u github.com/golang/lint/golint
	echo "done"
	golint="$(type -p "${GOPATH}/bin/golint")"
fi

## ensure we have the misspell tool
## https://github.com/client9/misspell
if ! misspell="$(type -p "${GOPATH}/bin/misspell")"; then
	echo -n "Installing misspell ... "
	go get -u github.com/client9/misspell/cmd/misspell
	echo "done"
	misspell="$(type -p "${GOPATH}/bin/misspell")"
fi

## ensure we have the gocyclo tool
## https://github.com/fzipp/gocyclo
if ! gocyclo="$(type -p "${GOPATH}/bin/gocyclo")"; then
	echo -n "Installing gocyclo ... "
	go get -u github.com/fzipp/gocyclo
	echo "done"
	gocyclo="$(type -p "${GOPATH}/bin/gocyclo")"
fi

## run the tests ignoring vendor and git directories where needed
test $(find . -name '*.go' -not -path "./.git/*" -not -path "./vendor/*" | xargs gofmt -l -s 2>&1 | wc -l) -gt 0         && echo "gofmt     failed" && exit 1 || echo "gofmt     succeeded"
test $(go vet ./... 2>&1 | grep -v ^vendor | grep -v ^exit\ status| wc -l) -gt 0                                         && echo "go vet    failed" && exit 1 || echo "go vet    succeeded"
test $(${golint} ./... 2>&1 | grep -v ^vendor | wc -l) -gt 0                                                             && echo "golint    failed" && exit 1 || echo "golint    succeeded"
test $(find . -name '*'    -not -path "./.git/*" -not -path "./vendor/*" | xargs ${misspell} 2>&1 | wc -l) -gt 0         && echo "misspell  failed" && exit 1 || echo "misspell  succeeded"
test $(find . -name '*.go' -not -path "./.git/*" -not -path "./vendor/*" | xargs ${gocyclo} -over 15 2>&1 | wc -l) -gt 0 && echo "gocyclo   failed" && exit 1 || echo "gocyclo   succeeded"
