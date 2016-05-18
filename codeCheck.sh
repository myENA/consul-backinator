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

## check formatting ignoring git and vendor
if test $(find . -name '*.go' -not -path "./.git/*" -not -path "./vendor/*" | xargs gofmt -l -s 2>&1 | wc -l) -gt 0; then
	echo "gofmt     failed" && exit 1
else
	echo "gofmt     succeeded"
fi

## run go vet ignoring vendor and the silly "Error" bug/feature
## https://github.com/golang/go/issues/6407
if test $(go vet ./... 2>&1 | egrep -v '^vendor/|\s+vendor/|/vendor/' | grep -v ^exit\ status | grep -v "possible formatting directive in Error call" | wc -l) -gt 0; then
	echo "go vet    failed" && exit 1
else
	echo "go vet    succeeded"
fi

## run go lint ignoring vendor
if test $(${golint} ./... 2>&1 | egrep -v '^vendor/|\s+vendor/|/vendor/' | wc -l) -gt 0; then
	echo "golint    failed" && exit 1
else
	echo "golint    succeeded"
fi

## check misspell ignoring git, vendor and 3rdparty
if test $(find . -name '*' -not -path "./.git/*" -not -path "./vendor/*" -not -path "./3rdparty/*" | xargs ${misspell} 2>&1 | wc -l) -gt 0; then
	echo "misspell  failed" && exit 1
else
	echo "misspell  succeeded"
fi

## check gocyclo ignoring git and vendor
if test $(find . -name '*.go' -not -path "./.git/*" -not -path "./vendor/*" | xargs ${gocyclo} -over 15 2>&1 | wc -l) -gt 0; then
	echo "gocyclo   failed" && exit 1
else
	echo "gocyclo   succeeded"
fi
