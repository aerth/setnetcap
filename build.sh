#!/bin/bash
###############################
# setnetcap builder, by aerth #
###############################
set -e

## uncomment next line to disable calls to /usr/bin/logger
if [ "$1" == "help" ]; then
  echo "./build.sh"
  echo "./build.sh nolog"
  exit 0
fi

if [ "$1" == "sudo" ]; then
  if [ -x "setnetcap" ]; then
  set -x
  sudo chown root setnetcap && \
  sudo chmod u+s setnetcap && \
  set +x
  echo "setnetcap is armed and dangerous"
  exit 0
fi
  echo "setnetcap chown+setuid failed, you must do it manually"
  exit 111
fi

if [ "$1" == "nolog" ]; then
  echo "building with no calls to /usr/bin/logger"
  TAGS=-tags='nolog'
fi
set -x
CGO_ENABLED=0 go build -v -x ${TAGS} -ldflags='-w -s' github.com/aerth/setnetcap && \
set +x
set -e
STR='
Success:
congratulations! setnetcap build successful
*important* now run these commands as root:
  chown root setnetcap
  chmod u+s setnetcap

Tip:
	you can also use "./build.sh sudo" to do it for you

Usage info:
  * move target exe to its final destination before using setnetcap
  * this is because file capabilities may get removed when the file is modified
'; echo "$STR"; exit 0
