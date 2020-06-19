#!/bin/bash -eu

## environment related functions that can be used by other scripts, use source lib/scripts/env.sh to use
## V0.0.2 : Tim Dadd : First version

BLACK=$(tput setaf 0)
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
LIME_YELLOW=$(tput setaf 190)
YELLOW=$(tput setaf 3)
POWDER_BLUE=$(tput setaf 153)
BLUE=$(tput setaf 4)
MAGENTA=$(tput setaf 5)
CYAN=$(tput setaf 6)
ORANGE=$(tput setaf 10)
WHITE=$(tput setaf 7)
BRIGHT=$(tput bold)
NORMAL=$(tput sgr0)
BLINK=$(tput blink)
REVERSE=$(tput smso)
UNDERLINE=$(tput smul)

# Show the command before execution $1=cmd, $2=echo before
function showDoCmd() {
  if [ $# -eq 2 ]; then echo -n "$2";fi
  echo "$1"
  bash -c "$1"
}