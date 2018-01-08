## Installation

```shell
go get github.com/cybermiles/explorer/services
cd $GOPATH/src/github.com/cybermiles/explorer/services
make all
```

## Initializing to connect to public node in testnets:

```shell
explorercli init --chain-id=gaia-2 --node=tcp://gaia-2-node0.testnets.interblock.io:46657
```

## Or connect to our own gaia node

1. Download configuration files for testnets:

   ```shell
   git clone https://github.com/tendermint/testnets $HOME/testnets
   GAIANET=$HOME/testnets/gaia-2/gaia
   cd $GAIANET
   ```

2. Start our gaia node:

   ```
   gaia start --home=$GAIANET
   ```

3. Initializing to connect to our own node:

   ```shell
   explorercli init --chain-id=gaia-2 --node=tcp://localhost:46657
   ```

## Start rest server

```shell
explorercli rest-server
```
