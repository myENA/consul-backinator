#!/usr/bin/env bash

## ensure we have the misspell tool
## https://github.com/client9/misspell
if ! misspell=$(type -p "${GOPATH}/bin/misspell"); then
	echo -n "Installing misspell ... "
	go install github.com/client9/misspell/cmd/misspell@latest
	echo "done"
	misspell=$(type -p "${GOPATH}/bin/misspell")
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

if  ! go vet "./..." ; then
	echo "go vet    failed"
	echo "output - ---"
	echo $vetTest
	echo "----"
	exit 1
else
	echo "go vet    succeeded"
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

