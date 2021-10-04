#!/usr/bin/env bash
#
# Use lower_case variables in the scripts and UPPER_CASE variables for override
# Use the constants.sh for env overrides
# Use the versions.sh to specify versions
#

AVALANCHE_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd ) # Directory above this script

# Set the PATHS
GOPATH="$(go env GOPATH)"
coreth_path="$GOPATH/pkg/mod/github.com/ava-labs/coreth@$coreth_version"

# Where AvalancheGo binary goes
build_dir="$AVALANCHE_PATH/build"
binary_manager_path="$build_dir/avalanchego"

# Latest Avalanchego binary
latest_avalanchego_path="$build_dir"/avalanchego-latest
latest_avalanchego_process_path="$latest_avalanchego_path/avalanchego-process"
latest_plugin_dir="$latest_avalanchego_path/plugins"
latest_evm_path="$latest_plugin_dir/evm"

# Previous AvalancheGo binary
prev_build_dir="$build_dir/avalanchego-preupgrade" # Where pre-db migration AvalancheGo binary goes
prev_avalanchego_process_path="$prev_build_dir/avalanchego-process"
prev_plugin_dir="$prev_build_dir/plugins"

# Avalabs docker hub
# avaplatform/avalanchego - defaults to local as to avoid unintentional pushes
# You should probably set it - export DOCKER_REPO='avaplatform/avalanchego'
avalanchego_dockerhub_repo=${DOCKER_REPO:-"local"}

# Current branch
current_branch=$(git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

git_commit=${AVALANCHEGO_COMMIT:-$( git rev-list -1 HEAD )}
