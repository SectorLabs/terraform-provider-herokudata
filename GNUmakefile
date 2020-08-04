GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PROVIDER_NAME="terraform-provider-herokudata"

default: build

build: fmtcheck
	go build -o $(PROVIDER_NAME)

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

.PHONY: build fmt fmtcheck errcheck
