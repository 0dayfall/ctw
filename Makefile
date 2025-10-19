.PHONY: build clean install test deb

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.1.0")
BINARY_NAME = ctw
BUILD_DIR = build
DEB_DIR = $(BUILD_DIR)/deb
INSTALL_PREFIX ?= /usr/local

build:
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/ctw

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PREFIX)/bin..."
	@install -d $(INSTALL_PREFIX)/bin
	@install -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PREFIX)/bin/$(BINARY_NAME)

test:
	@echo "Running tests..."
	go test ./...

# Build .deb package for Ubuntu/Debian
deb: build
	@echo "Building .deb package version $(VERSION)..."
	@mkdir -p $(DEB_DIR)/DEBIAN
	@mkdir -p $(DEB_DIR)/usr/bin
	@mkdir -p $(DEB_DIR)/usr/share/doc/$(BINARY_NAME)
	
	# Copy binary
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(DEB_DIR)/usr/bin/
	
	# Copy documentation
	@cp README.md $(DEB_DIR)/usr/share/doc/$(BINARY_NAME)/
	@cp LICENSE $(DEB_DIR)/usr/share/doc/$(BINARY_NAME)/
	
	# Create control file
	@echo "Package: $(BINARY_NAME)" > $(DEB_DIR)/DEBIAN/control
	@echo "Version: $(VERSION)" >> $(DEB_DIR)/DEBIAN/control
	@echo "Section: utils" >> $(DEB_DIR)/DEBIAN/control
	@echo "Priority: optional" >> $(DEB_DIR)/DEBIAN/control
	@echo "Architecture: amd64" >> $(DEB_DIR)/DEBIAN/control
	@echo "Maintainer: 0dayfall <maintainer@example.com>" >> $(DEB_DIR)/DEBIAN/control
	@echo "Description: Command-line toolkit for Twitter v2 API" >> $(DEB_DIR)/DEBIAN/control
	@echo " A Go-based CLI tool for working with Twitter v2 REST endpoints." >> $(DEB_DIR)/DEBIAN/control
	@echo " Supports tweets, users, DMs, media upload, and more." >> $(DEB_DIR)/DEBIAN/control
	
	# Build the package
	@dpkg-deb --build $(DEB_DIR) $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_amd64.deb
	@echo "Package created: $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_amd64.deb"

# Cross-compile for different architectures
deb-multi: build-amd64 build-arm64
	@echo "Building multi-architecture .deb packages..."

build-amd64:
	@echo "Building for amd64..."
	@mkdir -p $(BUILD_DIR)
	GOARCH=amd64 GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-amd64 ./cmd/ctw
	@$(MAKE) deb-arch ARCH=amd64

build-arm64:
	@echo "Building for arm64..."
	@mkdir -p $(BUILD_DIR)
	GOARCH=arm64 GOOS=linux go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-arm64 ./cmd/ctw
	@$(MAKE) deb-arch ARCH=arm64

deb-arch:
	@mkdir -p $(DEB_DIR)-$(ARCH)/DEBIAN
	@mkdir -p $(DEB_DIR)-$(ARCH)/usr/bin
	@mkdir -p $(DEB_DIR)-$(ARCH)/usr/share/doc/$(BINARY_NAME)
	
	@cp $(BUILD_DIR)/$(BINARY_NAME)-$(ARCH) $(DEB_DIR)-$(ARCH)/usr/bin/$(BINARY_NAME)
	@cp README.md $(DEB_DIR)-$(ARCH)/usr/share/doc/$(BINARY_NAME)/
	@cp LICENSE $(DEB_DIR)-$(ARCH)/usr/share/doc/$(BINARY_NAME)/
	
	@echo "Package: $(BINARY_NAME)" > $(DEB_DIR)-$(ARCH)/DEBIAN/control
	@echo "Version: $(VERSION)" >> $(DEB_DIR)-$(ARCH)/DEBIAN/control
	@echo "Section: utils" >> $(DEB_DIR)-$(ARCH)/DEBIAN/control
	@echo "Priority: optional" >> $(DEB_DIR)-$(ARCH)/DEBIAN/control
	@echo "Architecture: $(ARCH)" >> $(DEB_DIR)-$(ARCH)/DEBIAN/control
	@echo "Maintainer: 0dayfall <maintainer@example.com>" >> $(DEB_DIR)-$(ARCH)/DEBIAN/control
	@echo "Description: Command-line toolkit for Twitter v2 API" >> $(DEB_DIR)-$(ARCH)/DEBIAN/control
	@echo " A Go-based CLI tool for working with Twitter v2 REST endpoints." >> $(DEB_DIR)-$(ARCH)/DEBIAN/control
	@echo " Supports tweets, users, DMs, media upload, and more." >> $(DEB_DIR)-$(ARCH)/DEBIAN/control
	
	@dpkg-deb --build $(DEB_DIR)-$(ARCH) $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_$(ARCH).deb
	@echo "Package created: $(BUILD_DIR)/$(BINARY_NAME)_$(VERSION)_$(ARCH).deb"
