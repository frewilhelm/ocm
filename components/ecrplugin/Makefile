NAME      = ecrplugin
PROVIDER  ?= ocm.software
GITHUBORG  ?= open-component-model
COMPONENT = $(PROVIDER)/plugins/$(NAME)
OCMREPO   ?= ghcr.io/$(GITHUBORG)/ocm
PLATFORMS ?= linux/amd64 linux/arm64 darwin/amd64 darwin/arm64
CTF_TYPE ?= directory


REPO_ROOT                                      := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/../..
ifeq ($(VERSION),)
VERSION                                        := $(shell go run ../../api/version/generate/release_generate.go print-rc-version $(CANDIDATE))
endif
COMMIT                                         = $(shell git rev-parse --verify HEAD)
# if EFFECTIVE_VERSION is not set, set it to VERSION+COMMIT
# this is not the same as '?=' because it will also set the value if EFFECTIVE_VERSION is set to an empty string
ifeq ($(EFFECTIVE_VERSION),)
EFFECTIVE_VERSION                              := $(VERSION)+$(COMMIT)
endif
GIT_TREE_STATE                                 := $(shell [ -z "$$(git status --porcelain 2>/dev/null)" ] && echo clean || echo dirty)

CMDSRCS=$(shell find $(REPO_ROOT)/cmds/$(NAME) -type f)
OCMSRCS=$(shell find $(REPO_ROOT)/api -type f) $(REPO_ROOT)/go.*

CREDS ?=
# Define the path to the binary
OCM_BIN = $(REPO_ROOT)/bin/ocm

# Rule to build the binary if it doesn't exist or if the source code has changed
$(OCM_BIN): $(REPO_ROOT)/cmds/ocm/main.go
	mkdir -p $(REPO_ROOT)/bin
	go build -ldflags $(BUILD_FLAGS) -o $(OCM_BIN) $(REPO_ROOT)/cmds/ocm

# Use the binary for the OCM command
OCM = $(OCM_BIN) $(CREDS)

GEN = $(REPO_ROOT)/gen/$(NAME)

 $(GEN):
	@mkdir -p $(GEN)

NOW         := $(shell date -u +%FT%T%z)
BUILD_FLAGS := "-s -w \
 -X ocm.software/ocm/api/version.gitVersion=$(EFFECTIVE_VERSION) \
 -X ocm.software/ocm/api/version.gitTreeState=$(GIT_TREE_STATE) \
 -X ocm.software/ocm/api/version.gitCommit=$(COMMIT) \
 -X ocm.software/ocm/api/version.buildDate=$(NOW)"

.PHONY: build
build: $(GEN)/build

$(GEN)/build: $(GEN) $(CMDSRCS) $(OCMSRCS)
	@for i in $(PLATFORMS); do \
    tag=$$(echo $$i | sed -e s:/:-:g); \
    echo GOARCH=$$(basename $$i) GOOS=$$(dirname $$i) CGO_ENABLED=0 go build -ldflags $(BUILD_FLAGS) -o $(GEN)/$(NAME).$$tag ../../cmds/$(NAME); \
    GOARCH=$$(basename $$i) GOOS=$$(dirname $$i) CGO_ENABLED=0 go build -ldflags $(BUILD_FLAGS) -o $(GEN)/$(NAME).$$tag ../../cmds/$(NAME) & \
    done; \
	wait
	@touch $(GEN)/build


.PHONY: ctf
ctf: $(GEN)/ctf

$(GEN)/ctf: $(OCM_BIN) $(GEN)/.exists $(GEN)/build component-constructor.yaml $(CHARTSRCS)
	@rm -rf "$(GEN)/ctf"
	$(OCM) add componentversions \
		--create \
		--file $(GEN)/ctf \
		--type $(CTF_TYPE) \
		--templater=spiff \
		COMPONENT="$(COMPONENT)" \
		NAME="$(NAME)" \
		VERSION="$(VERSION)" \
		PROVIDER="$(PROVIDER)" \
		COMMIT="$(COMMIT)" \
		GEN="$(GEN)" \
		PLATFORMS="$(PLATFORMS)" \
		component-constructor.yaml
	touch "$(GEN)/ctf"

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: push
push: $(GEN)/ctf $(GEN)/push.$(NAME)

$(GEN)/push.$(NAME): $(GEN)/ctf $(OCM_BIN)
	$(OCM) transfer ctf -f $(GEN)/ctf $(OCMREPO)
	@touch $(GEN)/push.$(NAME)

.PHONY: plain-push
plain-push: $(GEN) $(OCM_BIN)
	$(OCM) transfer ctf -f $(GEN)/ctf $(OCMREPO)
	@touch $(GEN)/push.$(NAME)

.PHONY: transport
transport: $(OCM_BIN)
ifneq ($(TARGETREPO),)
	$(OCM) transfer component -Vc  $(OCMREPO)//$(COMPONENT):$(VERSION) $(TARGETREPO)
endif

$(GEN)/.exists:
	@mkdir -p $(GEN)
	@touch $@

.PHONY: info
info:
	@echo "ROOT:     $(REPO_ROOT)"
	@echo "VERSION:  $(VERSION)"
	@echo "COMMIT;   $(COMMIT)"

.PHONY: describe
describe: $(GEN)/ctf $(OCM_BIN)
	$(OCM) get resources --lookup $(OCMREPO) -r -o treewide $(GEN)/ctf

.PHONY: descriptor
descriptor: $(GEN)/ctf $(OCM_BIN)
	$(OCM) get component -S v3alpha1 -o yaml $(GEN)/ctf

.PHONY: clean
clean:
	rm -rf $(GEN)
