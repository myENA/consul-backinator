#!/usr/bin/env bash

## ensure we have the golint tool
## https://github.com/golang/lint
if ! golint=$(type -p "${GOPATH}/bin/golint"); then
	echo -n "Installing golint ... "
	go get -u github.com/golang/lint/golint
	echo "done"
	golint=$(type -p "${GOPATH}/bin/golint")
fi

## ensure we have the misspell tool
## https://github.com/client9/misspell
if ! misspell=$(type -p "${GOPATH}/bin/misspell"); then
	echo -n "Installing misspell ... "
	go get -u github.com/client9/misspell/cmd/misspell
	echo "done"
	misspell=$(type -p "${GOPATH}/bin/misspell")
fi

## ensure we have the gocyclo tool
## https://github.com/fzipp/gocyclo
if ! gocyclo=$(type -p "${GOPATH}/bin/gocyclo"); then
	echo -n "Installing gocyclo ... "
	go get -u github.com/fzipp/gocyclo
	echo "done"
	gocyclo=$(type -p "${GOPATH}/bin/gocyclo")
fi

## check formatting ignoring git and vendor
fmtTest=$(find . -name '*.go' -not -path './.git/*' -not -path './vendor/*' | xargs gofmt -l -s 2>&1)
if [ ! -z "$fmtTest" ]; then
	echo "gofmt     failed"
	echo "$fmtTest"
	exit 1
else
	echo "gofmt     succeeded"
fi

## run go vet ignoring vendor and the silly "Error" bug/feature
## https://github.com/golang/go/issues/6407
vetTest=$(go vet ./... 2>&1 | egrep -v '^vendor/|\s+vendor/|/vendor/|^exit\ status|\ possible\ formatting\ directive\ in\ Error\ call')
if [ ! -z "$vetTest" ]; then
	echo "go vet    failed"
	echo "$vetTest"
	exit 1
else
	echo "go vet    succeeded"
fi

## run go lint ignoring vendor
lintTest=$(${golint} ./... 2>&1 | egrep -v '^vendor/|\s+vendor/|/vendor/')
if [ ! -z "$lintTest" ]; then
	echo "golint    failed"
	echo "$lintTest"
	exit 1
else
	echo "golint    succeeded"
fi

## check misspell ignoring git, vendor and 3rdparty
spellTest=$(find . -name '*' -not -path './.git/*' -not -path './vendor/*' -not -path './3rdparty/*' | xargs ${misspell} 2>&1 | echo)
if [ ! -z "$spellTest" ]; then
	echo "misspell  failed"
	echo "$spellTest"
	exit 1
else
	echo "misspell  succeeded"
fi

## check gocyclo ignoring git and vendor
cycloTest=$(find . -name '*.go' -not -path './.git/*' -not -path './vendor/*' | xargs ${gocyclo} -over 15 2>&1 | echo)
if [ ! -z "$cycloTest" ]; then
	echo "gocyclo   failed"
	echo "$cycloTest"
	exit 1
else
	echo "gocyclo   succeeded"
fi
