#!/bin/bash
# WEBPRISM Third-Party Dependencies Setup Script
#
# This script downloads required third-party proto files for building WEBPRISM.
# It ensures consistent versions across all environments.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
THIRD_PARTY_DIR="$PROJECT_ROOT/third_party"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🔧 Setting up third-party dependencies for WEBPRISM..."
echo ""

# ============================================
# googleapis - Google API Definitions
# ============================================
# Repository: https://github.com/googleapis/googleapis
# We use a specific commit to ensure reproducibility

GOOGLEAPIS_REPO="https://github.com/googleapis/googleapis.git"
GOOGLEAPIS_COMMIT="2f9af297c84c55c8b871ba4495e01ade42476c92" # Latest stable as of 2024-01
GOOGLEAPIS_DIR="$THIRD_PARTY_DIR/googleapis"

echo "📦 Downloading googleapis proto files..."
echo "   Repository: $GOOGLEAPIS_REPO"
echo "   Commit: $GOOGLEAPIS_COMMIT"
echo ""

# Check if googleapis directory exists
if [ -d "$GOOGLEAPIS_DIR" ]; then
    echo -e "${YELLOW}⚠️  googleapis directory already exists${NC}"
    read -p "   Do you want to re-download? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${GREEN}✓ Using existing googleapis${NC}"
        echo ""
        echo "Third-party dependencies setup complete!"
        exit 0
    fi
    echo "   Removing existing directory..."
    rm -rf "$GOOGLEAPIS_DIR"
fi

# Create third_party directory if it doesn't exist
mkdir -p "$THIRD_PARTY_DIR"

# Clone googleapis repository (sparse checkout for efficiency)
echo "   Cloning googleapis repository..."
cd "$THIRD_PARTY_DIR"

# Use sparse checkout to download only what we need
git clone --filter=blob:none --no-checkout --depth=1 --single-branch "$GOOGLEAPIS_REPO" googleapis
cd googleapis
git checkout "$GOOGLEAPIS_COMMIT"

# Sparse checkout: only get google/api and google/rpc
git sparse-checkout init --cone
git sparse-checkout set google/api google/rpc

echo ""
echo -e "${GREEN}✓ googleapis downloaded successfully${NC}"

# Verify required files exist
echo ""
echo "🔍 Verifying required proto files..."

REQUIRED_FILES=(
    "google/api/annotations.proto"
    "google/api/http.proto"
    "google/api/field_behavior.proto"
)

ALL_OK=true
for file in "${REQUIRED_FILES[@]}"; do
    if [ -f "$GOOGLEAPIS_DIR/$file" ]; then
        echo -e "${GREEN}✓${NC} $file"
    else
        echo -e "${RED}✗${NC} $file (missing)"
        ALL_OK=false
    fi
done

echo ""
if [ "$ALL_OK" = true ]; then
    echo -e "${GREEN}✓ All required proto files are present${NC}"
    echo ""
    echo "Third-party dependencies setup complete!"
    echo ""
    echo "You can now run:"
    echo "  make proto    # Generate protobuf code"
    echo "  make build    # Build the project"
else
    echo -e "${RED}✗ Some required files are missing${NC}"
    echo "Please check the setup or report this issue."
    exit 1
fi
