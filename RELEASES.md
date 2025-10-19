# Release Notes

## v0.1.1 (2025-10-19)

### Added
- `--version` and `-v` flags to display version information
- Version output now includes:
  - Version number
  - Git commit hash
  - Build timestamp

### Fixed
- Media upload test multipart form handling

### Installation
```bash
brew tap 0dayfall/tap
brew upgrade ctw
```

---

## v0.1.0 (2025-10-19)

Initial release of ctw - Command-line Twitter Client for automation.

### Features
- **Complete Twitter v2 API Support**
  - Tweets (create, delete, lookup)
  - Users (lookup, follow, block)
  - Media upload (images, videos with chunked upload)
  - Real-time streaming with filtered stream
  - Engagement (likes, retweets, bookmarks)
  - Timelines (user, mentions, home)
  - Direct messages
  - Tweet search and counts

- **Automation-First Design**
  - JSON output for easy parsing
  - Standard Unix exit codes
  - Environment variable configuration
  - Works with pipes and standard tools (jq, grep, awk)
  - Single binary, no dependencies

- **Distribution**
  - Homebrew tap (0dayfall/tap)
  - .deb packages (Ubuntu/Debian)
  - .rpm packages (RHEL/Fedora)
  - Pre-built binaries for Linux, macOS, Windows

- **Documentation**
  - Comprehensive automation guide (AUTOMATION.md)
  - Streaming guide (STREAMING.md)
  - Keyword screening reference
  - Installation instructions
  - Example scripts

### Installation
```bash
# macOS
brew tap 0dayfall/tap
brew install ctw

# Ubuntu/Debian
wget https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_linux_amd64.deb
sudo dpkg -i ctw_0.1.0_linux_amd64.deb

# From source
go install github.com/0dayfall/ctw/cmd/ctw@latest
```

### Quick Start
```bash
export BEARER_TOKEN="your_twitter_bearer_token"
ctw search recent --query "golang" | jq -r '.data[].text'
```

### Links
- GitHub: https://github.com/0dayfall/ctw
- Releases: https://github.com/0dayfall/ctw/releases
- Documentation: https://github.com/0dayfall/ctw#readme
