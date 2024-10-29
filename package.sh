#!/bin/bash

VERSION=$1

if [ -z "$VERSION" ]; then
  echo "Error: Version not specified."
  echo "Usage: ./package.sh <version>"
  exit 1
fi

DIST_DIR="dist"
BUILD_DIR="build"

PLATFORMS=("windows" "macos" "linux")

# output path
mkdir -p "$DIST_DIR"

# binary packaging
for platform in "${PLATFORMS[@]}"; do
  TAR_FILE="$DIST_DIR/go-port-tester-$VERSION-$platform-bin.tar.gz"
  SHA_FILE="$DIST_DIR/go-port-tester-$VERSION-$platform-bin.sha512"

  echo "Packaging $platform files for version $VERSION..."
  (cd "$BUILD_DIR/$platform" && tar -czf "../../$TAR_FILE" .)

  echo "Creating SHA-512 checksum for $platform bin package..."
  sha512sum "$TAR_FILE" > "$SHA_FILE"
done

# source packaging
SOURCE_TAR_FILE="$DIST_DIR/go-port-tester-$VERSION-source.tar.gz"
SOURCE_SHA_FILE="$DIST_DIR/go-port-tester-$VERSION-source.sha512"

echo "Packaging source files for version $VERSION..."
tar --exclude="$DIST_DIR" --exclude="$BUILD_DIR" -czf "$SOURCE_TAR_FILE" .

echo "Creating SHA-512 checksum for source package..."
sha512sum "$SOURCE_TAR_FILE" > "$SOURCE_SHA_FILE"

echo "Packaging completed for version $VERSION. Files are located in the $DIST_DIR directory."
