#!/bin/bash

set -e

CONFIG_DIR="$HOME/.config/aicmdtools"
- [1f79bf9](git@github.com:Piotr1215/aicmdtools.git/1f79bf9@) - (HEAD -> main, origin/main, origin/HEAD) further improve instructions
- [e0c628b](git@github.com:Piotr1215/aicmdtools.git/e0c628b@) - clarify install instructions
- [36d9b3a](git@github.com:Piotr1215/aicmdtools.git/36d9b3a@) - (tag: v1.0.0) installation process fix
- [de73413](git@github.com:Piotr1215/aicmdtools.git/de73413@) - readme
- [1ca4785](git@github.com:Piotr1215/aicmdtools.git/1ca4785@) - upgrade readme
- [d879d5c](git@github.com:Piotr1215/aicmdtools.git/d879d5c@) - add command to see models
- [f276715](git@github.com:Piotr1215/aicmdtools.git/f276715@) - trim response
- [3a89424](git@github.com:Piotr1215/aicmdtools.git/3a89424@) - pass through the model
- [8ecb34b](git@github.com:Piotr1215/aicmdtools.git/8ecb34b@) - Revert "Removing bubbletea, too many things don't work"
- [27b17fa](git@github.com:Piotr1215/aicmdtools.git/27b17fa@) - Removing bubbletea, too many things don't work
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
