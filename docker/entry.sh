#!/usr/bin/env ash

COMMAND=$1
PATH=/app/build:$PATH

case $COMMAND in
    "run")
        exec /usr/local/bin/nudge
    ;;
    
    *)
        exec $@
    ;;
esac
