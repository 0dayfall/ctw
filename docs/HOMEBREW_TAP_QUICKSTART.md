# Quick Start: Homebrew Tap

## TL;DR - Create Your Tap in 5 Minutes

### 1. Create Repository

On GitHub, create a new repository: **`homebrew-tap`**

```bash
# Repository name MUST be: homebrew-tap
# Owner: 0dayfall
# Public repository
```

### 2. Initialize Tap

```bash
git clone https://github.com/0dayfall/homebrew-tap.git
cd homebrew-tap
mkdir -p Formula
cat > README.md << 'EOF'
# Homebrew Tap for ctw

## Install
```bash
brew tap 0dayfall/tap
brew install ctw
```
EOF

git add .
git commit -m "Initial tap setup"
git push origin main
```

### 3. Create GitHub Token

1. GitHub → Settings → Developer settings → Personal access tokens → **Generate new token (classic)**
2. Name: `HOMEBREW_TAP_TOKEN`
3. Scope: ✅ `public_repo`
4. Copy token

### 4. Add Token to Environment

```bash
# Add to ~/.bashrc or ~/.zshrc
export HOMEBREW_TAP_GITHUB_TOKEN="ghp_your_token_here"

# Reload shell
source ~/.zshrc  # or ~/.bashrc
```

### 5. Create First Release

```bash
cd /path/to/ctw

# Ensure token is set
echo $HOMEBREW_TAP_GITHUB_TOKEN

# Tag and release
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
goreleaser release --clean
```

GoReleaser will automatically create the formula in your tap!

## Test Installation

```bash
# Add tap
brew tap 0dayfall/tap

# Install
brew install ctw

# Verify
ctw --version
```

## Update README.md

Add to your ctw README.md:

```markdown
## Installation

### macOS (Homebrew)

\```bash
brew tap 0dayfall/tap
brew install ctw
\```
```

## Done! 🍺

Users can now install with:
```bash
brew tap 0dayfall/tap && brew install ctw
```

---

For detailed setup and troubleshooting, see [HOMEBREW_TAP_SETUP.md](HOMEBREW_TAP_SETUP.md)
