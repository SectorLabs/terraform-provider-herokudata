GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PROVIDER_NAME="terraform-provider-herokudata"
PROVIDER_VERSION="0.1"

default: build

build: fmtcheck
	go build -o $(PROVIDER_NAME)

release:
	go build -o "$(PROVIDER_NAME)_v$(PROVIDER_VERSION)"

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: build fmt fmtcheck errcheck
