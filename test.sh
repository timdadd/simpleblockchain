#!/bin/bash -eu

# check https://www.gnu.org/software/bash/manual/html_node/The-Set-Builtin.html#The-Set-Builtin for options
# Exit immediately if error, unset varible usage is an error, if a pipe command fails then use output of previous
set -o nounset
set -o errexit
set -o pipefail

show_help () {
    echo "Usage: $0 {chapter number}" >&2
    echo "   3 ..... First customer, BabaYaga"
    echo "   4 ..... BabaYaga pays rent to Caesar"
    echo "     ..... No number, build only"
}

RED=$(tput setaf 1)
YELLOW=$(tput setaf 3)
CYAN=$(tput setaf 6)

projd=$PWD # make a note of the project directory
if [ ! -d $projd/cmd ];then
  echo "${RED}Why is there no ${YELLOW}cmd${RED} directory?"
  echo "${CYAN}This script should be run from project root${WHITE}"
  exit 1
fi

source scripts/env.sh
source scripts/go.sh

# Every blockchain has a "Genesis" file. The Genesis file is used to distribute
# the first tokens to early blockchain participants.
showDoCmd "cp -f baseline/genesis.json db"
showDoCmd "cp -f baseline/tx.db db"

go mod tidy
go mod vendor
goFmt
goVet
showDoCmd "go build -o tbb" $GREEN
showDoCmd "./tbb version" $CYAN$'\n'
showDoCmd "./tbb balances list" $YELLOW$'\n'

if [ $# -eq 0 ];then exit;fi

chapter=$1
if [ $chapter -lt 3 -a $chapter -gt 4 ]; then
  show_help
  exit
fi

echo $WHITE"Running Test 3 - First customer"

showDoCmd "rm db/tx.db"
showDoCmd "./tbb balances list" $YELLOW

## Andrej purchases 3 shots of vodka from his own bar
showDoCmd "./tbb tx add --from=andrej --to=andrej --value=3 --data=vodka" ${POWDER_BLUE}
if [ $chapter -eq 3 ]; then showDoCmd "./tbb balances list";fi

# Andrej also decides he should be getting 100 tokens per day for maintaining
# the database and having such a brilliant disruptive idea. (700 per week)
showDoCmd "./tbb tx add --from=andrej --to=andrej --value=700 --data=services" ${POWDER_BLUE}
if [ $chapter -eq 3 ]; then showDoCmd "./tbb balances list";fi

# To bring traffic to his bar, Andrej announces an exclusive 100% bonus for everyone who
# purchases the TBB tokens in the next 24 hours.
# Bingo! He gets his first customer called BabaYaga. BabaYaga pre-purchases 1000€ worth of tokens
showDoCmd "./tbb tx add --from=andrej --to=babayaga --value=2000  --data=1000€" ${LIME_YELLOW}
if [ $chapter -eq 3 ]; then showDoCmd "./tbb balances list";fi

# She immediately spends 1 TBB for a vodka shot.
showDoCmd "./tbb tx add --from=babayaga --to=andrej --value=1 --data=vodka" ${LIME_YELLOW}
if [ $chapter -eq 3 ]; then showDoCmd "./tbb balances list";fi

# Another rewarding day
showDoCmd "./tbb tx add --from=andrej --to=andrej --value=100 --data=services" ${POWDER_BLUE}
echo $WHITE"Chapter 3 processed"
showDoCmd "./tbb balances list"

if [ $chapter -eq 3 ]; then
  showDoCmd "./tbb balances state"
  exit
fi

echo $WHITE"Running Chapter 4 - BabaYaga pays rent to Caesar and Andrej takes his cut"
# Rent payment
showDoCmd "./tbb tx add --from=babayaga --to=caesar --value=1000 --data=rent" ${LIME_YELLOW}
if [ $chapter -eq 4 ]; then showDoCmd "./tbb balances list";fi

# Hidden transaction charge
showDoCmd "./tbb tx add --from=babayaga --to=andrej --value=50 --data=hidden_fee" ${RED}
if [ $chapter -eq 4 ]; then showDoCmd "./tbb balances list";fi

# Another rewarding day
showDoCmd "./tbb tx add --from=andrej --to=andrej --value=100 --data=services" ${POWDER_BLUE}
echo $WHITE"Chapter 4 processed"
showDoCmd "./tbb balances list"

if [ $chapter -eq 4 ]; then showDoCmd "./tbb balances state";fi
