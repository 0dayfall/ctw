# GoReleaser Setup Guide

This guide will help you set up automated releases with GoReleaser for the `ctw` project.

## Prerequisites

- Go 1.18 or later
- Git
- GitHub account
- GitHub CLI (`gh`) installed

## Step 1: Install GoReleaser

### macOS
```bash
brew install goreleaser/tap/goreleaser
``# Add to your ~/.bashrc or ~/.zshrc
echo 'export HOMEBREW_TAP_GITHUB_TOKEN="ghp_your_token_here"' >> ~/.zshrc
source ~/.zshrc`

### Linux
```bash
# Using Go
go install github.com/goreleaser/goreleaser@latest

# Or download binary
curl -sfL https://goreleaser.com/static/run | bash
```

Verify installation:
```bash
goreleaser --version
```

## Step 2: Create Homebrew Tap Repository

Create a new public GitHub repository for your Homebrew tap:

```bash
gh repo create homebrew-tap --public --description "Homebrew tap for ctw"
```

Or manually:
1. Go to https://github.com/new
2. Name it `homebrew-tap`
3. Make it public
4. Create the repository

## Step 3: Create GitHub Personal Access Token

You need a token with `repo` scope for GoReleaser to push to your tap.

### Using GitHub CLI:
```bash
gh auth token
```

### Or manually:
1. Go to https://github.com/settings/tokens/new
2. Name it "GoReleaser Homebrew Tap"
3. Select scope: `repo` (all sub-scopes)
4. Click "Generate token"
5. Copy the token

### Set the token as environment variable:
```bash
# Add to your ~/.zshrc or ~/.bashrc
export HOMEBREW_TAP_GITHUB_TOKEN="ghp_your_token_here"

# Or set it temporarily
export HOMEBREW_TAP_GITHUB_TOKEN="ghp_your_token_here"
```

## Step 4: Update GoReleaser Config

The `.goreleaser.yml` file is already configured. Update the repository owner if needed:

```yaml
brews:
  - name: ctw
    repository:
      owner: 0dayfall  # Update this to your GitHub username
      name: homebrew-tap
```

## Step 5: Test GoReleaser Locally

Test the configuration without creating a release:

```bash
goreleaser check
goreleaser release --snapshot --clean
```

This creates a test build in the `dist/` directory without pushing anything.

## Step 6: Create Your First Release

### Tag a version:
```bash
git add .
git commit -m "Prepare for v0.1.0 release"
git push

git tag -a v0.1.0 -m "Initial release v0.1.0"
git push origin v0.1.0
```

### Run GoReleaser:
```bash
goreleaser release --clean
```

This will:
1. Build binaries for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)
2. Create `.deb` and `.rpm` packages
3. Generate checksums
4. Create a GitHub release with all artifacts
5. Update your Homebrew tap automatically

## Step 7: Verify the Release

1. Check your GitHub releases:
   ```bash
   gh release list
   ```

2. Check your homebrew-tap repository:
   ```bash
   gh repo view 0dayfall/homebrew-tap
   ```

3. Test installation:
   ```bash
   brew tap 0dayfall/tap
   brew install ctw
   ctw --version
   ```

## Future Releases

For subsequent releases, just tag and run:

```bash
# Update version
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# Release
goreleaser release --clean
```

## CI/CD Integration (Optional)

This repo includes GitHub Actions workflows:

- `CI` runs `go test ./...` on pushes and pull requests.
- `Release` runs tests on tags and then runs GoReleaser. It also runs the optional smoke test if a `BEARER_TOKEN` secret is configured.

Add `HOMEBREW_TAP_GITHUB_TOKEN` to your repository secrets:
1. Go to repository Settings → Secrets and variables → Actions
2. Click "New repository secret"
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Value: Your personal access token

Optional: add `BEARER_TOKEN` to enable the real API smoke test in release workflows.

Now just push a tag and GitHub Actions will handle the release automatically!

## Troubleshooting

### "repository not found" error
- Ensure `homebrew-tap` repository exists and is public
- Verify your token has `repo` scope
- Check the owner name in `.goreleaser.yml`

### "dirty working directory" error
```bash
git status  # Check for uncommitted changes
git add . && git commit -m "Clean up"
```

### Test without releasing
```bash
goreleaser release --snapshot --clean --skip-publish
```

## Resources

- GoReleaser docs: https://goreleaser.com/
- Homebrew tap docs: https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap
