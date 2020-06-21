# the HTTTP API

[Why is this package called `node`?](https://en.bitcoin.it/wiki/Full_node)
 
This provides a simple RESTful API as follows:

##  http://.../balances/list
Provides a json list of account balances
### Example JSON Response
```json
{
   "balances" : {
      "andrej" : 8851,
      "babayaga" : 1049,
      "caesar" : 1000,
      "tim" : 20000
   },
   "block_hash" : "5591d6cea7ff917d1b5c3a827e43821e800a228d0fcfa516b01d71e4c705919e"
}
```

##  http://.../tx/add
This adds a transaction to the blockchain.  The body of the 'POST' provides details of the transaction.
### Example JSON Request
```json
{
  "from": "andrej",
  "to": "babayaga",
  "value": 100
}
```
### Example JSON Response
```json
{
  "hash":"5591d6cea7...",
  "block" : {
    "header" : {
      "parent":"7e97e062f587...",
      "time":1592716425
    },
    "payload": [{
      "from":"andrej",
      "to":"babayaga",
      "value":100,
      "data":""
    }]
  }
}
```
```json
{
  "block_hash" : "5591d6cea7f..."
}
```
