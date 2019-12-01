#!/bin/sh
SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)

# set link

path="$SCRIPTPATH/src/github.com/leanote"
if [ ! -d "$path" ]; then
	mkdir -p "$path"
fi
rm -rf $SCRIPTPATH/src/github.com/htmambo/leanote # 先删除
ln -s ../../../../ $SCRIPTPATH/src/github.com/htmambo/leanote

# set GOPATH
export GOPATH=$SCRIPTPATH

script="$SCRIPTPATH/leanote-linux-amd64"
chmod 777 $script
$script -importPath github.com/htmambo/leanote