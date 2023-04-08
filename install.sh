#!/bin/bash

set -e

CONFIG_DIR="$HOME/.config/goai"
SRC_DIR="$(pwd)"

echo "Installing goai..."
go install .

echo "Creating configuration directory..."
mkdir -p "${CONFIG_DIR}"

echo "Copying yolo.yaml and prompt.txt to ${CONFIG_DIR} ..."
cp "${SRC_DIR}/yolo.yaml" "${CONFIG_DIR}/yolo.yaml"
cp "${SRC_DIR}/prompt.txt" "${CONFIG_DIR}/prompt.txt"
