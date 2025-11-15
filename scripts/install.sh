#!/bin/bash

set -e

CONFIG_DIR="$HOME/.config/aicmdtools"
SRC_DIR="$(pwd)"
CONFIG_FILES_DIR="${SRC_DIR}/config"

echo "Creating configuration directory..."
mkdir -p "${CONFIG_DIR}"

file_list=$(/usr/bin/ls ${CONFIG_FILES_DIR} | sed 's/^/- /')
echo -e "Copying:\n${file_list}\nto ${CONFIG_DIR} ..."

cp "${CONFIG_FILES_DIR}/config.yaml" "${CONFIG_DIR}/config.yaml"
cp "${CONFIG_FILES_DIR}/prompt.txt" "${CONFIG_DIR}/prompt.txt"
cp "${CONFIG_FILES_DIR}/chat-prompt.txt" "${CONFIG_DIR}/chat-prompt.txt"
cp "${CONFIG_FILES_DIR}/comp-graph-prompt.txt" "${CONFIG_DIR}/comp-graph-prompt.txt"
cp "${CONFIG_FILES_DIR}/aifix-prompt.txt" "${CONFIG_DIR}/aifix-prompt.txt"

echo "Configuration files have been copied to ${CONFIG_DIR}"
echo "Installation complete!"
