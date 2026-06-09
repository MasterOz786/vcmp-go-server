PLUGIN_NAME ?= goserver04rel64
PLUGIN_DIR  ?= plugins

.PHONY: all build build-linux build-windows tidy clean

all: build

tidy:
	go mod tidy

build: tidy
	@mkdir -p $(PLUGIN_DIR)
	CGO_ENABLED=1 go build -buildmode=c-shared -o $(PLUGIN_DIR)/$(PLUGIN_NAME).so .

build-linux: tidy
	@mkdir -p $(PLUGIN_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -buildmode=c-shared -o $(PLUGIN_DIR)/$(PLUGIN_NAME).so .

build-windows: tidy
	@mkdir -p $(PLUGIN_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -o $(PLUGIN_DIR)/$(PLUGIN_NAME).dll .

clean:
	rm -f $(PLUGIN_DIR)/$(PLUGIN_NAME).so $(PLUGIN_DIR)/$(PLUGIN_NAME).dll $(PLUGIN_NAME).h
