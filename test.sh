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
    echo "   8 ..... HTTP API"
    echo "     ..... No number, build only"
}

## tx_postman 1:from 2:to 3:value 4:data 5:colour
tx_postman () {
  cmd=$'curl -s --location --request POST \'http:/localhost:8080/tx/add\'
  --header \'Content-Type: application/json\'
  --data-raw \'{"from": "'$1'","to": "'$2'","value": '$3"}'"
echo $5$cmd
echo $cmd | bash
echo
#bash -c "$cmd"
}

RED=$(tput setaf 1)
YELLOW=$(tput setaf 3)
CYAN=$(tput setaf 6)

projd=$PWD # make a note of the project directory
if [ ! -d $projd/cli ];then
  echo "${RED}Why is there no ${YELLOW}cli${RED} directory?"
  echo "${CYAN}This script should be run from project root${WHITE}"
  exit 1
fi

datad="$projd/data"

source scripts/env.sh
source scripts/go.sh

go mod tidy
go mod vendor
goFmt
goVet
showDoCmd "go build -o tbb" $GREEN

# Every blockchain has a "Genesis" file. The Genesis file is used to distribute
# the first tokens to early blockchain participants.
echo $CYAN"Clear database"
killall tbb || true
showDoCmd "rm -f $datad/thispeernode.json"
showDoCmd "rm -f $datad/db/*.*"

showDoCmd "./tbb version" $CYAN$'\n'
showDoCmd "./tbb balances list" $YELLOW$'\n'

if [ $# -eq 0 ];then exit;fi

chapter=$1
if [ $chapter -lt 3 -a $chapter -gt 8 ]; then
  show_help
  exit
fi



if [ $chapter -ge 3 ]; then
  echo $WHITE"Running Chapter 3 - First customer"
  ## Andrej purchases 3 shots of vodka from his own bar
  showDoCmd "./tbb tx add --from=andrej --to=andrej --value=3 --data=vodka" ${POWDER_BLUE}
  if [ $chapter -eq 3 ]; then showDoCmd "./tbb balances list";fi

  # Andrej also decides he should be getting 100 tokens per day for maintaining
  # the database and having such a brilliant disruptive idea. (700 per week)
  showDoCmd "./tbb tx add --from=andrej --to=andrej --value=700 --data=reward" ${POWDER_BLUE}
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
  showDoCmd "./tbb tx add --from=andrej --to=andrej --value=100 --data=reward" ${POWDER_BLUE}
  echo $GREEN"Chapter 3 processed"

  showDoCmd "./tbb balances list" $YELLOW
fi


if [ $chapter -ge 4 ]; then
  echo $WHITE"Running Chapter 4 - BabaYaga pays rent to Caesar and Andrej takes his cut"
  # Rent payment
  showDoCmd "./tbb tx add --from=babayaga --to=caesar --value=1000 --data=rent" ${LIME_YELLOW}
  if [ $chapter -eq 4 ]; then showDoCmd "./tbb balances list";fi

  # Hidden transaction charge
  showDoCmd "./tbb tx add --from=babayaga --to=andrej --value=50 --data=hidden_fee" ${RED}
  if [ $chapter -eq 4 ]; then showDoCmd "./tbb balances list";fi

  # Another rewarding day
  showDoCmd "./tbb tx add --from=andrej --to=andrej --value=100 --data=reward" ${POWDER_BLUE}
  echo $GREEN"Chapter 4 processed"
  showDoCmd "./tbb balances list"

  if [ $chapter -eq 4 ]; then showDoCmd "./tbb balances state";fi
fi

if [ $chapter -ge 8 ]; then
  echo $WHITE"Running Chapter 8 - Andrej pays BabaYaga 100 units via the RESTful API"
  echo "${WHITE}Starting the node"
  ./tbb run &
  sleep 1
  showDoCmd "curl -s --http2 http://localhost:8080/balances/list | json_pp" $CYAN
  tx_postman andrej babayaga 100 gift $POWDER_BLUE
  ## This next line shows the wrong balance because the state is persisted in memory of other the API process
  echo $GREEN"Chapter 8 processed"
  showDoCmd "curl -s --http2 http://localhost:8080/balances/list | json_pp"
  ## Burn out compensation
  showDoCmd "./tbb tx add --from=andrej --to=andrej --value=24700 --data=reward" ${POWDER_BLUE}
  echo $GREEN"Burn out compensation added"
  showDoCmd "./tbb balances list"
fi

if [ $chapter -ge 10 ]; then
  echo $WHITE"Running Chapter 10 - Peer-to-Peer DB Sync"
  showDoCmd "curl -s --http2 curl -X GET http://localhost:8080/node/status | json_pp" $CYAN
#  tx_postman andrej babayaga 100 gift $POWDER_BLUE
#  ## This next line shows the wrong balance because the state is persisted in memory of other the API process
#  echo $GREEN"Chapter 8 processed"
#  showDoCmd "curl -s --http2 http://localhost:8080/balances/list | json_pp"
#  ## Burn out compensation
#  showDoCmd "./tbb tx add --from=andrej --to=andrej --value=24700 --data=reward" ${POWDER_BLUE}
#  echo $GREEN"Burn out compensation added"
  showDoCmd "./tbb balances list"
fi




if [ $chapter ]; then echo $WHITE"All done up until chapter $chapter";fi

blockdb="$datad/db/block.db"
echo $blockdb
if [ -f $blockdb ]; then tail $blockdb;fi

