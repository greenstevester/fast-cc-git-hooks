#!/bin/bash
set -e

# GitHub Action entrypoint for conventional commit validation

echo "<<Made for Boo>>"
echo "🔍 Checking conventional commits..."

CONFIG_FILE="${1:-.fast-cc-git-hooks/config.yaml}"
BASE_BRANCH="${2:-$GITHUB_BASE_REF}"
FAIL_ON_ERROR="${3:-true}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Initialize variables
INVALID_COMMITS=""
HAS_ERRORS=false

# Get the commit range to check
if [ "$GITHUB_EVENT_NAME" = "pull_request" ]; then
    # For PRs, check commits between base and head
    echo "Checking commits in pull request #${GITHUB_PULL_REQUEST}"
    COMMIT_RANGE="${GITHUB_BASE_REF}..${GITHUB_HEAD_REF}"
    
    # Fetch the base branch to ensure we have the commits
    git fetch origin "${BASE_BRANCH}" --depth=50 || true
    
    # Get all commits in the PR
    COMMITS=$(git rev-list "origin/${BASE_BRANCH}..HEAD" 2>/dev/null || git rev-list HEAD~10..HEAD)
else
    # For pushes, check the last 10 commits or since last tag
    echo "Checking recent commits..."
    LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "")
    if [ -n "$LAST_TAG" ]; then
        COMMITS=$(git rev-list "${LAST_TAG}..HEAD")
    else
        COMMITS=$(git rev-list HEAD~10..HEAD 2>/dev/null || git rev-list HEAD)
    fi
fi

# Check if config file exists and use it
CONFIG_ARGS=""
if [ -f "$CONFIG_FILE" ]; then
    echo "Using config file: $CONFIG_FILE"
    CONFIG_ARGS="--config $CONFIG_FILE"
fi

# Validate each commit message
echo ""
echo "Validating commit messages..."
echo "================================"

for COMMIT in $COMMITS; do
    # Get the commit message
    MESSAGE=$(git log --format=%B -n 1 "$COMMIT")
    SUBJECT=$(git log --format=%s -n 1 "$COMMIT")
    HASH=$(git log --format=%h -n 1 "$COMMIT")
    
    # Validate the commit message
    if fcgh validate "$MESSAGE" $CONFIG_ARGS 2>/dev/null; then
        echo -e "${GREEN}✅${NC} $HASH: $SUBJECT"
    else
        echo -e "${RED}❌${NC} $HASH: $SUBJECT"
        ERROR_OUTPUT=$(fcgh validate "$MESSAGE" $CONFIG_ARGS 2>&1 || true)
        echo -e "${YELLOW}   Issue: ${NC}$ERROR_OUTPUT"
        INVALID_COMMITS="${INVALID_COMMITS}${HASH}: ${SUBJECT}\n"
        HAS_ERRORS=true
    fi
done

echo "================================"
echo ""

# Set outputs for GitHub Actions
if [ "$GITHUB_ACTIONS" = "true" ]; then
    if [ "$HAS_ERRORS" = "true" ]; then
        echo "valid=false" >> $GITHUB_OUTPUT
        echo "invalid-commits<<EOF" >> $GITHUB_OUTPUT
        echo -e "$INVALID_COMMITS" >> $GITHUB_OUTPUT
        echo "EOF" >> $GITHUB_OUTPUT
    else
        echo "valid=true" >> $GITHUB_OUTPUT
        echo "invalid-commits=" >> $GITHUB_OUTPUT
    fi
fi

# Summary
if [ "$HAS_ERRORS" = "true" ]; then
    echo -e "${RED}❌ Some commits do not follow conventional commit format${NC}"
    echo ""
    echo "Invalid commits:"
    echo -e "$INVALID_COMMITS"
    echo ""
    echo "Please update your commit messages to follow the conventional format:"
    echo "  <type>(<scope>): <subject>"
    echo ""
    echo "Example: feat(api): add user authentication endpoint"
    echo ""
    echo "Valid types: feat, fix, docs, style, refactor, perf, test, build, ci, chore"
    
    if [ "$FAIL_ON_ERROR" = "true" ]; then
        exit 1
    fi
else
    echo -e "${GREEN}✅ All commits follow conventional commit format!${NC}"
fi