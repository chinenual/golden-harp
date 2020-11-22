#! /bin/bash

VERSION=`grep Version version.go | sed -e 's/^.* //' -e 's/"//g'`
tgt=GoldenHarpManager-${VERSION}.msi

rm -f $tgt

wixl -v \
	-a x86 \
	-D VERSION=${VERSION} \
	-D SourceDir=. \
	-o $tgt \
	windows-installer.wxs

