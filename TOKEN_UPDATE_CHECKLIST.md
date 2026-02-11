# Token Update Checklist

## Quick Checklist ✓

Use this checklist to fix the "Bad credentials" error:

### Part 1: Generate New Token

- [ ] Go to https://github.com/settings/tokens
- [ ] Click "Generate new token (classic)"
- [ ] Enter token name: "Sneat Go Backend CI" (or similar)
- [ ] Set expiration: 
  - [ ] 90 days (recommended)
  - [ ] 180 days
  - [ ] No expiration (less secure)
- [ ] Select scopes:
  - [ ] ✅ `repo` (Full control of private repositories)
- [ ] Click "Generate token"
- [ ] Copy the token immediately (you won't see it again!)
- [ ] Store it temporarily in a secure location

### Part 2: Update Repository Secret

- [ ] Go to https://github.com/sneat-co/sneat-go-backend/settings/secrets/actions
- [ ] Find secret named: `GH_PAT_READWRITE_REPOS`
- [ ] Click "Update" button
- [ ] Paste the new token value
- [ ] Click "Update secret" or "Add secret"
- [ ] Confirm the secret shows as updated

### Part 3: Verify the Fix

- [ ] Go to https://github.com/sneat-co/sneat-go-backend/actions
- [ ] Find the most recent failed workflow
- [ ] Click "Re-run all jobs" or "Re-run failed jobs"
- [ ] Wait for workflow to complete
- [ ] Check that:
  - [ ] No "Bad credentials" error appears
  - [ ] Tests pass successfully
  - [ ] Version tagging completes (if on main branch)

### Part 4: Document (Optional but Recommended)

- [ ] Note token creation date: ________________
- [ ] Note token expiration date: ________________
- [ ] Set calendar reminder for renewal: ________________
- [ ] Update team documentation if needed

---

## Troubleshooting

If it still fails after updating the token:

### Check 1: Token Permissions
The token MUST have `repo` scope. Verify:
- [ ] Go to https://github.com/settings/tokens
- [ ] Find your token
- [ ] Confirm it has `repo` scope checked

### Check 2: Secret Name
Verify the secret name is exactly:
- [ ] `GH_PAT_READWRITE_REPOS` (case-sensitive, no spaces)

### Check 3: Token Owner
The token must be from an account that has write access to this repository:
- [ ] Token is from a user who is a collaborator
- [ ] Or token is from the repository owner's account

### Check 4: Organization Settings
If this is an organization repository:
- [ ] Check organization's OAuth App policy
- [ ] Check organization's Personal Access Token policy
- [ ] Ensure tokens are allowed for CI/CD

---

## Alternative: Use GITHUB_TOKEN (Test First)

If you want to avoid token management:

- [ ] Edit `.github/workflows/ci.yml`
- [ ] Change line 34 to: `GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}`
- [ ] Commit and push
- [ ] Test if it works with your GOPRIVATE dependencies
- [ ] If `go get` fails, revert and stick with PAT approach

---

## Status

Current Status:
- [ ] Not started
- [ ] Token generated
- [ ] Secret updated
- [ ] ✅ Verified working

Last Updated: ________________
Updated By: ________________

---

## Need Help?

See these documents for more information:
- `VERIFICATION_SUMMARY.md` - Quick explanation of the issue
- `WORKFLOW_ANALYSIS.md` - Detailed technical analysis
- `ALTERNATIVE_SOLUTION.md` - Alternative approaches
