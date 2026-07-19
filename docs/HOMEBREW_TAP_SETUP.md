# Setting Up Homebrew Tap for ctw

This guide shows how to create and configure a Homebrew tap for automated formula updates via GoReleaser.

## What is a Homebrew Tap?

A tap is a third-party GitHub repository containing Homebrew formulas. Users can install your software with:

```bash
brew tap 0dayfall/tap
brew install ctw
```

## Quick Setup

### 1. Create the Tap Repository

Create a new GitHub repository named `homebrew-tap`:

```bash
# On GitHub, create: 0dayfall/homebrew-tap
# Then clone it locally:
git clone https://github.com/0dayfall/homebrew-tap.git
cd homebrew-tap

# Create Formula directory
mkdir -p Formula

# Create README
cat > README.md << 'EOF'
# Homebrew Tap for 0dayfall

Homebrew formulas for 0dayfall projects.

## Install

```bash
brew tap 0dayfall/tap
brew install ctw
```

## Available Formulas

- **ctw** - Command-line toolkit for Twitter v2 API automation
EOF

# Commit and push
git add .
git commit -m "Initial tap setup"
git push origin main
```

### 2. Create GitHub Token

GoReleaser needs a token to push formula updates to your tap:

1. Go to **GitHub Settings** → **Developer settings** → **Personal access tokens** → **Tokens (classic)**
2. Click **"Generate new token (classic)"**
3. Set name: `HOMEBREW_TAP_TOKEN`
4. Select scopes:
   - ✅ `public_repo` (or `repo` if tap is private)
5. Click **"Generate token"**
6. **Copy the token** (you won't see it again!)

### 3. Add Token to Environment

```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc)
export HOMEBREW_TAP_GITHUB_TOKEN="ghp_your_token_here"

# Or add to repository secrets for GitHub Actions
# Settings → Secrets and variables → Actions → New repository secret
# Name: HOMEBREW_TAP_GITHUB_TOKEN
# Value: ghp_your_token_here
```

### 4. Test the Setup

The GoReleaser configuration (`.goreleaser.yml`) is already set up. Test it:

```bash
# Dry run to validate configuration
goreleaser release --snapshot --clean

# Check that builds work
ls -lh dist/

# When ready, create a release
git tag v0.1.0
git push origin v0.1.0
goreleaser release --clean
```

GoReleaser will automatically:
1. Build binaries for all platforms
2. Create GitHub release with artifacts
3. Update the Homebrew formula in `homebrew-tap`
4. Commit and push the formula update

## Using Your Tap

### Install from Tap

Once the first release is published:

```bash
# Add your tap
brew tap 0dayfall/tap

# Install ctw
brew install ctw

# Verify installation
ctw --version
```

### Update Formula

Users can update to the latest version:

```bash
brew update
brew upgrade ctw
```

### Uninstall

```bash
brew uninstall ctw
brew untap 0dayfall/tap
```

## GoReleaser Configuration

Your `.goreleaser.yml` already has the tap configured:

```yaml
brews:
  - repository:
      owner: 0dayfall              # Your GitHub username
      name: homebrew-tap           # Tap repository name
      branch: main                 # Branch to push to
      token: "{{ .Env.HOMEBREW_TAP_GITHUB_TOKEN }}"
    directory: Formula             # Where formulas go
    homepage: https://github.com/0dayfall/ctw
    description: Command-line toolkit for Twitter v2 API
    license: MIT
    install: |
      bin.install "ctw"
    test: |
      system "#{bin}/ctw", "--version"
```

## Release Workflow

### Manual Release

```bash
# 1. Ensure clean working directory
git status

# 2. Create and push tag
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# 3. Run GoReleaser (uses HOMEBREW_TAP_GITHUB_TOKEN)
goreleaser release --clean
```

### Automated Release (GitHub Actions)

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

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v5
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

Then just push a tag to trigger the release:

```bash
git tag v0.1.0
git push origin v0.1.0
# GitHub Actions will handle the rest
```

## Verifying the Tap

After your first release, check that the formula was created:

```bash
# Visit: https://github.com/0dayfall/homebrew-tap/tree/main/Formula
# You should see: ctw.rb

# Or clone and check locally
git clone https://github.com/0dayfall/homebrew-tap.git
cat homebrew-tap/Formula/ctw.rb
```

The formula should look like:

```ruby
class Ctw < Formula
  desc "Command-line toolkit for Twitter v2 API"
  homepage "https://github.com/0dayfall/ctw"
  version "0.1.0"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_Darwin_arm64.tar.gz"
      sha256 "..."
    end
    if Hardware::CPU.intel?
      url "https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_Darwin_x86_64.tar.gz"
      sha256 "..."
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_Linux_arm64.tar.gz"
      sha256 "..."
    end
    if Hardware::CPU.intel?
      url "https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_Linux_x86_64.tar.gz"
      sha256 "..."
    end
  end

  def install
    bin.install "ctw"
  end

  test do
    system "#{bin}/ctw", "--version"
  end
end
```

## Troubleshooting

### "Error: No available formula with the name"

Make sure:
- You've pushed at least one release tag
- GoReleaser successfully updated the tap
- Formula exists in `homebrew-tap/Formula/ctw.rb`

### "Permission denied" during release

Check that `HOMEBREW_TAP_GITHUB_TOKEN` has `public_repo` scope.

### Formula not updating

1. Check GoReleaser logs for errors
2. Verify token is set: `echo $HOMEBREW_TAP_GITHUB_TOKEN`
3. Check tap repository for commits
4. Try `brew untap 0dayfall/tap && brew tap 0dayfall/tap`

### Test locally before publishing

```bash
# Build snapshot without publishing
goreleaser release --snapshot --clean

# Install from local build
brew install ./dist/homebrew/Formula/ctw.rb
```

## Best Practices

1. **Version Tags**: Use semantic versioning (`v1.0.0`, `v1.0.1`)
2. **Changelog**: Keep `CHANGELOG.md` updated
3. **Testing**: Test formula locally before releasing
4. **Documentation**: Update README with installation instructions
5. **Automation**: Use GitHub Actions for consistent releases

## Update README

Add installation instructions to your main README:

```markdown
## Installation

### macOS (Homebrew)

```bash
brew tap 0dayfall/tap
brew install ctw
```

### Ubuntu/Debian

Download the `.deb` package from [releases](https://github.com/0dayfall/ctw/releases):

```bash
wget https://github.com/0dayfall/ctw/releases/download/v0.1.0/ctw_0.1.0_amd64.deb
sudo dpkg -i ctw_0.1.0_amd64.deb
```

### From Source

```bash
git clone https://github.com/0dayfall/ctw.git
cd ctw
go build -o ctw ./cmd/ctw
sudo mv ctw /usr/local/bin/
```
```

## Resources

- [Homebrew Documentation](https://docs.brew.sh/)
- [Creating Taps](https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap)
- [GoReleaser Homebrew](https://goreleaser.com/customization/homebrew/)
- [Formula Cookbook](https://docs.brew.sh/Formula-Cookbook)

## Summary

**Steps to Create Tap:**
1. ✅ Create `homebrew-tap` repository on GitHub
2. ✅ Generate GitHub token with `public_repo` scope
3. ✅ Export token as `HOMEBREW_TAP_GITHUB_TOKEN`
4. ✅ Create release tag and run GoReleaser
5. ✅ Formula is automatically created/updated in tap

**Users Install With:**
```bash
brew tap 0dayfall/tap
brew install ctw
```

That's it! Your Homebrew tap is ready for automated releases. 🍺
