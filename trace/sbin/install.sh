#!/usr/bin/env bash
sbin="`dirname $0`"
sbin="`cd $sbin; pwd`"

source $sbin/common.sh

#Builds and installs a Go app
USAGE="
Usage: <-a application path>
\t -a relative path in \$GOPATH/src/
"

OPTIND=1         
while getopts "a:" opt; do
  case $opt in
    a)
      APP="$OPTARG"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2; usage "$USAGE"; exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2; usage "$USAGE"; exit 1
      ;;
  esac
done

if [ -z $APP ]; then
  usage "$USAGE"; exit 1
fi

echo "Building and installing $APP ..."
go install $APP

#If the app is successfully built and installed, the return code is 0
if [ $? == 0 ]; then
  echo "$APP is installed to $GOPATH/bin"
else
  echo "ERROR: failed to install $APP"
  exit 1
fi
