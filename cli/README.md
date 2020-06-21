# The cli package

This package manages all the command line stuff by using the 
[cobra library](https://github.com/spf13/cobra)

| File     | Description                                 | Command                |
| -------- | ------------------------------------------- | ---------------------- |
| balances | show balances and status                    | `./tbb balances list`   |
| flags    | Global flags to use with CLI comands        | `./tbb --datadir=$HOME/.tbb`   |
| run      | Starts the HTTP service                     | `./tbb run -p=8088`   |
| tx       | Add a transaction to the blockchain         | `./tbb tx add --from=from --to=to --value=amount --data=reason` |
| version  | Version info                                | `./tbb version` |
| state    | This establishes the current state of the blockchain ||

