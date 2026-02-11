# Fix for "Bad Credentials" Error in GitHub Actions

## Problem
The GitHub Actions workflow is failing with "Bad credentials" error when using the reusable workflow from `strongo/go-ci-action`.

## Root Cause
The reusable workflow at `strongo/go-ci-action/.github/workflows/workflow.yml` has a bug:
- It defines `GH_TOKEN` (uppercase) as a required secret input
- But references it as `secrets.gh_token` (lowercase) in three places

This causes the actions to receive undefined/null tokens, resulting in authentication failures.

## Required Changes in strongo/go-ci-action Repository

The following changes need to be made in the `strongo/go-ci-action` repository's `.github/workflows/workflow.yml` file:

### 1. Line 152: Fix goveralls action
```yaml
# Before:
github-token: ${{ secrets.gh_token }}

# After:
github-token: ${{ secrets.GH_TOKEN }}
```

### 2. Line 172: Fix commented LSIF action
```yaml
# Before:
#     github_token: ${{ secrets.gh_token }}

# After:
#     github_token: ${{ secrets.GH_TOKEN }}
```

### 3. Line 180: Fix github-tag-action (main cause of error)
```yaml
# Before:
github_token: ${{ secrets.gh_token }}

# After:
github_token: ${{ secrets.GH_TOKEN }}
```

### 4. Line 163: Fix formatting (optional but recommended)
```yaml
# Before:
- if : ${{ inputs.min_test_coverage_percent != '' }}

# After:
- if: ${{ inputs.min_test_coverage_percent != '' }}
```

## How to Apply the Fix

### Option 1: Manual Fix
1. Fork or clone the `strongo/go-ci-action` repository
2. Create a new branch: `git checkout -b fix/correct-secret-reference`
3. Make the four changes listed above in `.github/workflows/workflow.yml`
4. Commit: `git commit -m "Fix: Use GH_TOKEN instead of gh_token for secret references"`
5. Push and create a Pull Request

### Option 2: Apply Patch File
A patch file is available at `upstream-fix.patch` in this repository.

To apply it:
```bash
cd /path/to/strongo/go-ci-action
git apply /path/to/sneat-go-backend/upstream-fix.patch
git commit -am "Fix: Use GH_TOKEN instead of gh_token for secret references"
```

## Changes Needed in This Repository (sneat-go-backend)

**NONE!** Once the upstream repository is fixed, this repository's workflow will work correctly as-is. The current configuration in `.github/workflows/ci.yml` is already correct:

```yaml
secrets:
  GH_TOKEN: ${{ secrets.GH_PAT_READWRITE_REPOS }}
```

## Temporary Workaround

If you need an immediate fix before the upstream repository is updated, you can:
1. Fork `strongo/go-ci-action` to your own account or the `sneat-co` organization
2. Apply the fix to your fork
3. Temporarily update `.github/workflows/ci.yml` to use your fork:
   ```yaml
   uses: your-org/go-ci-action/.github/workflows/workflow.yml@fix/correct-secret-reference
   ```

## Timeline

The bug was introduced in commit `29cf8d73c6906c26f7eef4295fa502de814f96e7` on August 12, 2024, when the secret was renamed from `gh_token` to `GH_TOKEN` but the usage references were not updated.

## Notes

- The token has **NOT** expired
- The error occurs because `secrets.gh_token` is undefined (doesn't exist)
- GitHub treats undefined secrets as empty/null values, causing "Bad credentials" errors
