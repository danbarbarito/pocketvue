#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Build directory
BUILD_DIR="builds"
BINARY_NAME="pocketvue"

echo -e "${GREEN}Starting build process...${NC}"

# Step 1: Build frontend
echo -e "${YELLOW}Step 1: Building frontend...${NC}"
cd frontend
pnpm build
cd ..

# Step 2: Create builds directory
echo -e "${YELLOW}Step 2: Creating builds directory...${NC}"
mkdir -p "$BUILD_DIR"

# Get absolute path for build directory
ABS_BUILD_DIR="$(cd "$BUILD_DIR" && pwd)"

# Step 3: Build backend for all platforms
echo -e "${YELLOW}Step 3: Building backend binaries...${NC}"

# Change to backend directory
cd backend

# Define platforms and architectures
GOOS_LIST=("linux" "windows" "darwin")
GOARCH_LIST=("amd64" "arm64" "arm" "s390x" "ppc64le")

# Function to check if combination should be ignored
should_ignore() {
    local goos=$1
    local goarch=$2
    
    # Windows exclusions
    if [[ "$goos" == "windows" ]]; then
        if [[ "$goarch" == "arm" ]] || [[ "$goarch" == "s390x" ]] || [[ "$goarch" == "ppc64le" ]]; then
            return 0
        fi
    fi
    
    # Darwin exclusions
    if [[ "$goos" == "darwin" ]]; then
        if [[ "$goarch" == "arm" ]] || [[ "$goarch" == "s390x" ]] || [[ "$goarch" == "ppc64le" ]]; then
            return 0
        fi
    fi
    
    return 1
}

# Build for each platform/arch combination
BUILT_COUNT=0
SKIPPED_COUNT=0

for goos in "${GOOS_LIST[@]}"; do
    for goarch in "${GOARCH_LIST[@]}"; do
        if should_ignore "$goos" "$goarch"; then
            echo -e "${YELLOW}Skipping ${goos}/${goarch} (excluded combination)${NC}"
            ((SKIPPED_COUNT++))
            continue
        fi
        
        # Set binary extension
        extension=""
        if [[ "$goos" == "windows" ]]; then
            extension=".exe"
        fi
        
        # Build binary name
        output_name="${BINARY_NAME}-${goos}-${goarch}${extension}"
        output_path="${ABS_BUILD_DIR}/${output_name}"
        
        echo -e "${GREEN}Building ${goos}/${goarch}...${NC}"
        
        # Build command (we're already in backend directory, so build from current dir)
        # -ldflags="-s -w" strips debug info and symbols, reducing binary size by ~25%
        if [[ "$goarch" == "arm" ]]; then
            env CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" GOARM=7 \
                go build -ldflags="-s -w" -o "$output_path" .
        else
            env CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
                go build -ldflags="-s -w" -o "$output_path" .
        fi
        
        if [[ $? -eq 0 ]]; then
            echo -e "${GREEN}✓ Built: ${output_name}${NC}"
            ((BUILT_COUNT++))
        else
            echo -e "${RED}✗ Failed: ${output_name}${NC}"
        fi
    done
done

# Return to root directory
cd ..

echo ""
echo -e "${GREEN}Build completed!${NC}"
echo -e "  Built: ${BUILT_COUNT} binaries"
echo -e "  Skipped: ${SKIPPED_COUNT} combinations"
echo -e "  Output directory: ${BUILD_DIR}/"

