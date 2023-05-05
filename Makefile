.PHONY: build run clean

build:
	@echo "Building binary..."
	$(eval commit_hash := $(shell git rev-parse HEAD))
	$(eval git_tag := $(shell git describe --tags --abbrev=0))
	mkdir -p build
	go build -ldflags "-X kirill/cmd.GitCommit=$(commit_hash) -X kirill/cmd.GitTag=$(git_tag)" -o build/kirill
	GOOS=linux GOARCH=amd64 go build -ldflags "-X kirill/cmd.GitCommit=$(commit_hash) -X kirill/cmd.GitTag=$(git_tag)" -o build/kirill_linux
	@echo "Done."

run: build
	@echo "Running binary..."
	build/kirill

clean:
	@echo "Cleaning..."
	rm -rf build
	@echo "Done."
