# edisonIsThePilot

You can find the latest version <a href="https://github.com/ssoudan/edisonIsThePilot">here</a>.

<a href="https://github.com/ssoudan">Sebastien Soudan</a> --
<a href="https://github.com/philixxx">Philippe Martinez</a>

This is intended to become an autopilot at some point.

For more information about the design check [this](DESIGN.md).

## Build the toolchain

Useful in case you are on Mac and want to build go toolchain for linux x86:

	$ cd $GOROOT/src
	$ export GOROOT_BOOTSTRAP=$GOROOT
	$ GOOS=linux GOARCH=386 ./make.bash --no-clean

## Build the project

    $ GOARCH=386 GOOS=linux go get ./...
	$ GOARCH=386 GOOS=linux go build cmd/edisonIsThePilot/edisonIsThePilot.go
	$ scp edisonIsThePilot root@edison.local.:

or you can use the Makefile:

    $ make 
    $ make deploy 

to build and copy everything to edison.local. Note this will also add a systemd service and start it.

## Licensing
Under Apache License v2.

Copyright 2015 Sebastien Soudan

Licensed under the Apache License, Version 2.0 (the "License"); you may not
use this file except in compliance with the License. You may obtain a copy
of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
License for the specific language governing permissions and limitations
under the License.