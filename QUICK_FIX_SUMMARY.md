# Quick Fix Summary

## What's the Problem?
Your GitHub Actions workflow fails with "Bad credentials" error. **The token has NOT expired.**

## Root Cause
The reusable workflow at `strongo/go-ci-action` has a bug where it uses `secrets.gh_token` (lowercase) instead of `secrets.GH_TOKEN` (uppercase).

## Solution

### In the `strongo/go-ci-action` Repository

You need to fix **3 lines** in `.github/workflows/workflow.yml`:

1. **Line 152** - Change:
   ```yaml
   github-token: ${{ secrets.gh_token }}
   ```
   To:
   ```yaml
   github-token: ${{ secrets.GH_TOKEN }}
   ```

2. **Line 172** - Change:
   ```yaml
   #     github_token: ${{ secrets.gh_token }}
   ```
   To:
   ```yaml
   #     github_token: ${{ secrets.GH_TOKEN }}
   ```

3. **Line 180** - Change:
   ```yaml
   github_token: ${{ secrets.gh_token }}
   ```
   To:
   ```yaml
   github_token: ${{ secrets.GH_TOKEN }}
   ```

### In This Repository (sneat-go-backend)

**NO CHANGES NEEDED!** Your configuration is already correct. Once the upstream is fixed, everything will work.

## How to Apply

1. Go to https://github.com/strongo/go-ci-action
2. Edit `.github/workflows/workflow.yml`
3. Make the 3 changes above (change `gh_token` to `GH_TOKEN` in all 3 places)
4. Commit with message: "Fix: Use GH_TOKEN instead of gh_token for secret references"

That's it! Your workflow will work after the upstream fix is merged.

## Need It Now?

If you can't wait for the upstream fix:
1. Fork `strongo/go-ci-action` to `sneat-co` organization
2. Make the fix in your fork
3. Temporarily change `.github/workflows/ci.yml` line 30 to:
   ```yaml
   uses: sneat-co/go-ci-action/.github/workflows/workflow.yml@main
   ```
