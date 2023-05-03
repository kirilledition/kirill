.PHONY: build run clean

build:
	@echo "Building binary..."
	$(eval commit_hash := $(shell git rev-parse HEAD))
	$(eval git_tag := $(shell git describe --tags --abbrev=0))
	go build -ldflags "-X kirill/cmd.GitCommit=$(commit_hash) -X kirill/cmd.GitTag=$(git_tag)" -o build/kirill
	@echo "Done."

run: build
	@echo "Running binary..."
	build/kirill

clean:
	@echo "Cleaning..."
	rm -f build/*
	@echo "Done."
