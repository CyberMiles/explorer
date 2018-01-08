package handlers

import (
  "fmt"
  "net/http"

  "github.com/tendermint/go-wire/data"
  rpcclient "github.com/tendermint/tendermint/rpc/client"

  "github.com/cosmos/cosmos-sdk/client"
  "github.com/cosmos/cosmos-sdk/client/commands"
)

func printResult(w http.ResponseWriter, res interface{}) error {
  json, err := data.ToJSON(res)
  if err != nil {
    return err
  }
  _, err = fmt.Fprintf(w, "%s\n", json)
  return err
}

func getSecureNode() (rpcclient.Client, error) {
  c := commands.GetNode()
  cert, err := commands.GetCertifier()
  if err != nil {
    return nil, err
  }
  return client.SecureClient(c, cert), nil
}
