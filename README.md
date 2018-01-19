# Cosmos Gaia Explorer

This project provides a blockchain explorer for the Cosmos testnet (currenty gaia-2).

http://explorer.cosmosvalidators.com/

Currently supported features:

* Balances and coin transaction history for an account [example](http://explorer.cosmosvalidators.com/#/account/7334A4B2668DE1CEF0DD7DBA695C29449EC3A0D0)
* Coin transactions and validators in a block [example](http://explorer.cosmosvalidators.com/#/block/178507)
* Raw request for each coin transaction [example](http://explorer.cosmosvalidators.com/#/tx/83527AC99E577CEF7408FA8BD2F660F7D95C69BC)
* Recent blocks [example](http://explorer.cosmosvalidators.com/)

Soon to be supported features:

* Stake transactions
* Recent transactions

## How to use

First, you need to build and deploy the [REST services](https://github.com/CyberMiles/explorer/tree/master/services). The REST services point to a gaia-2 network node (default is gaia-2-node0.testnets.interblock.io) and provide decoded transactions in JSON format. You can see a [full list of currently implemented service endpoints here](https://explorerservices.docs.apiary.io/#reference).

Second, deploy and start the [node.js web application](https://github.com/CyberMiles/explorer/tree/master/ui), which utilizes the REST services to display data.

Enjoy!
