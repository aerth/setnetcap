# setnetcap

simple little static binary that you can setuid root and allow daemons to setcap themselves

perfect for **automated build scripts** for updating **live running web servers**

## installation

step 1, as any user:

```

CGO_ENABLED=0 go get -v -ldflags='-s' github.com/aerth/setnetcap

```

step 2, as root, changing '/path/to' to the location setnetcap will live

```
mv setnetcap /path/to/ && \
chown root /path/to/setnetcap && \
chmod u+s /path/to/setnetcap && \
echo "setnetcap is owned by root, setuid, and can provide additional net bind abilities"

```

Now, any user with access to 'setnetcap' can serve on privileged ports (such as 80, 443)

### use case 1

  * linux machine, with /sbin/setcap installed
  * unprivileged daemon user running web server which is listening on port 80, 443
  * web server gets triggered to rebuild itself (by webhook, cron, whatever)

Problem: after rebuilding server daemon binary, restarting will result in an error, "permission denied" without running /sbin/setcap as root.

Solution: with setnetcap installed correctly, it can be used in the build script such as:

```
#!/bin/sh
go build -o serverd
setnetcap serverd && echo "setnetcap works"
```

### security thoughts

  * any user with access to setnetcap can use it, maybe hide it somewhere with no read access
  * low ports (X, syslog, auth, time) potentially be hijacked
  * better than sudo? i dont know

### credits

Copyright (c) 2016-2017 aerth <aerth@riseup.net> (MIT License)
