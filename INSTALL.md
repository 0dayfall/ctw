# Installation Guide

## macOS (Homebrew)

### From Homebrew Tap

Once you've set up your Homebrew tap (see "Distribution" section below):

```bash
brew tap 0dayfall/tap
brew install ctw
```

### Update

```bash
brew upgrade ctw
```

### Uninstall

```bash
brew uninstall ctw
```

## Ubuntu/Debian Package Installation

### Quick Install (Pre-built .deb)

If you have a pre-built `.deb` package:

```bash
sudo dpkg -i ctw_<version>_amd64.deb
```

### Building from Source

#### Prerequisites

- Go 1.18 or later
- `make`
- `dpkg-deb` (for building .deb packages)

#### Build the .deb Package

1. Clone the repository:
   ```bash
   git clone https://github.com/0dayfall/ctw.git
   cd ctw
   ```

2. Build the package:
   ```bash
   make deb
   ```

3. Install the package:
   ```bash
   sudo dpkg -i build/ctw_<version>_amd64.deb
   ```

#### Multi-Architecture Builds

To build for both amd64 and arm64:

```bash
make deb-multi
```

This creates:
- `build/ctw_<version>_amd64.deb`
- `build/ctw_<version>_arm64.deb`

### Uninstall

```bash
sudo dpkg -r ctw
```

## Manual Installation (without .deb)

### Build Binary

```bash
make build
```

### Install to System

```bash
sudo make install
```

This installs to `/usr/local/bin` by default. To install elsewhere:

```bash
sudo make install INSTALL_PREFIX=/opt/ctw
```

### Configuration

After installation, configure your Twitter API credentials:

```bash
ctw init
```

This verifies access and writes a config file at:
- macOS/Linux: `~/.config/ctw/config.toml`
- Windows: `%APPDATA%\\ctw\\config.toml`

You can also configure via environment variables:

```bash
export BEARER_TOKEN="<your twitter bearer token>"
# Optional:
export USER_AGENT="my-client/1.0"
```

Consider adding these to your `~/.bashrc` or `~/.zshrc` for persistence.

See `config.example.toml` for a ready-to-copy template.

## Verify Installation

```bash
ctw --help
ctw --version
```

## First Steps

1. Test your credentials:
   ```bash
   ctw users lookup --usernames twitter
   ```

2. Try uploading media:
   ```bash
   ctw media upload --file image.jpg --category tweet_image
   ```

3. Create a tweet:
   ```bash
   ctw tweets create --text "Hello from ctw!"
   ```

## Distribution

### Using GoReleaser (Recommended)

[GoReleaser](https://goreleaser.com/) automates building binaries for multiple platforms, creating `.deb`/`.rpm` packages, and updating your Homebrew tap.

#### One-time Setup

1. Install GoReleaser:
   ```bash
   # macOS
   brew install goreleaser/tap/goreleaser
   
   # Linux
   go install github.com/goreleaser/goreleaser@latest
   ```

2. Create a Homebrew tap repository:
   ```bash
   # Create a new GitHub repository named 'homebrew-tap'
   gh repo create homebrew-tap --public
   ```

3. Set GitHub token for Homebrew tap updates:
   ```bash
   export HOMEBREW_TAP_GITHUB_TOKEN="<your github token>"
   ```

#### Creating a Release

1. Tag a new version:
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   git push origin v0.1.0
   ```

2. Run GoReleaser:
   ```bash
   goreleaser release --clean
   ```

This will:
- Build binaries for Linux, macOS, Windows (amd64 & arm64)
- Create `.deb` and `.rpm` packages
- Generate checksums
- Create a GitHub release with all artifacts
- Update your Homebrew tap automatically

#### Test Release Locally

```bash
goreleaser release --snapshot --clean
```

### Manual Distribution

#### GitHub Releases

Upload the `.deb` packages to GitHub Releases:

```bash
gh release create v0.1.0 \
  build/ctw_0.1.0_amd64.deb \
  build/ctw_0.1.0_arm64.deb \
  --title "Release v0.1.0" \
  --notes "Initial release"
```

Users can then download and install:

```bash
wget https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_amd64.deb
sudo dpkg -i ctw_0.1.0_amd64.deb
```

### PPA (Advanced)

For a full Ubuntu PPA setup, see: https://help.launchpad.net/Packaging/PPA

### Manual Homebrew Tap Setup

If not using GoReleaser's automatic tap updates:

1. Create a `homebrew-tap` repository on GitHub

2. Copy `Formula/ctw.rb` to your tap repository

3. Update the formula with release URLs and SHA256 checksums:
   ```bash
   # Calculate SHA256 for each release archive
   shasum -a 256 ctw_0.1.0_darwin_amd64.tar.gz
   ```

4. Users install with:
   ```bash
   brew tap 0dayfall/tap
   brew install ctw
   ```

### APT Repository (Advanced)

To host your own APT repository, use tools like:
- `aptly` (https://www.aptly.info/)
- `reprepro` (https://mirrorer.alioth.debian.org/)

## Development Installation

For development work:

```bash
go mod tidy
go run ./cmd/ctw --help
```

Run tests:

```bash
make test
# or
go test ./...
```
