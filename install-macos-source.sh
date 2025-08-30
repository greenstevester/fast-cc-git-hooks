#!/bin/bash
# fcgh macOS Build from Source Installation Script
# Use this if the regular installation script fails due to missing releases

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🛠️  fcgh Build from Source Installation${NC}"
echo -e "${YELLOW}This script will build fcgh from source code${NC}"
echo ""

# Check if Go is installed
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo -e "${BLUE}💡 Install Go first:${NC}"
    echo -e "${YELLOW}   brew install go${NC}"
    echo -e "${YELLOW}   # or download from https://golang.org/dl/${NC}"
    exit 1
fi

GO_VERSION=$(go version | cut -d' ' -f3)
echo -e "${GREEN}✅ Found Go: $GO_VERSION${NC}"

# Check if git is installed
if ! command -v git >/dev/null 2>&1; then
    echo -e "${RED}❌ Git is not installed${NC}"
    echo -e "${BLUE}💡 Install Git first:${NC}"
    echo -e "${YELLOW}   brew install git${NC}"
    exit 1
fi

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Clone the repository
echo -e "${BLUE}📥 Cloning repository...${NC}"
git clone https://github.com/greenstevester/fast-cc-git-hooks.git
cd fast-cc-git-hooks

# Build the project
echo -e "${BLUE}🔨 Building fcgh...${NC}"
if ! make build; then
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi

# Check if binaries were created
if [[ ! -f "build/fcgh" ]]; then
    echo -e "${RED}❌ fcgh binary not found after build${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Build successful${NC}"

# Determine installation directory
INSTALL_DIR="/usr/local/bin"
if [[ ! -d "$INSTALL_DIR" ]]; then
    echo -e "${YELLOW}⚠️  $INSTALL_DIR doesn't exist, creating it...${NC}"
    sudo mkdir -p "$INSTALL_DIR"
fi

# Install binaries
echo -e "${BLUE}📦 Installing binaries to $INSTALL_DIR...${NC}"
sudo cp build/fcgh "$INSTALL_DIR/"
[[ -f "build/cc" ]] && sudo cp build/cc "$INSTALL_DIR/"
[[ -f "build/ccc" ]] && sudo cp build/ccc "$INSTALL_DIR/"

# Make sure they're executable
sudo chmod +x "$INSTALL_DIR/fcgh"
[[ -f "$INSTALL_DIR/cc" ]] && sudo chmod +x "$INSTALL_DIR/cc"
[[ -f "$INSTALL_DIR/ccc" ]] && sudo chmod +x "$INSTALL_DIR/ccc"

# Check if /usr/local/bin is in PATH
if [[ ":$PATH:" != *":/usr/local/bin:"* ]]; then
    echo -e "${YELLOW}⚠️  /usr/local/bin is not in your PATH${NC}"
    echo -e "${BLUE}📝 Adding /usr/local/bin to your shell profile...${NC}"
    
    # Detect shell and add to appropriate profile
    if [[ "$SHELL" == *"zsh"* ]] || [[ -n "$ZSH_VERSION" ]]; then
        PROFILE="$HOME/.zshrc"
        echo 'export PATH="/usr/local/bin:$PATH"' >> "$PROFILE"
        echo -e "${GREEN}✅ Added to $PROFILE${NC}"
    elif [[ "$SHELL" == *"bash"* ]] || [[ -n "$BASH_VERSION" ]]; then
        PROFILE="$HOME/.bash_profile"
        [[ ! -f "$PROFILE" ]] && PROFILE="$HOME/.bashrc"
        echo 'export PATH="/usr/local/bin:$PATH"' >> "$PROFILE"
        echo -e "${GREEN}✅ Added to $PROFILE${NC}"
    else
        echo -e "${YELLOW}⚠️  Unknown shell. Please manually add /usr/local/bin to your PATH${NC}"
    fi
    
    echo -e "${BLUE}💡 Run 'source $PROFILE' or restart your terminal to update PATH${NC}"
fi

# Clean up
cd /
rm -rf "$TMP_DIR"

# Verify installation
echo -e "${BLUE}🔍 Verifying installation...${NC}"
if command -v fcgh >/dev/null 2>&1; then
    VERSION_OUTPUT=$(fcgh version 2>/dev/null || echo "fcgh installed")
    echo -e "${GREEN}✅ fcgh: $VERSION_OUTPUT${NC}"
else
    echo -e "${RED}❌ fcgh not found in PATH${NC}"
    exit 1
fi

if command -v cc >/dev/null 2>&1; then
    echo -e "${GREEN}✅ cc: installed${NC}"
else
    echo -e "${YELLOW}⚠️  cc not found (may not have been built)${NC}"
fi

if command -v ccc >/dev/null 2>&1; then
    echo -e "${GREEN}✅ ccc: installed${NC}"
else
    echo -e "${YELLOW}⚠️  ccc not found (may not have been built)${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Installation completed successfully!${NC}"
echo ""
echo -e "${BLUE}📚 Quick Start:${NC}"
echo -e "   ${YELLOW}fcgh setup-ent${NC}     # Set up git hooks with enterprise features"
echo -e "   ${YELLOW}cc${NC}                # Preview commit message (if available)"
echo -e "   ${YELLOW}ccc${NC}               # Generate and commit automatically (if available)"
echo ""
echo -e "${BLUE}💡 Need help? Run:${NC}"
echo -e "   ${YELLOW}fcgh --help${NC}"
echo ""
echo -e "${GREEN}🚀 Happy committing!${NC}"