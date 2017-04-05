#!/bin/sh
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
  echo running sudo chown root setnetcap && \
  sudo chown root setnetcap && \
  echo running sudo chmod u+s setnetcap && \
  sudo chmod u+s setnetcap && \
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

CGO_ENABLED=0 go build ${TAGS} -ldflags='-w -s' github.com/aerth/setnetcap && \
echo 'congratulations! setnetcap build successful' && \
echo '*important* now run these commands as root:' && \
echo '' && \
echo '  chown root setnetcap' && \
echo '  chmod u+s setnetcap' && \
echo '' && \
echo '*important*: move setnetcap to its final destination before running these commands'
echo "you can also use './build.sh sudo' to do it for you"
exit 0
