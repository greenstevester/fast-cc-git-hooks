#!/bin/bash
# fcgh Cross-Platform Installation Script
# Detects platform and installs fcgh, ccg, and ccdo

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 fcgh Cross-Platform Installer${NC}"
echo ""

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $OS in
    darwin)
        OS_NAME="macOS"
        if [[ "$ARCH" == "arm64" ]]; then
            BINARY_ARCH="darwin_arm64"
            echo -e "${BLUE}🍎 Detected Apple Silicon (M1/M2/M3) Mac${NC}"
        elif [[ "$ARCH" == "x86_64" ]]; then
            BINARY_ARCH="darwin_amd64"
            echo -e "${BLUE}🍎 Detected Intel Mac${NC}"
        else
            echo -e "${RED}❌ Unsupported macOS architecture: $ARCH${NC}"
            exit 1
        fi
        ;;
    linux)
        OS_NAME="Linux"
        if [[ "$ARCH" == "x86_64" ]]; then
            BINARY_ARCH="linux_amd64"
            echo -e "${BLUE}🐧 Detected Linux x86_64${NC}"
        elif [[ "$ARCH" == "aarch64" ]] || [[ "$ARCH" == "arm64" ]]; then
            BINARY_ARCH="linux_arm64"
            echo -e "${BLUE}🐧 Detected Linux ARM64${NC}"
        else
            echo -e "${RED}❌ Unsupported Linux architecture: $ARCH${NC}"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}❌ Unsupported operating system: $OS${NC}"
        echo -e "${YELLOW}💡 Try manual installation from: https://github.com/greenstevester/fast-cc-git-hooks/releases${NC}"
        exit 1
        ;;
esac

# Get latest release info from GitHub
echo -e "${BLUE}📦 Fetching latest release information...${NC}"
LATEST_RELEASE=$(curl -s https://api.github.com/repos/greenstevester/fast-cc-git-hooks/releases/latest)
VERSION=$(echo "$LATEST_RELEASE" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [[ -z "$VERSION" ]]; then
    echo -e "${RED}❌ Failed to get latest version${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Latest version: $VERSION${NC}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Download the release
DOWNLOAD_URL="https://github.com/greenstevester/fast-cc-git-hooks/releases/download/$VERSION/fcgh_${VERSION}_${BINARY_ARCH}.tar.gz"
echo -e "${BLUE}⬇️  Downloading fcgh $VERSION for $BINARY_ARCH...${NC}"

# Try downloading with error handling
if ! curl -L -f -o fcgh.tar.gz "$DOWNLOAD_URL"; then
    echo -e "${RED}❌ Failed to download from releases. Trying alternative URL format...${NC}"
    
    # Try alternative URL format (without version prefix)
    ALT_URL="https://github.com/greenstevester/fast-cc-git-hooks/releases/download/$VERSION/fast-cc-git-hooks_${VERSION}_${BINARY_ARCH}.tar.gz"
    echo -e "${BLUE}⬇️  Trying alternative URL: fast-cc-git-hooks format...${NC}"
    
    if ! curl -L -f -o fcgh.tar.gz "$ALT_URL"; then
        echo -e "${RED}❌ Release not found. This might mean:${NC}"
        echo -e "${YELLOW}   1. No release has been published yet${NC}"
        echo -e "${YELLOW}   2. The release is still being built${NC}"
        echo -e "${YELLOW}   3. The URL format has changed${NC}"
        echo ""
        echo -e "${BLUE}💡 Alternative: Build from source${NC}"
        echo -e "${YELLOW}   git clone https://github.com/greenstevester/fast-cc-git-hooks.git${NC}"
        echo -e "${YELLOW}   cd fast-cc-git-hooks${NC}"
        echo -e "${YELLOW}   make build-all-tools${NC}"
        echo -e "${YELLOW}   sudo cp build/fcgh build/ccg build/ccdo /usr/local/bin/${NC}"
        echo ""
        exit 1
    fi
fi

# Validate downloaded file
if [[ ! -f "fcgh.tar.gz" ]] || [[ ! -s "fcgh.tar.gz" ]]; then
    echo -e "${RED}❌ Downloaded file is empty or missing${NC}"
    exit 1
fi

# Check if it's a valid tar.gz file
if ! file fcgh.tar.gz | grep -q "gzip compressed"; then
    echo -e "${RED}❌ Downloaded file is not a valid gzip archive${NC}"
    echo -e "${YELLOW}File contents:${NC}"
    head -n 5 fcgh.tar.gz
    exit 1
fi

# Extract
echo -e "${BLUE}📂 Extracting files...${NC}"
if ! tar -xzf fcgh.tar.gz; then
    echo -e "${RED}❌ Failed to extract archive${NC}"
    exit 1
fi

# Make binaries executable
chmod +x fcgh ccg ccdo 2>/dev/null || chmod +x fcgh ccg  # ccdo might not exist in older versions

# Determine installation directory
if [[ "$OS" == "darwin" ]]; then
    INSTALL_DIR="/usr/local/bin"
else
    INSTALL_DIR="/usr/local/bin"
fi

if [[ ! -d "$INSTALL_DIR" ]]; then
    echo -e "${YELLOW}⚠️  $INSTALL_DIR doesn't exist, creating it...${NC}"
    sudo mkdir -p "$INSTALL_DIR"
fi

# Install binaries
echo -e "${BLUE}📦 Installing binaries to $INSTALL_DIR...${NC}"
sudo cp fcgh "$INSTALL_DIR/"
sudo cp ccg "$INSTALL_DIR/"
[[ -f ccdo ]] && sudo cp ccdo "$INSTALL_DIR/"

# Remove quarantine attributes on macOS (prevents security warnings)
if [[ "$OS" == "darwin" ]]; then
    echo -e "${BLUE}🔓 Removing quarantine attributes...${NC}"
    sudo xattr -d com.apple.quarantine "$INSTALL_DIR/fcgh" 2>/dev/null || true
    sudo xattr -d com.apple.quarantine "$INSTALL_DIR/ccg" 2>/dev/null || true
    [[ -f "$INSTALL_DIR/ccdo" ]] && sudo xattr -d com.apple.quarantine "$INSTALL_DIR/ccdo" 2>/dev/null || true
fi

# Clean up
cd /
rm -rf "$TMP_DIR"

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

# Verify installation
echo -e "${BLUE}🔍 Verifying installation...${NC}"
if command -v fcgh >/dev/null 2>&1; then
    VERSION_OUTPUT=$(fcgh version 2>/dev/null || echo "fcgh installed")
    echo -e "${GREEN}✅ fcgh: $VERSION_OUTPUT${NC}"
else
    echo -e "${RED}❌ fcgh not found in PATH${NC}"
    exit 1
fi

if command -v ccg >/dev/null 2>&1; then
    echo -e "${GREEN}✅ ccg: installed${NC}"
else
    echo -e "${RED}❌ ccg not found in PATH${NC}"
fi

if command -v ccdo >/dev/null 2>&1; then
    echo -e "${GREEN}✅ ccdo: installed${NC}"
else
    echo -e "${YELLOW}⚠️  ccdo not found (may not be included in this version)${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Installation completed successfully!${NC}"
echo ""
echo -e "${BLUE}📚 Quick Start:${NC}"
echo -e "   ${YELLOW}fcgh setup-ent${NC}     # Set up git hooks with enterprise features"
echo -e "   ${YELLOW}ccg${NC}                # Preview commit message"
echo -e "   ${YELLOW}ccdo${NC}               # Generate and commit automatically"
echo ""
echo -e "${BLUE}💡 Need help? Run:${NC}"
echo -e "   ${YELLOW}fcgh --help${NC}"
echo -e "   ${YELLOW}ccg --help${NC}"
echo ""
echo -e "${GREEN}🚀 Happy committing!${NC}"