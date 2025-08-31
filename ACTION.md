# GitHub Action: Conventional Commits Check

Use `fast-cc-git-hooks` as a GitHub Action to automatically validate conventional commit messages in your pull requests.

## Quick Setup

Add this workflow to your repository at `.github/workflows/pr-validation.yml`:

```yaml
name: PR Commit Validation

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  validate-commits:
    name: Validate Conventional Commits
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
      with:
        fetch-depth: 0  # Fetch all history for commit validation
    
    - name: Validate Commit Messages
      uses: greenstevester/fast-cc-git-hooks@main
      with:
        fail-on-error: 'true'
```

## Inputs

| Input | Description | Required | Default |
|-------|-------------|----------|---------|
| `config-file` | Path to the configuration file | No | `.fast-cc-git-hooks/config.yaml` |
| `base-branch` | Base branch to compare against | No | `${{ github.base_ref }}` |
| `fail-on-error` | Whether to fail the action if commits are invalid | No | `true` |

## Outputs

| Output | Description |
|--------|-------------|
| `valid` | Whether all commits are valid (true/false) |
| `invalid-commits` | List of invalid commit messages found |

## Example with Custom Configuration

```yaml
- name: Validate Commit Messages
  uses: greenstevester/fast-cc-git-hooks@main
  with:
    config-file: '.github/fcgh-config.yaml'
    fail-on-error: 'true'
```

## Example with PR Comments

```yaml
name: PR Commit Validation

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  validate-commits:
    name: Validate Conventional Commits
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout code
      uses: actions/checkout@v4
      with:
        fetch-depth: 0
    
    - name: Validate Commit Messages
      id: validate
      uses: greenstevester/fast-cc-git-hooks@main
      with:
        fail-on-error: 'false'  # Don't fail, just report
    
    - name: Comment on PR if validation fails
      if: steps.validate.outputs.valid == 'false'
      uses: actions/github-script@v7
      with:
        script: |
          github.rest.issues.createComment({
            issue_number: context.issue.number,
            owner: context.repo.owner,
            repo: context.repo.repo,
            body: \`## ❌ Commit Message Validation Failed
            
            Some commits in this PR do not follow the [Conventional Commits](https://www.conventionalcommits.org/) format.
            
            **Invalid commits:**
            \${{ steps.validate.outputs.invalid-commits }}
            
            **Expected format:** \\\`<type>(<scope>): <subject>\\\`
            
            **Valid types:** feat, fix, docs, style, refactor, perf, test, build, ci, chore
            
            **Example:** \\\`feat(api): add user authentication endpoint\\\`
            
            Please update your commit messages using \\\`git rebase -i\\\` or squash commits with proper messages.\`
          })
```

## What This Action Does

1. **Fetches commit history** for the pull request
2. **Validates each commit message** against conventional commit format
3. **Reports results** with clear success/failure indicators
4. **Provides detailed feedback** on what's wrong with invalid commits
5. **Optionally fails the workflow** to block merging of invalid PRs

## Conventional Commit Format

The action validates against this format:
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Valid Types
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests
- `build`: Changes that affect the build system
- `ci`: Changes to CI configuration
- `chore`: Other changes that don't modify src or test files

### Examples of Valid Commits
```
feat(auth): add JWT token validation
fix: resolve memory leak in image processing
docs(api): update authentication examples
style: fix code formatting
refactor(utils): simplify date formatting
perf: optimize database queries
test: add unit tests for user service
build: update dependencies
ci: add new workflow for releases
chore: update .gitignore
```

## Configuration

Create a `.fast-cc-git-hooks/config.yaml` file in your repository to customize validation rules:

```yaml
# Allowed commit types
types:
  - feat
  - fix
  - docs
  - style
  - refactor
  - perf
  - test
  - build
  - ci
  - chore

# Allowed scopes (optional)
scopes:
  - api
  - web
  - cli
  - db
  - auth

# Maximum subject length
max_subject_length: 72

# Require JIRA ticket references (enterprise)
require_jira_ticket: true
jira_project_keys:
  - PROJ
  - DEV
  - TECH
```

## Benefits

- ✅ **Automatic validation** on every PR
- ✅ **Clear feedback** to developers
- ✅ **Consistent commit history** across your team
- ✅ **Integration with semantic versioning** tools
- ✅ **Customizable rules** for your project needs
- ✅ **Zero configuration** required for basic usage

## Troubleshooting

### Action fails with "No commits found"
- Ensure `fetch-depth: 0` is set in the checkout step
- Check that the base branch exists and is accessible

### Custom config not being used
- Verify the config file path is correct
- Ensure the YAML syntax is valid
- Check that the config file is committed to the repository

### Commits pass locally but fail in Action
- The Action uses stricter validation by default
- Check that your local fcgh version matches the Action version
- Verify enterprise features aren't enabled unexpectedly