#!/bin/bash

function tab () {
    local cmd=""
    local cdto="$PWD"
    local args="$@"

    echo "pwd: $PWD"
    echo "cdto: $cdto"
    echo "args: $args"

    if [ -d "$1" ]; then
        cdto=`cd "$1"; pwd`
        args="${@:2}"
    fi

    if [ -n "$args" ]; then
        cmd="$args"
    fi

    echo "cmd: $args"

    if [ $TERM_PROGRAM = "Apple_Terminal" ]; then
        echo "tab in $TERM_PROGRAM; cdto=$cdto; cmd=$cmd;"
        osascript
            -e "tell application \"Terminal\"" \
                -e "tell application \"System Events\" to keystroke \"t\" using {command down}" \
                -e "do script \"cd $cdto; clear $cmd\" in front window" \
            -e "end tell"
            > /dev/null
    elif [ $TERM_PROGRAM = "iTerm.app" ]; then
      echo "cmd: $cmd"
        echo "tab in iTerm.app; cdto: $cdto; cmd: $cmd;"
        osascript
            -e "tell application \"iTerm2\"" \
#                -e "tell current terminal" \
#                    -e "launch session \"Default Session\"" \
#                    -e "tell the last session" \
#                        -e "write text \"cd \"$cdto\"$cmd\"" \
#                    -e "end tell" \
#                -e "end tell" \
            -e "end tell" \
#            > /dev/null
    else
      echo "tab in $TERM_PROGRAM"
    fi
}

cwd=$(pwd)

tab "$cwd" sh ./serve_gae.sh
tab "$cwd" sh ./server_fb_emulator.sh
