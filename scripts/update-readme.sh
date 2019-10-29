#!/bin/bash

# Borrowed from fogg.
# I would have written this directly in the Makefile, but that was difficult.

CMD="$1"

TMP=`mktemp`
TMP2=`mktemp`
./"$BASE_BINARY_NAME" -doc > "$TMP"
sed '/^<!-- START -->$/,/<!-- END -->/{//!d;}' README.md | sed "/^<!-- START -->$/r $TMP" > $TMP2

case "$CMD" in
    update)
        mv $TMP2 README.md
    ;;
    check)
        diff $TMP2 README.md >/dev/null
    ;;
esac

exit $?
