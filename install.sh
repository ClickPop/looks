#!/bin/bash

OS=$(uname)
ARCH=$(uname -m)

handleOS() {
  if [[ "$OS" != "Darwin" && "$OS" != "Linux" ]] ; then
    echo "unsupported operating system"
  else
    handleArch
  fi
}

handleArch() {
  if [[ "$ARCH" != "arm64" && "$ARCH" != "x84_64" ]] ; then
    echo "unsupported cpu architecture"
  else
    if [[ "$ARCH" == "x84_64" ]] ; then
      ARCH="amd64"
    fi
    mkdir ~/.looks
    lowerCase "looks-$OS-$ARCH"
    curl -L https://github.com/clickpop/looks/releases/latest/download/$FILENAME >> ~/.looks/looks
    chmod 777 ~/.looks/looks
  fi
}

lowerCase() {
  FILENAME=$(echo "$1" | awk '{print tolower($0)}')
}

handleOS