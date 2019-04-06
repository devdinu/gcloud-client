#!/bin/bash

set -ex

_TMUX_DEF_DIR="~/.config/tmuxinator/"
_TMUX_DIR="${TMUXINATOR_CONFIG:-$_TMUX_DEF_DIR}"
_CONFIG_DIR="${CONFIG_DIR:-${HOME}/.gcloud-client/}"
GOPATH="${GOPATH:-$TMPDIR}"
_PROJ_DIR="$GOPATH/src/github.com/devdinu"

echo "Adding tmuxinator templates to ${_TMUX_DIR}"
echo "Using $GOPATH for lib installation"
echo "Project dir $_PROJ_DIR"

mkdir -p ${GOPATH}/{src,bin}
mkdir -p ${_TMUX_DIR} ${_PROJ_DIR}


[[ -d "$_PROJ_DIR/gcloud-client" ]] || (echo "sym linking dir to ${_PROJ_DIR} " && ln -s $PWD ${_PROJ_DIR})

pushd $_PROJ_DIR/gcloud-client
cp -vf ./scripts/templates/*.yml ${_TMUX_DIR}/
echo "copied tmuxinator project templates."
go install .
echo "installation completed"
popd
