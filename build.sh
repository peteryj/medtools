#!/bin/bash

VERSION=`go build && ./medtools -v`

build() {
    GOOS=$1 bee pack -f zip -exs .zip:.sh:.go:.DS_Store:.tmp
    mv medtools.zip medtools-$VERSION-$1.zip
}

build linux
build darwin
build windows
