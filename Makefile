# Makefile for a go project
#
# Author: Jon Eisen
# 	site: joneisen.me
# 	
# Targets:
# 	all: Builds the code
# 	build: Builds the code
# 	fmt: Formats the source files
# 	clean: cleans the code
# 	install: Installs the code to the GOPATH
# 	iref: Installs referenced projects
#	test: Runs the tests
#	
#  Blog post on it: http://joneisen.me/post/25503842796
#
FLAGS=GOARCH=386 GOOS=linux
GITID=$(shell git rev-parse --short HEAD)

SSH=ssh -o StrictHostKeyChecking=no root@edison.local.
SCP=scp -o StrictHostKeyChecking=no

# Go parameters
GOCMD=$(FLAGS) go
GOBUILD=$(GOCMD) build -ldflags "-X main.Version='$(GITID)'" 
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=go test -cover
GODEP=$(GOCMD) test -i
GOFMT=gofmt -w

# Package lists
TOPLEVEL_PKG := github.com/ssoudan/edisonIsThePilot
INT_LIST :=  #<-- Interface directories
IMPL_LIST := conf control alarm dashboard pilot gps steering stepper drivers/pwm drivers/mcp4725 drivers/sincos drivers/gpio drivers/motor  infrastructure/logger infrastructure/pid  #<-- Implementation directories
CMD_LIST := cmd/edisonIsThePilot cmd/mario cmd/ap100Control cmd/systemCalibration cmd/motorControl cmd/ledControl cmd/motorCalibration cmd/alarmControl #<-- Command directories

# List building
ALL_LIST = $(INT_LIST) $(IMPL_LIST) $(CMD_LIST)

BUILD_LIST = $(foreach int, $(ALL_LIST), $(int)_build)
CLEAN_LIST = $(foreach int, $(ALL_LIST), $(int)_clean)
INSTALL_LIST = $(foreach int, $(ALL_LIST), $(int)_install)
IREF_LIST = $(foreach int, $(ALL_LIST), $(int)_iref)
TEST_LIST = $(foreach int, $(ALL_LIST), $(int)_test)
FMT_LIST = $(foreach int, $(ALL_LIST), $(int)_fmt)

# All are .PHONY for now because dependencyness is hard
.PHONY: $(CLEAN_LIST) $(TEST_LIST) $(FMT_LIST) $(INSTALL_LIST) $(BUILD_LIST) $(IREF_LIST)

all: build
build: $(BUILD_LIST)
clean: $(CLEAN_LIST)
install: $(INSTALL_LIST)
test: $(TEST_LIST)
iref: $(IREF_LIST)
fmt: $(FMT_LIST)
deploy: build test
	$(SSH) systemctl stop edisonIsThePilot
	sleep 3
	$(SCP) mario edisonIsThePilot motorControl systemCalibration alarmControl ledControl motorCalibration edisonIsThePilot.service root@edison.local.:
	$(SSH) cp edisonIsThePilot.service /lib/systemd/system
	$(SSH) systemctl daemon-reload
	$(SSH) systemctl start edisonIsThePilot

$(BUILD_LIST): %_build: %_fmt %_iref
	$(GOBUILD) $(TOPLEVEL_PKG)/$*
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(TOPLEVEL_PKG)/$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*
$(IREF_LIST): %_iref:
	$(GODEP) $(TOPLEVEL_PKG)/$*
$(TEST_LIST): %_test:
	$(GOTEST) $(TOPLEVEL_PKG)/$*
$(FMT_TEST): %_fmt:
	$(GOFMT) ./$*
