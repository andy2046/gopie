#!/usr/bin/env bash

set -euo pipefail

echo 'usage: make [target] ...'
echo
echo 'targets:'
fgrep -h "##" ./Makefile | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
