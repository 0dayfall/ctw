# Quick Start Guide for GoReleaser

## Your Setup is Complete! ✅

GoReleaser is installed and your `.goreleaser.yml` is valid.

## Next Steps

### 1. Create Your Homebrew Tap Repository

```bash
gh repo create homebrew-tap --public --description "Homebrew tap for ctw"
```

Or manually at: https://github.com/new (name it `homebrew-tap`)

### 2. Set Your GitHub Token

Create a personal access token with `repo` scope:

```bash
# Get your current token (if using gh CLI)
gh auth token

# Or create a new one at:
# https://github.com/settings/tokens/new

# Set it as an environment variable
export HOMEBREW_TAP_GITHUB_TOKEN="ghp_your_token_here"

# Add to your ~/.zshrc for persistence
echo 'export HOMEBREW_TAP_GITHUB_TOKEN="ghp_your_token_here"' >> ~/.zshrc
```

### 3. Test Locally (Without Publishing)

```bash
# Test the build without creating a release
goreleaser release --snapshot --clean --skip=publish

# This creates files in dist/ directory
ls -la dist/
```

### 4. Create Your First Release

```bash
# Make sure everything is committed
git add .
git commit -m "Prepare for v0.1.0 release"
git push

# Tag the release
git tag -a v0.1.0 -m "Initial release v0.1.0"
git push origin v0.1.0

# Run GoReleaser (this will publish!)
goreleaser release --clean
```

### What GoReleaser Will Do:

1. ✅ Build binaries for:
   - Linux (amd64, arm64)
   - macOS (amd64 Intel, arm64 Apple Silicon)
   - Windows (amd64, arm64)

2. ✅ Create packages:
   - `.deb` for Ubuntu/Debian
   - `.rpm` for Fedora/RedHat
   - `.tar.gz` archives with README, LICENSE, INSTALL.md

3. ✅ Generate checksums for all artifacts

4. ✅ Create GitHub release with all files

5. ✅ Update your Homebrew tap automatically!

### After Release, Users Can Install:

**macOS (Homebrew):**
```bash
brew tap 0dayfall/tap
brew install ctw
```

**Ubuntu/Debian:**
```bash
wget https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_Linux_x86_64.deb
sudo dpkg -i ctw_0.1.0_Linux_x86_64.deb
```

**Linux (RPM):**
```bash
wget https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_Linux_x86_64.rpm
sudo rpm -i ctw_0.1.0_Linux_x86_64.rpm
```

**Windows:**
Download from GitHub releases page

## Troubleshooting

### "repository not found" error
- Ensure `homebrew-tap` repository exists and is public
- Verify token has `repo` scope
- Check owner name in `.goreleaser.yml` (currently set to `0dayfall`)

### "dirty working directory" error
```bash
git status  # Check for uncommitted changes
git add . && git commit -m "Clean up"
```

### Test without publishing
```bash
goreleaser release --snapshot --clean --skip=publish
```

## CI/CD Automation (Optional)

Want GitHub Actions to automatically release when you push a tag?

Create `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.18'
      
      - uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

Then add `HOMEBREW_TAP_GITHUB_TOKEN` as a repository secret in GitHub.

## Ready to Test?

Run this to test locally without publishing:

```bash
goreleaser release --snapshot --clean --skip=publish && ls -la dist/
```

This will show you exactly what will be created when you do a real release!
