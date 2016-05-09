#!/usr/bin/env bash
#
## package declarations
BUILD_NAME="consul-backinator"

## required for glide
export GO15VENDOREXPERIMENT=1

## simple usage example
showUsage() {
	printf "Usage: $0 [-u|-i|-d]
	-u    Update vendor directory from glide.yaml using 'glide up' and build
	-i    Install vendor directory from glide.lock using 'glide install' and build
	-d    Remove existing glide.lock and vendor directory and exit\n\n"
	exit 0
}

## read options
while getopts ":uid" opt; do
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
		*)
			showUsage
		;;
	esac
done

## remove options
shift $((OPTIND-1))

## continue to build
printf "Building ... "

## build it
go build -o "${BUILD_NAME}"

## go build return
RETVAL=$?

## check status
if [ $RETVAL -ne 0 ]; then
	printf "Error during build!\n"
else
	printf "done.\nUsage: ./${BUILD_NAME} -h\n"
fi

## exit same as build
exit $RETVAL
