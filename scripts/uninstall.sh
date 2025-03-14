#!/bin/sh
set -e

BIN_NAME="pdfmc"
INSTALL_DIR="/usr/local/bin"

# Check if the binary exists
if [ -f "$INSTALL_DIR/$BIN_NAME" ]; then
    echo "Removing $BIN_NAME from $INSTALL_DIR..."
    rm -f "$INSTALL_DIR/$BIN_NAME"
    echo "$BIN_NAME uninstalled successfully."
else
    echo "$BIN_NAME is not installed in $INSTALL_DIR."
fi
