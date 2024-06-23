.DEFAULT_GOAL := build

.PHONY: build

build:
	go fmt ./...
	go build -ldflags "-w -s" ./cmd/gui
	fyne package --exe gui --icon Icon.png

clean:
	go clean
	rm -rf thermpro_exporter.app

