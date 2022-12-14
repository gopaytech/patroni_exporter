ENVVAR=CGO_ENABLED=0 GO111MODULE=on
PKGS=./...

GO?=go
GOTEST?=$(GO) test
GOFMT?=$(GO) fmt
GOOS?=$(shell $(GO) env GOHOSTOS)
GOARCH?=$(shell $(GO) env GOHOSTARCH)

ifdef LDFLAGS
  LDFLAGS_FLAG=--ldflags "${LDFLAGS}"
else
  LDFLAGS_FLAG=
endif

.PHONY: build
build:
	mkdir -p .build/${GOOS}-${GOARCH}
	$(ENVVAR) GOOS=$(GOOS) go build -o .build/${GOOS}-${GOARCH}/patroni_exporter ${LDFLAGS_FLAG} ./cmd/main.go

.PHONY: test
test:
	$(GOTEST) $(PKGS)

.PHONY: fmt
fmt:
	$(GOFMT) $(PKGS)

.PHONY: run
run:
	.build/${GOOS}-${GOARCH}/patroni_exporter --patroni.host="http://localhost" --patroni.port=8008
