SRC := $(shell find . -type f -name '*.go')
CROSSARCH := amd64 386
CROSSOS := darwin linux windows
TIMETRACKING_CROSS := $(foreach os,$(CROSSOS),$(foreach arch,$(CROSSARCH),dist/timetracking/$(os).$(arch)))

.PHONY: clean build cross install

build: dist/timetracking/native
cross: $(TIMETRACKING_CROSS)
install: build
	mv dist/timetracking/native $(GOBIN)/timetracking
	

clean:
	rm -r dist

dist/timetracking/native: $(SRC)
	go build -ldflags="-X main.v=$(shell git describe)" -o $@ cmd/timetracking/*.go

$(TIMETRACKING_CROSS): $(SRC)
	gox \
		-osarch="$(shell echo "$@" | cut -d'/' -f3- | sed 's/\./\//')" \
		-ldflags="-X main.v=$(shell git describe)" \
		-output="$@" \
		./cmd/timetracking
	if [ -f "$@.exe" ]; then mv "$@.exe" "$@"; fi
