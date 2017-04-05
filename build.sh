#!/bin/sh
###############################
# setnetcap builder, by aerth #
###############################
set -e
CGO_ENABLED=0 go build -ldflags='-w -s' github.com/aerth/setnetcap && \
echo 'setnetcap build successful' && \
echo '*important* now run these commands as root:' && \
echo '' && \
echo '  chown root setnetcap' && \
echo '  chmod u+s setnetcap' && \
echo '' && \
echo '*important*: move setnetcap to its final destination before running these commands'
exit 0
