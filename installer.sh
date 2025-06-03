#!/bin/bash

# Convcommit Installer Script
# This script downloads and installs the convcommit tool globally

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Convcommit Installer${NC}"
echo "=============================="
echo "This script will install convcommit globally on your system."

# Determine OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture to Go architecture naming
case "${ARCH}" in
    x86_64) ARCH="amd64" ;;
    i686|i386) ARCH="386" ;;
    aarch64|arm64) ARCH="arm64" ;;
    armv7*) ARCH="arm" ;;
    *)
        echo -e "${RED}Unsupported architecture: ${ARCH}${NC}"
        exit 1
        ;;
esac

# Set installation directory based on OS
INSTALL_DIR=""
if [ "$OS" = "darwin" ] || [ "$OS" = "linux" ]; then
    INSTALL_DIR="/usr/local/bin"
elif [ "$OS" = "windows" ]; then
    echo -e "${RED}Windows installation is not supported by this script.${NC}"
    echo "Please download the Windows binary from the GitHub releases page."
    exit 1
else
    echo -e "${RED}Unsupported operating system: ${OS}${NC}"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}Go is not installed. Installing convcommit from pre-built binary.${NC}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    # Download the latest release
    echo "Downloading convcommit..."
    GITHUB_REPO="danfordChris/conventional_commit_CLI"
    DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/latest/download/convcommit_${OS}_${ARCH}.tar.gz"

    if ! curl -sL "$DOWNLOAD_URL" -o convcommit.tar.gz; then
        echo -e "${RED}Failed to download convcommit.${NC}"
        echo "Please check your internet connection and try again."
        exit 1
    fi

    # Extract the archive
    tar -xzf convcommit.tar.gz

    # Install the binary
    echo "Installing convcommit to ${INSTALL_DIR}..."
    sudo mv convcommit "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/convcommit"

    # Clean up
    cd - > /dev/null
    rm -rf "$TMP_DIR"
else
    echo -e "${GREEN}Go is installed. Building and installing convcommit from source.${NC}"

    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    # Clone the repository
    echo "Cloning the repository..."
    GITHUB_REPO="danfordChris/conventional_commit_CLI"
    git clone "https://github.com/${GITHUB_REPO}.git" .

    # Build the binary
    echo "Building convcommit..."
    cd cmd/convcommit
    go build -o convcommit

    # Install the binary
    echo "Installing convcommit to ${INSTALL_DIR}..."
    sudo mv convcommit "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/convcommit"

    # Clean up
    cd - > /dev/null
    rm -rf "$TMP_DIR"
fi

# Verify installation
if command -v convcommit &> /dev/null; then
    echo -e "${GREEN}Convcommit has been successfully installed!${NC}"
    echo "You can now use the 'convcommit' command from anywhere."
    echo "Run 'convcommit --help' for usage information."
else
    echo -e "${RED}Installation failed.${NC}"
    echo "Please check the error messages above and try again."
    exit 1
fi
