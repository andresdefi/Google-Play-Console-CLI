#!/usr/bin/env bash
# gpc installer - https://github.com/andresdefi/Google-Play-Console-CLI
# Usage: curl -sSfL https://raw.githubusercontent.com/andresdefi/Google-Play-Console-CLI/main/install.sh | bash

set -euo pipefail

REPO="andresdefi/Google-Play-Console-CLI"
INSTALL_DIR="/usr/local/bin"
BINARY="gpc"

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    darwin) OS="darwin" ;;
    linux)  OS="linux" ;;
    *)
        echo "Error: unsupported operating system: $OS"
        exit 1
        ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)
        echo "Error: unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Detected platform: ${OS}/${ARCH}"

# Get latest release tag
echo "Fetching latest release..."
LATEST=$(curl -sSf "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
if [ -z "$LATEST" ]; then
    echo "Error: could not determine latest release"
    exit 1
fi

VERSION="${LATEST#v}"
echo "Latest version: ${LATEST}"

# Build download URL
ARCHIVE="gpc_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${ARCHIVE}"

# Download and extract
TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${URL}..."
curl -sSfL "$URL" -o "${TMP_DIR}/${ARCHIVE}"

echo "Extracting..."
tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"

# Install
echo "Installing to ${INSTALL_DIR}/${BINARY}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
else
    sudo mv "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
fi
chmod +x "${INSTALL_DIR}/${BINARY}"

echo ""
echo "gpc installed successfully!"
"${INSTALL_DIR}/${BINARY}" version
