#!/bin/bash

# Build script for TerraDrift Watcher
# Builds binaries for multiple platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get version from git tag or use dev
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo -e "${GREEN}Building TerraDrift Watcher ${VERSION}${NC}"
echo "Commit: ${COMMIT}"
echo "Date: ${DATE}"

# Create dist directory
mkdir -p dist

# Build flags
LDFLAGS="-s -w -X 'github.com/terradrift-watcher/cmd.version=${VERSION}' -X 'github.com/terradrift-watcher/cmd.commit=${COMMIT}' -X 'github.com/terradrift-watcher/cmd.date=${DATE}'"

# Function to build for a specific platform
build_platform() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT=$3
    
    echo -e "${YELLOW}Building for ${GOOS}/${GOARCH}...${NC}"
    
    GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="${LDFLAGS}" -o "dist/${OUTPUT}" .
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Built ${OUTPUT}${NC}"
    else
        echo -e "${RED}✗ Failed to build ${OUTPUT}${NC}"
        exit 1
    fi
}

# Build for all platforms
build_platform "linux" "amd64" "terradrift-watcher-linux-amd64"
build_platform "linux" "arm64" "terradrift-watcher-linux-arm64"
build_platform "darwin" "amd64" "terradrift-watcher-darwin-amd64"
build_platform "darwin" "arm64" "terradrift-watcher-darwin-arm64"
build_platform "windows" "amd64" "terradrift-watcher-windows-amd64.exe"

# Create checksums
echo -e "${YELLOW}Creating checksums...${NC}"
cd dist
if command -v sha256sum &> /dev/null; then
    sha256sum * > checksums.txt
elif command -v shasum &> /dev/null; then
    shasum -a 256 * > checksums.txt
else
    echo -e "${RED}Warning: No checksum tool found${NC}"
fi
cd ..

# Create archives for release
echo -e "${YELLOW}Creating release archives...${NC}"
cd dist
for file in terradrift-watcher-*; do
    if [[ "$file" == *.exe ]]; then
        zip "${file%.exe}.zip" "$file" ../config.example.yml ../README.md ../LICENSE
        echo -e "${GREEN}✓ Created ${file%.exe}.zip${NC}"
    else
        tar czf "${file}.tar.gz" "$file" ../config.example.yml ../README.md ../LICENSE
        echo -e "${GREEN}✓ Created ${file}.tar.gz${NC}"
    fi
done
cd ..

echo -e "${GREEN}Build complete! Binaries are in the dist/ directory.${NC}"
echo ""
echo "To test the binary for your current platform:"
echo "  ./dist/terradrift-watcher-$(go env GOOS)-$(go env GOARCH)$(if [[ $(go env GOOS) == 'windows' ]]; then echo '.exe'; fi) --version" 