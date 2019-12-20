#!/usr/bin/env bash

set -euo pipefail

for f in pkg/*; do
  if [[ -d $f ]]; then
    i=$(basename $f)
    echo
    echo === Testing pkg $i
    dir=$GOPATH/src/github.com/andy2046/gopie/pkg/$i
    cd $dir
    GOCACHE=off go test -v -race
    go fmt
    if [[ $i =~ "nocopy" || $i =~ "spinlock" ]]; then
      echo "ignore go vet"
    else
      go vet
    fi
    golint
    cd -
    echo === Tested pkg $i
    echo
  fi
done;
