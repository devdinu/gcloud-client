#!/bin/bash

set -ex

DEST=$1
_TMUX_DEF_DIR="~/.config/tmuxinator"
_TMUX_DIR="${TMUXINATOR_CONFIG:-$_TMUX_DEF_DIR}"

GOPATH="$PWD/.gobuild"
_PROJ_DIR="$GOPATH/src/github.com/devdinu/gcloud-client"

echo "Adding tmuxinator templates to ${_TMUX_DIR}"
echo "Using $GOPATH for lib installation, proj dir: $_PROJ_DIR"

mkdir -p ${GOPATH}/{src,bin,pkg}
mkdir -p ${_TMUX_DIR} ${_PROJ_DIR}

cp -vf ./scripts/templates/*.yml "${_TMUX_DIR}/*"
echo "copied tmuxinator project templates."

cp -r $PWD/* "$_PROJ_DIR/"
pushd $_PROJ_DIR && go get . && GOBIN=$GOPATH/bin/ go install  && popd

cp -v "$GOPATH/bin/gcloud-client" $DEST/

rm -rf "$_PROJ_DIR"
echo "installation completed."
