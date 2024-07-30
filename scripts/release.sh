#!/bin/bash

# Check if version argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <version>"
    exit 1
fi

version=$1

# Create a release tag
git tag v$version
git push origin v$version

# Publish the package
GOPROXY=proxy.golang.org go list -m github.com/bk1031/rincon-go@v$version

echo "Package released successfully for version $version"