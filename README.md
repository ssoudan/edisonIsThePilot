# edisonIsThePilot

Sebastien Soudan

This is intended to become an autopilot at some point.

## Build the toolchain

	$ cd $GOROOT/src
	$ export GOROOT_BOOTSTRAP=$GOROOT
	$ GOOS=linux GOARCH=386 ./make.bash --no-clean

## Build the project

    $ GOARCH=386 GOOS=linux go get
	$ GOARCH=386 GOOS=linux go build edisonIsThePilot.go
	$ scp edisonIsThePilot root@edison.local.:%