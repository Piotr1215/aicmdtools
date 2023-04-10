#!/bin/bash

version_file="./cmd/goai/main.go"

# Get the current version from the version.go file
current_version=$(grep -oP 'version = "\K[^"]+' $version_file)

# Increment the version number
IFS='.' read -ra version_parts <<<"$current_version"
((version_parts[2]++))
new_version="${version_parts[0]}.${version_parts[1]}.${version_parts[2]}"

# Update the version.go file with the new version number
sed -i "s/version = \"$current_version\"/version = \"$new_version\"/g" $version_file

# Print the new version
echo "Updated version: $new_version"
