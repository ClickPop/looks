#!/bin/bash

OS=$(uname)
ARCH=$(uname -m)
LOOKS_PATH=~/.looks

handleOS() {
  if [[ "$OS" != "Darwin" && "$OS" != "Linux" ]] ; then
    echo "unsupported operating system"
  else
    handleArch
  fi
}

getLooks() {
  mkdir ~/.looks
  lowerCase "looks-$OS-$ARCH"
  curl -L https://github.com/clickpop/looks/releases/latest/download/$FILENAME >> "$LOOKS_PATH/looks"
  chmod 777 "$LOOKS_PATH/looks"
}

handlePath() {
  IN_PATH=false
  IFS=':' read -ra ADDR <<< "$PATH"
  for i in "${ADDR[@]}"; do
    if [[ "$i" == "$LOOKS_PATH" ]] ; then
      IN_PATH=true
    fi
  done
  if [[ ! $IN_PATH ]] ; then
    if [[ -f ~/.bashrc ]] ; then
      echo "profile"
      handlePathPersist $HOME/.bashrc
    fi

    if [[ -f ~/.zshrc ]] ; then
      handlePathPersist $HOME/.zshrc
    fi
  fi
}

handlePathPersist() {
  IN_FILE=false
  NEW_PATH="export PATH=\"\$PATH:$LOOKS_PATH\""
  if grep -q "$NEW_PATH" $1 ; then
    echo "In file"
    IN_FILE=true
  fi
  if [[ ! $IN_FILE ]] ; then
    echo $NEW_PATH >> $1
  fi
}

handleArch() {
  if [[ "$ARCH" != "arm64" && "$ARCH" != "x84_64" ]] ; then
    echo "unsupported cpu architecture"
  else
    if [[ "$ARCH" == "x84_64" ]] ; then
      ARCH="amd64"
    fi
    getLooks
    handlePath
  fi
}

lowerCase() {
  FILENAME=$(echo "$1" | awk '{print tolower($0)}')
}

handleOS
echo -e "\nplease restart your terminal to use looks"