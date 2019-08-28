#!/usr/bin/env bash
#
## package declarations
BUILD_NAME="consul-backinator"
RELEASE_VERSION="1.6.5"
RELEASE_BUILD=0

## simple usage example
showUsage() {
	printf "Usage: $0 [-d|-c|-r]
	-d    Remove binary and vendor directory
	-c    Clean distribution files
	-r    Build and package release binaries\n\n"
	exit 0
}

## install gox if needed
ensureGox() {
	which gox > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		printf "Installing gox ... "
		go get github.com/mitchellh/gox
	fi
}

## we require module support
export GO111MODULE=on

## exit toggle
should_exit=false

## read options
while getopts ":dcr" opt; do
	case $opt in
		d)
			printf "Removing binary and vendor directory ... "
			rm -rf "${BUILD_NAME}" vendor
			printf "done.\n"
			should_exit=true
		;;
		c)
			printf "Cleaning dist directory ... "
			rm -rf ./dist/
			printf "done.\n"
			should_exit=true
		;;
		r)
			ensureGox
			RELEASE_BUILD=1
		;;
		*)
			showUsage
		;;
	esac
done

## remove options
shift $((OPTIND-1))

## exiting?
if [ $should_exit == true ]; then
	exit 0
fi

## check release option
if [ $RELEASE_BUILD -eq 1 ]; then
	## clean dist directory
	rm -rf ./dist/

	## call gox to build our binaries
	CGO_ENABLED=0 gox \
	-osarch="linux/amd64 darwin/amd64 freebsd/amd64 windows/amd64 windows/386" \
	-ldflags="-X main.appVersion=${RELEASE_VERSION} -s -w" \
	-output="./dist/${BUILD_NAME}-${RELEASE_VERSION}-{{.Arch}}-{{.OS}}/${BUILD_NAME}-${RELEASE_VERSION}"

	## gox return
	RETVAL=$?

else

	## build it
	CGO_ENABLED=0 go build -o "${BUILD_NAME}" \
	-ldflags="-X main.appVersion=${RELEASE_VERSION} -s -w"

	## go build return
	RETVAL=$?
fi

## check build status
if [ $RETVAL -ne 0 ]; then
	printf "\nError during build!\n"
	exit $RETVAL
fi

## check release option
if [ $RELEASE_BUILD -eq 1 ]; then
	## package binaries
	printf "\nPackaging ... "

	## package files
	pushd ./dist/ > /dev/null >&1
	find . -maxdepth 1 -type d -name \*-\* \
	-exec tar -czf {}.tar.gz {} > /dev/null >&1 \; \
	-exec zip -m -r {}.zip {} > /dev/null >&1 \;

	## package binaries
	printf "signing ... "

	## generate checksums and sign
	shasum -a256 *.tar.gz *.zip >> ${BUILD_NAME}-${RELEASE_VERSION}-SHA256SUMS
	gpg2 -u "r&d@ena.com" -b ${BUILD_NAME}-${RELEASE_VERSION}-SHA256SUMS > /dev/null >&1
	popd > /dev/null >&1

	## all done
	printf "done.\nRelease files may be found in the ./dist/ directory.\n"
else
	## all done
	printf "Usage: ./${BUILD_NAME} -h\n"
fi

## exit same as build
exit $RETVAL
