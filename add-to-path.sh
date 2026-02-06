#!/bin/bash

# Script to add Go bin to PATH

GOPATH_BIN="$(go env GOPATH)/bin"
SHELL_RC=""

# Detect shell
if [ -n "$ZSH_VERSION" ]; then
    SHELL_RC="$HOME/.zshrc"
elif [ -n "$BASH_VERSION" ]; then
    SHELL_RC="$HOME/.bashrc"
else
    echo "⚠️  Could not detect shell. Please manually add the following to your shell config:"
    echo "export PATH=\"\$PATH:$GOPATH_BIN\""
    exit 0
fi

# Check if already in PATH
if echo "$PATH" | grep -q "$GOPATH_BIN"; then
    echo "✅ $GOPATH_BIN is already in your PATH"
    exit 0
fi

# Check if already in shell config
if [ -f "$SHELL_RC" ] && grep -q "GOPATH.*bin" "$SHELL_RC"; then
    echo "✅ Go bin path is already configured in $SHELL_RC"
    echo "   Please run: source $SHELL_RC"
    exit 0
fi

echo "📝 Adding Go bin to PATH in $SHELL_RC"

# Add to shell config
cat >> "$SHELL_RC" << 'EOF'

# Added by noaas-default-route setup
export PATH="$PATH:$(go env GOPATH)/bin"
EOF

echo "✅ Added to $SHELL_RC"
echo ""
echo "To use golangci-lint directly, run:"
echo "  source $SHELL_RC"
echo ""
echo "Or start a new terminal session."
echo ""
echo "Note: The Makefile already uses the full path, so 'make lint' will work without this."
