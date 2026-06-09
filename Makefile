PLUGIN_NAME ?= goserver04rel64
HEADER := include/plugin.h

.PHONY: all deps build build-linux build-windows clean

all: build

deps:
	cd scripts && go run fetch_plugin.go

$(HEADER):
	$(MAKE) deps

build: $(HEADER)
	CGO_ENABLED=1 go build -buildmode=c-shared -o $(PLUGIN_NAME).so .

build-linux: $(HEADER)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o $(PLUGIN_NAME).so .

build-windows: $(HEADER)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -o $(PLUGIN_NAME).dll .

clean:
	rm -f $(PLUGIN_NAME).so $(PLUGIN_NAME).dll $(PLUGIN_NAME).h
