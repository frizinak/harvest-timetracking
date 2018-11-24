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
	go build -o $@ cmd/timetracking/*.go

dist: $(SRC)
	gox \
		-osarch="$(CROSS)" \
		-output="dist/{{.OS}}.{{.Arch}}" \
		./cmd/slek

	touch dist

$(TIMETRACKING_CROSS): $(SRC)
	gox \
		-osarch="$(shell echo "$@" | cut -d'/' -f3- | sed 's/\./\//')" \
		-output="$@" \
		./cmd/timetracking
	if [ -f "$@.exe" ]; then mv "$@.exe" "$@"; fi
