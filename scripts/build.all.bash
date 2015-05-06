#!/bin/bash

if [ "$(hostname)" != "golang" ]; then
  echo "ERROR: This is meant to be run inside container specified "
  echo "       by Buildme.yml"
  exit 1
fi

# only supporting more common archs for now, let me know if you need
# more, or add it to the list
archs="darwin_386 darwin_amd64 freebsd_386 freebsd_amd64 linux_386 linux_amd64 windows_386 windows_amd64"

if [ "$1" != "" ]; then
  archs=$1
fi

for arch in $archs
do
  split=(${arch//_/ })
  goos=${split[0]}
  goarch=${split[1]}

  src=bin/dockme.go
  target=builds/$goos/$goarch

  [[ "windows" == "$goos" ]] && target=$target/dockme.exe

  if ! test -d "$(dirname $target)"
  then
    echo "mkdir -pv $(dirname $target)"
    mkdir -pv $(dirname $target)
  fi

  echo "GOOS=$goos GOARCH=$goarch go build -x -o $target $src"
  GOOS=$goos GOARCH=$goarch go build -x -o $target/dockme $src
  md5sum $target/dockme | tee $target/dockme.md5
done

