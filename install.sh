#!/bin/bash

set -e

CONFIG_DIR="$HOME/.config/goai"
SRC_DIR="$(pwd)"
CONFIG_FILES_DIR="${SRC_DIR}/config"

echo "Creating configuration directory..."
mkdir -p "${CONFIG_DIR}"

echo "Copying config.yaml and prompt.txt to ${CONFIG_DIR} ..."
cp "${CONFIG_FILES_DIR}/config.yaml" "${CONFIG_DIR}/config.yaml"
cp "${CONFIG_FILES_DIR}/prompt.txt" "${CONFIG_DIR}/prompt.txt"
