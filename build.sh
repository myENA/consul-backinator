#!/usr/bin/env bash
#
## package declarations
BUILD_NAME="consul-backinator"
RELEASE_VERSION="1.0"
RELEASE_BUILD=0

## required for glide
export GO15VENDOREXPERIMENT=1

## simple usage example
showUsage() {
	printf "Usage: $0 [-u|-i|-d]
	-u    Update vendor directory from glide.yaml using 'glide up' and build
	-i    Install vendor directory from glide.lock using 'glide install' and build
	-d    Remove existing glide.lock and vendor directory and exit
	-r    Build and package release binaries\n\n"
	exit 0
}

## read options
while getopts ":uidr" opt; do
	case $opt in
		u)
			printf "Updating vendor directory ... "
			glide -q up > /dev/null 2>&1
		;;
		i)
			printf "Installing from glide.lock ..."
			glide -q install > /dev/null 2>&1
		;;
		d)
			printf "Removing glide.lock and vendor directory ...\n"
			rm -rf glide.lock vendor
			exit 0
		;;
		r)
			## toggle release building
			RELEASE_BUILD=1
		;;
		*)
			showUsage
		;;
	esac
done

## remove options
shift $((OPTIND-1))

## check release option
if [[ $RELEASE_BUILD -eq 1 ]]; then
	## clean dist directory
	rm -rf ./dist/

	## check for gox
	if [[ ! -x $(which gox) ]]; then
		printf "Installing gox ... "
		go get github.com/mitchellh/gox
	fi

	## build release
	printf "Building release ... \n"

	## call gox to build our binaries
	gox \
	-osarch="linux/amd64 darwin/amd64 freebsd/amd64 openbsd/amd64 windows/386 windows/amd64" \
	-ldflags="-X main.appVersion=${RELEASE_VERSION}" \
	-output="./dist/${BUILD_NAME}-${RELEASE_VERSION}-{{.Arch}}-{{.OS}}/${BUILD_NAME}-${RELEASE_VERSION}"

	## gox return
	RETVAL=$?

else

	## build binaries
	printf "Building ... "

	## build it
	go build -o "${BUILD_NAME}"

	## go build return
	RETVAL=$?
fi

## check build status
if [ $RETVAL -ne 0 ]; then
	printf "\nError during build!\n"
	exit $RETVAL
fi

## check release option
if [[ $RELEASE_BUILD -eq 1 ]]; then
	## package binaries
	printf "Packaging ... \n"

	## package files
	pushd ./dist/ > /dev/null 2>&1
	find . -type d -name \*-\* -maxdepth 1 \
	-exec zip -m -r {}.zip {}/ \;
	popd > /dev/null 2>&1

	## all done
	printf "done\n"
else
	## all done
	printf "done.\nUsage: ./${BUILD_NAME} -h\n"
fi

## exit same as build
exit $RETVAL
